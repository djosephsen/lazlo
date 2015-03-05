package lib

import (
)

var Logger = newLogger()
const M = "messages"
const E = "events"
const T = "timers"
const L = "links"

type Broker struct {
	SlackMeta		*ApiResponse
	Config			*Config
	Socket			*websocket.Conn
	Modules			map[string] *Module
	Brain				*Brain
	ApiResponses	map[int32]chan map[string]interface{}
	CallbackIndex	map[string]map[string]*interface{}
	ReadFilters    []*ReadFilter
	WriteFilters   []*WriteFilter
	MID				int32
	WriteThread		*WriteThread
	SigChan			chan os.Signal
	SyncChan			chan bool
}

type Module struct{
	Name				string
	Usage				string
	Run				func(Broker)
	SigChan			chan os.Signal
}		   

type WriteThread struct{
	broker			*Broker
	Chan           chan Event
	SyncChan        chan bool
}

type ReadFilter struct{
	Name				string
	Usage				string
	Run      func(thingy map[string]interface{}) map[string]interface{}
}		   

type WriteFilter struct{
	Name				string
	Usage				string
	Run      func(e *Event)
}		   

func NewBroker() (*Broker, error){
//return a Broker instance

	broker := &Broker{
		Config:			newConfig(),
		ApiResponses:   make(map[int32]chan map[string]interface{}),
		MessageChans:   make([]chan Event),
		EventChans:     make([]chan map[string]interface{}),
		TimerChans:     make([]chan time.Time),
		LinkChans:      make([]chan map[string]interface{}),
		MID:            0,
		WriteThread:    &WriteThread{
			Chan:			make(chan Event),
			SyncChan:		make(chan Bool),
		},
		SigChan:        make(chan os.Signal),
		SyncChan:       make(chan bool),
	}
	broker.WriteThread.Broker=broker

	//connect to slack and establish an RTM websocket
	socket,meta,err := getMeASocket(broker)
	if err != nil{
		return err
	}
	broker.SlackMeta=&meta
	broker.Socket=&socket

	var brain Brain
	brain,err = newBrain()
	if err != nil{
		return err
	}
	broker.Brain=&brain
	return broker,nil
}

func (b *broker) Start(){
	go b.WriteThread.Start()
	for {
		thingy := make(map[string]interface{})
		b.Ws.ReadJSON(&thingy)
		go b.This(thingy)
}

func (w *WriteThread) Start(){
   Logger.Debug(`Write-Thread Started`)
   stop := false
   for !stop {
      select{
      case e := <-w.Chan:
         Logger.Debug(`WriteThread:: Outbound `,e.Type,` channel: `,e.Channel,`. text: `,e.Text)
         if ejson, _ := json.Marshal(e); len(ejson) >= 16000 {
            e = Event{
            ID: e.ID,
            Type: e.Type,
            Channel: e.Channel,
            Text: fmt.Sprintf("ERROR! Response too large. %v Bytes!", len(ejson)),
            }
         }
            w.Broker.Ws.WriteJSON(e)
            time.Sleep(time.Second * 1)
      case stop = <- w.SyncChan:
         stop = true
         }
      }
   b.SyncChan <- true
}

 //probably need to make this thread-safe (for now only the write thread uses it)
func (b *Broker) NextMID() int32{
   b.MID += 1
   Logger.Debug(`incrementing MID to `, b.MID)
   return b.MID
}

func (b *Broker) This(thingy map[string]interface{}){
   if b.Modules == nil{ return }
//run the pre-handeler filters
   if b.ReadFilters != nil{
      for _,filter := range b.ReadFilters{ //run the read filters
         thingy = filter.Run(thingy)
      }
   }
   // stop here if a prefilter delted our thingy
   if len(thingy) == 0 { return }

   // if it's an api response send it to whomever is listening for it
   if replyVal, isReply := thingy[`reply_to`]; isReply{
      if replyVal != nil{ // sometimes the api returns: "reply_to":null
         b.handleApiReply(thingy)
      }
   }

	typeOfThingy:=thingy[`type`]
   switch typeOfThingy{
   case nil:
      return
   case `message`:
      b.handleMessage(thingy)
   default:
      b.handleEvent(thingy)
   }
}

func (b *Broker) Register(things ...interface{}){
   for _,thing := range things{
      switch t := thing.(type) {
      case Module:
			m:=thing.(Module)
         Logger.Debug(`registered Module: `,m.Name)
         b.Modules=append(b.Modules, &m)
 		case ReadFilter:
         r:=thing.(ReadFilter)
         Logger.Debug(`registered Read Filter: `, r.Name)
         b.PreFilters=append(b.ReadFilters, &r)
      case WriteFilter:
         w:=thing.(WriteFilter)
         Logger.Debug(`registered Write Filter: `, w.Name)
         b.WriteFilters=append(b.WriteFilters, &w)
      default:
         weirdType:=fmt.Sprintf(`%T`,t)
         Logger.Error(`sorry I cant register this handler because I don't know what a `,weirdType, ` is`)
      }
   }
}

func (b *Broker) StartModules(){
	for _,module := range b.Modules{
		go module.Run(b)
	}
}

func (b *Broker) handleApiReply(thingy map[string]interface{}){
   chanID:=int32(thingy[`reply_to`].(float64))
   Logger.Debug(`Broker:: reply message, to: `, thingy[`reply_to`])
   if callBackChannel, exists := b.APIResponses[chanID]; exists{
      callBackChannel <- thingy
      delete(b.APIResponses,chanID) //dont leak channels
      Logger.Debug(`deleted callback: `,chanID)
   } else {
      Logger.Debug(`no such channel: `,chanID)
   }
}

func (b *Broker) handleMessage(thingy map[string]interface{}){
		if b.messageCallbacks == nil { return }
      message := new(Event)
      jthingy,_ := json.Marshal(thingy)
      json.Unmarshal(jthingy, message)
      message.Broker = b
   	botNamePat := fmt.Sprintf(`^(?:@?%s[:,]?)\s+(?:${1})`, b.Config.Name)
   	for _, cbInterface := range b.cbIndex[M]{
		callback := cbInterface.(messageCallback)
      var r *regexp.Regexp
      if callback.Method == `RESPOND`{
         r = regexp.MustCompile(strings.Replace(botNamePat,"${1}", callback.Val, 1))
      }else{
         r = regexp.MustCompile(callback.Pattern)
      }
      if r.MatchString(message.Text){
         match:=r.FindAllStringSubmatch(message.Text, -1)[0]
         Logger.Debug(`Broker:: running callback: `, callback.Name)
         callback.Chan <- struct{ Message: message, Match: match }
      }
   }
}

func (b *Broker) handleEvent(thingy map[string]interface{}){
	if b.eventCallbacks == nil { return }
	for _,cbInterface := range b.cbIndex[`events`]{
		callback := cbInterface.(eventCallback)
		if keyVal, keyExists := thingy[callback.Key]; keyExists && replyVal != nil{
      	if matches,_ := regex.MatchString(callback.Val, keyVal); matches{
				callback.Chan <- thingy
			}
		}
	}
}

func (b *Broker) MessageCallback(pattern string, method string) (MessageCallback){
	callback := &messageCallback {
		ID:			fmt.Sprintf("message:%s",len(b.cbIndex[M])),
		Pattern: 	pattern,
		Method: 		method,
		Chan:			new(chan struct{Message: string, Match: string}),
	}

	if err := RegisterCallback(callback); err != nil{
		Logger.Debug("error registering callback ", callback.ID, ":: ",err)
		return nil
	}
	return callback
}
	
func (b *Broker) EventCallback(key string, val string) EventCallback{
	callback := &EventCallback{
		ID:			fmt.Sprintf("event:%s",len(b.cbIndex[E])),
		Key: 			key,
		Val: 			val,
		Chan:			new(chan map[string]interface{}),
	}
	if err := RegisterCallback(callback); err != nil{
		Logger.Debug("error registering callback ", callback.ID, ":: ",err)
		return nil
	}
	return callback
}

func (b *Broker) TimerCallback(thingy map[string]interface{}){
	callback := &TimerCallback{
		ID:			fmt.Sprintf("event:%s",len(b.cbIndex[E])),
		Schedule: 	schedule,
		Chan:			new(chan time.Time)
	}
	if err := RegisterCallback(callback); err != nil{
		Logger.Debug("error registering callback ", callback.ID, ":: ",err)
		return nil
	}
	return callback
}
//func (b *Broker) LinkCallback(thingy map[string]interface{}){

