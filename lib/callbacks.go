package lib

import(
)

type MessageCallback struct{
   ID             string
   Pattern        string
   Respond        bool // if true, only respond if the bot is mentioned by name
   Chan           chan struct{Message: string, Match: string}
}

type EventCallback struct{
   ID           string
   Key            string
   Val            string
   Chan           chan map[string]interface{}
}
type TimerCallback struct{
   ID           string
   Schedule     string
   Chan         chan time.Time
	State			 string
	Next			 time.Time
}
type LinkCallback struct{
   ID           string
   Chan         map[string]interface{}
}

func (b *Broker) RegisterCallback(callback *interface{}) error{
	switch callback.(type){
		case MessageCallback:
			m:=callback.(MessageCallback)
			b.cbIndex[M][m.ID] = &callback
		case EventCallback:
			e:=callback.(EventCallback)
			b.cbIndex[E][e.ID] = &callback
		case TimerCallback:
			t:=callback.(TimerCallback)
			b.cbIndex[T][t.ID] = &callback
		case LinkCallback:
			l:=callback.(LinkCallback)
			b.cbIndex[L][l.ID] = &callback
		default:
			err:=fmt.Errorf("unknown type in register callback: %T",cbObj.(type))
			Logger.Error(err)
			return err
		}
return nil
}

func (b *Broker) MessageCallback(pattern string, method string) (MessageCallback){
   callback := &messageCallback {
      ID:         fmt.Sprintf("message:%s",len(b.cbIndex[M])),
      Pattern:    pattern,
      Method:     method,
      Chan:       new(chan struct{Message: string, Match: string}),
   }

   if err := RegisterCallback(callback); err != nil{
      Logger.Debug("error registering callback ", callback.ID, ":: ",err)
      return nil
   }
   return callback
}

func (b *Broker) EventCallback(key string, val string) EventCallback{
   callback := &EventCallback{
      ID:         fmt.Sprintf("event:%s",len(b.cbIndex[E])),
      Key:        key,
      Val:        val,
      Chan:       new(chan map[string]interface{}),
   }

   if err := callback.Start(); err != nil{
      Logger.Debug("error registering callback ", callback.ID, ":: ",err)
      return nil
	}

   if err := RegisterCallback(callback); err != nil{
      Logger.Debug("error registering callback ", callback.ID, ":: ",err)
      return nil
   }

   return callback
}

func (b *Broker) TimerCallback(schedule string) TimerCallback{
   callback := &TimerCallback{
      ID:         fmt.Sprintf("timer:%s",len(b.cbIndex[E])),
      Schedule:   schedule,
      Chan:       new(chan time.Time)
   }
   if err := RegisterCallback(callback); err != nil{
      Logger.Debug("error registering callback ", callback.ID, ":: ",err)
      return nil
   }
   return callback
}
//func (b *Broker) LinkCallback(thingy map[string]interface{}){

