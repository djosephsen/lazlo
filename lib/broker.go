package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ccding/go-logging/logging"
	"github.com/gorilla/websocket"
	"os"
	"regexp"
	"strings"
	"time"
	"net/url"
)

// Logger is a global reference to our logging object
var Logger = newLogger()

// These contstats define the four types of callbacks that lazlo can hand you
const M = "messages"
const E = "events"
const T = "timers"
const L = "links"

// Broker is the all-knowing repository of references
type Broker struct {
	SlackMeta    *ApiResponse
	Config       *Config
	Socket       *websocket.Conn
	Modules      map[string]*Module
	Brain        Brain
	ApiResponses map[int32]chan map[string]interface{}
	cbIndex      map[string]map[string]interface{} //cbIndex[type][id]=pointer
	ReadFilters  []*ReadFilter
	WriteFilters []*WriteFilter
	MID          int32
	WriteThread  *WriteThread
	SigChan      chan os.Signal
	SyncChan     chan bool
	ThreadCount  int32
}

// The Module type represents a user-defined plug-in. Build one of these
// and add it to loadModules.go for Lazlo to run your thingy on startup
type Module struct {
	Name  string
	Usage string
	Run   func(*Broker)
}

// The WriteThread serielizes and sends messages to the slack RTM interface
type WriteThread struct {
	broker   *Broker
	Chan     chan Event
	SyncChan chan bool
}

// ReadFilter is a yet-to-be-implemented hook run on all inbound
// events from slack before the broker gets a hold of them
type ReadFilter struct {
	Name  string
	Usage string
	Run   func(thingy map[string]interface{}) map[string]interface{}
}

// WriteFilter is a yet-to-be-implemented hook run on all outbound
// events from slack before the broker gets a hold of them
type WriteFilter struct {
	Name  string
	Usage string
	Run   func(e *Event)
}

// NewBroker instantiates a new broker
func NewBroker() (*Broker, error) {

	broker := &Broker{
		MID:          0,
		Config:       newConfig(),
		Modules:      make(map[string]*Module),
		ApiResponses: make(map[int32]chan map[string]interface{}),
		cbIndex:      make(map[string]map[string]interface{}),
		WriteThread: &WriteThread{
			Chan:     make(chan Event),
			SyncChan: make(chan bool),
		},
		SigChan:  make(chan os.Signal),
		SyncChan: make(chan bool),
	}
	//correctly set the log level
	Logger.SetLevel(logging.GetLevelValue(strings.ToUpper(broker.Config.LogLevel)))

	broker.cbIndex[M] = make(map[string]interface{})
	broker.cbIndex[E] = make(map[string]interface{})
	broker.cbIndex[T] = make(map[string]interface{})
	broker.cbIndex[L] = make(map[string]interface{})
	broker.WriteThread.broker = broker

	//connect to slack and establish an RTM websocket
	socket, meta, err := broker.getASocket()
	if err != nil {
		return nil, err
	}
	broker.Socket = socket
	broker.SlackMeta = meta

	broker.Brain, err = broker.newBrain()
	if err != nil {
		return nil, err
	}
	//	broker.Brain = brain
	if err = broker.Brain.Open(); err != nil {
		Logger.Error(`couldn't open mah brain! `, err)
		return broker, err
	}
	return broker, nil
}

// Stop gracefully stops lazlo
func (broker *Broker) Stop() {
	// make sure the write thread finishes before we stop
	broker.WriteThread.SyncChan <- true
}

// It's called Start(), I mean srsly.
func (broker *Broker) Start() {
	go broker.StartHttp()
	go broker.WriteThread.Start()
	Logger.Debug(`Broker:: entering read-loop`)
	for {
		thingy := make(map[string]interface{})
		broker.Socket.ReadJSON(&thingy)
		go broker.This(thingy)
	}
}

// StartModules launches each user-provided plugin registered in loadMOdules.go
func (b *Broker) StartModules() {
	for _, module := range b.Modules {
		go module.Run(b)
	}
}

// WriteThread.Start starts the writethread
func (w *WriteThread) Start() {
	Logger.Debug(`Write-Thread Started`)
	stop := false
	for !stop {
		select {
		case e := <-w.Chan:
			Logger.Debug(`WriteThread:: Outbound `, e.Type, ` channel: `, e.Channel, `. text: `, e.Text)
			ejson := stupidUTFHack(e)
			if len(ejson) >= 16000 {
				e = Event{
					ID:      e.ID,
					Type:    e.Type,
					Channel: e.Channel,
					Text:    fmt.Sprintf("ERROR! Response too large. %v Bytes!", len(ejson)),
				}
				ejson = stupidUTFHack(e)
			}
			if matches, _ := regexp.MatchString(`<[hH#@].+>`, string(ejson)); matches {
				Logger.Debug(`message formtting detected; sending via api`)
				e.Broker = w.broker
				apiPostMessage(e)
			} else {
				w.broker.Socket.WriteMessage(1, ejson)
			}
			Logger.Debug(string(ejson))
			time.Sleep(time.Second * 1)
		case stop = <-w.SyncChan:
			stop = true
		}
	}
	//signal main that we're done
	w.broker.SyncChan <- true
}

//This stupid hack un-does the utf-escaping performed  by json.Marshal()
//because although Slack correctly parses utf, it doesn't recognize
//utf-escaped markup like <http://myurl.com|myurl>
// UPDATE: I can remove this Once I re-figure-out out how the hell it works
func stupidUTFHack(thingy interface{}) []byte {
	jThingy, _ := json.Marshal(thingy)
	jThingy = bytes.Replace(jThingy, []byte("\\u003c"), []byte("<"), -1)
	jThingy = bytes.Replace(jThingy, []byte("\\u003e"), []byte(">"), -1)
	jThingy = bytes.Replace(jThingy, []byte("\\u0026"), []byte("&"), -1)
	return jThingy
}

//NextMID() ensures our outbound messages have a unique ID number
// (a requirement of the slack rtm api)
func (b *Broker) NextMID() int32 {
	//probably need to make this thread-safe (for now only the write thread uses it)
	b.MID += 1
	Logger.Debug(`incrementing MID to `, b.MID)
	return b.MID
}

func (b *Broker) This(thingy map[string]interface{}) {
	if b.Modules == nil {
		Logger.Debug(`Broker:: Got a `, thingy[`type`], ` , but no modules are loaded!`)
		return
	}
	//run the pre-handeler filters
	if b.ReadFilters != nil {
		for _, filter := range b.ReadFilters { //run the read filters
			thingy = filter.Run(thingy)
		}
	}
	// stop here if a prefilter delted our thingy
	if len(thingy) == 0 {
		return
	}

	Logger.Debug(`broker:: got a `, thingy[`type`])
	// if it's an api response send it to whomever is listening for it
	if replyVal, isReply := thingy[`reply_to`]; isReply {
		if replyVal != nil { // sometimes the api returns: "reply_to":null
			b.handleApiReply(thingy)
		}
	}

	typeOfThingy := thingy[`type`]
	switch typeOfThingy {
	case nil:
		return
	case `message`:
		b.handleMessage(thingy)
	default:
		b.handleEvent(thingy)
	}
}

func (b *Broker) Register(things ...interface{}) {
	// this is where we register user-provided plug-in code of various description
	for _, thing := range things {
		switch t := thing.(type) {
		case *Module:
			m := thing.(*Module)
			Logger.Debug(`registered Module: `, m.Name)
			b.Modules[m.Name] = m
		case *ReadFilter:
			r := thing.(*ReadFilter)
			Logger.Debug(`registered Read Filter: `, r.Name)
			b.ReadFilters = append(b.ReadFilters, r)
		case *WriteFilter:
			w := thing.(*WriteFilter)
			Logger.Debug(`registered Write Filter: `, w.Name)
			b.WriteFilters = append(b.WriteFilters, w)
		default:
			weirdType := fmt.Sprintf(`%T`, t)
			Logger.Error(`sorry I cant register this handler because I don't know what a `, weirdType, ` is`)
		}
	}
}

func (b *Broker) handleApiReply(thingy map[string]interface{}) {
	chanID := int32(thingy[`reply_to`].(float64))
	Logger.Debug(`Broker:: caught a reply to: `, chanID)
	if callBackChannel, exists := b.ApiResponses[chanID]; exists {
		callBackChannel <- thingy
		//dont leak channels
		Logger.Debug(`deleting callback: `, chanID)
		close(callBackChannel)
		<-callBackChannel
		delete(b.ApiResponses, chanID)
	} else {
		Logger.Debug(`no such channel: `, chanID)
	}
}

func (b *Broker) handleMessage(thingy map[string]interface{}) {
	if b.cbIndex[M] == nil {
		return
	}
	message := new(Event)
	jthingy, _ := json.Marshal(thingy)
	json.Unmarshal(jthingy, message)
	message.Broker = b
	botNamePat := fmt.Sprintf(`^(?:@?%s[:,]?)\s+(?:${1})`, b.Config.Name)
	for _, cbInterface := range b.cbIndex[M] {
		callback := cbInterface.(*MessageCallback)
		var r *regexp.Regexp
		if callback.Respond {
			r = regexp.MustCompile(strings.Replace(botNamePat, "${1}", callback.Pattern, 1))
		} else {
			r = regexp.MustCompile(callback.Pattern)
		}
		if r.MatchString(message.Text) {
			match := r.FindAllStringSubmatch(message.Text, -1)[0]
			Logger.Debug(`Broker:: running callback: `, callback.ID)
			callback.Chan <- PatternMatch{Event: message, Match: match}
		}
	}
}

func (b *Broker) handleEvent(thingy map[string]interface{}) {
	if b.cbIndex[E] == nil {
		return
	}
	for _, cbInterface := range b.cbIndex[E] {
		callback := cbInterface.(EventCallback)
		if keyVal, keyExists := thingy[callback.Key]; keyExists && keyVal != nil {
			if matches, _ := regexp.MatchString(callback.Val, keyVal.(string)); matches {
				callback.Chan <- thingy
			}
		}
	}
}

// this is the primary interface to Slack's write socket. Use this to send events.
func (b *Broker) Send(e *Event) chan map[string]interface{} {
	e.ID = b.NextMID()
	b.ApiResponses[e.ID] = make(chan map[string]interface{}, 1)
	Logger.Debug(`created APIResponse: `, e.ID)
	b.WriteThread.Chan <- *e
	return b.ApiResponses[e.ID]
}

// Say something in the named channel (or the default channel if none specified)
func (b *Broker) Say(s string, channel ...string) chan map[string]interface{} {
	var c string
	if channel != nil {
		c = channel[0]
	} else {
		c = b.DefaultChannel()
	}
	resp := b.Send(&Event{
		Type:    `message`,
		Channel: c,
		Text:    s,
	})
	return resp
}

// send a reply to any sort of thingy that contains an ID and Channel attribute
func (b *Broker) Respond(text string, thing *interface{}, isReply bool) chan map[string]interface{} {
	var id, channel string
	var exists bool

	thingy := *thing
	switch thingy.(type) {
	case Event:
		eThingy := thingy.(Event)
		if eThingy.User != `` && eThingy.Channel != `` {
			id = eThingy.User
			channel = eThingy.Channel
		} else {
			return nil
		}
	case map[string]interface{}:
		mThingy := thingy.(map[string]interface{})
		if id, exists = mThingy[`id`].(string); !exists || id == `` {
			return nil
		}
		if channel, exists = mThingy[`channel`].(string); !exists || channel == `` {
			return nil
		}
		id = mThingy[`id`].(string)
		channel = mThingy[`channel`].(string)
	default:
		return nil
	}

	var replyText string
	if isReply {
		replyText = fmt.Sprintf(`%s: %s`, b.SlackMeta.GetUserName(id), text)
	} else {
		replyText = text
	}

	return b.Send(&Event{
		Type:    `message`,
		Channel: channel,
		Text:    replyText,
	})
}

//Get a direct message channel ID so we can DM the given user
func (b *Broker) GetDM(ID string) string {
	req := ApiRequest{ //use the web api so we don't block waiting for the read thread
		URL:	`https://slack.com/api/im.open`,
		Values: make(url.Values),
		Broker: b,
	}
	reply, err := MakeAPIReq(req)
	if err != nil{
		Logger.Error(`error making api request for dm channel: `,err)
		return ``
	}else{
		return reply.Channel.ID
	}
}

//returns the Team's default channel
func (b *Broker) DefaultChannel() string {
	for _, c := range b.SlackMeta.Channels {
		if c.IsGeneral {
			return c.ID
		}
	}
	return b.SlackMeta.Channels[0].ID
}
