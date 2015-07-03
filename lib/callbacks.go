package lib

import (
	"fmt"
	"time"
)

type MessageCallback struct {
	ID      string
	Pattern string
	Respond bool // if true, only respond if the bot is mentioned by name
	Chan    chan PatternMatch
	SlackChan string // if set filter message callbacks to this Slack channel
}

type PatternMatch struct {
	Event *Event
	Match []string
}

type EventCallback struct {
	ID   string
	Key  string
	Val  string
	Chan chan map[string]interface{}
}

type TimerCallback struct {
	ID       string
	Schedule string
	State    string
	Next     time.Time
	Chan     chan time.Time
	stop     chan bool //send true on this channel to stop the timer
}

type QuestionCallback struct{
	ID       string
	User     string
	DMChan	string
	Question string
	Answer	chan string
	asked    bool
}

type QuestionQueue struct{
	in     chan *QuestionCallback
}

func (b *Broker) RegisterCallback(callback interface{}) error {
	switch callback.(type) {
	case *MessageCallback:
		m := callback.(*MessageCallback)
		b.cbIndex[M][m.ID] = callback
		Logger.Debug("New Callback Registered, id:", m.ID)
	case *EventCallback:
		e := callback.(*EventCallback)
		b.cbIndex[E][e.ID] = callback
		Logger.Debug("New Callback Registered, id:", e.ID)
	case *TimerCallback:
		t := callback.(*TimerCallback)
		t.Start()
		b.cbIndex[T][t.ID] = callback
		Logger.Debug("New Callback Registered, id:", t.ID)
	case *LinkCallback:
		l := callback.(*LinkCallback)
		b.cbIndex[L][l.ID] = callback
		Logger.Debug("New Callback Registered, id:", l.ID)
	case *QuestionCallback:
		q := callback.(*QuestionCallback)
		b.cbIndex[Q][q.ID] = callback
		Logger.Debug("New Callback Registered, id:", q.ID)
	default:
		err := fmt.Errorf("unknown type in register callback: %T", callback)
		Logger.Error(err)
		return err
	}
	return nil
}

func (b *Broker) DeRegisterCallback(callback interface{}) error {
	switch callback.(type) {
	case *MessageCallback:
		m := callback.(*MessageCallback)
		delete(b.cbIndex[M], m.ID)
		Logger.Debug("De-Registered callback, id: ", m.ID)
	case *EventCallback:
		e := callback.(*EventCallback)
		delete(b.cbIndex[E], e.ID)
		Logger.Debug("De-Registered callback, id: ", e.ID)
	case *TimerCallback:
		t := callback.(*TimerCallback)
		t.Stop() // dont leak timers
		delete(b.cbIndex[T], t.ID)
		Logger.Debug("De-Registered callback, id: ", t.ID)
	case *LinkCallback:
		l := callback.(*LinkCallback)
		l.Delete() //dont leak httproutes
		delete(b.cbIndex[L], l.ID)
		Logger.Debug("De-Registered callback, id: ", l.ID)
	case *QuestionCallback:
		q := callback.(*QuestionCallback)
		delete(b.cbIndex[Q], q.ID)
		Logger.Debug("De-Registered callback, id: ", q.ID)
	default:
		err := fmt.Errorf("unknown type in de-register callback: %T", callback)
		Logger.Error(err)
		return err
	}
	return nil
}

func (b *Broker) MessageCallback(pattern string, respond bool, channel ...string) *MessageCallback {
	callback := &MessageCallback{
		ID:      fmt.Sprintf("message:%d", len(b.cbIndex[M])),
		Pattern: pattern,
		Respond: respond,
		Chan:    make(chan PatternMatch),
	}

	if channel != nil{
		callback.SlackChan = channel[0] // todo: support an array of channels
	}

	if err := b.RegisterCallback(callback); err != nil {
		Logger.Debug("error registering callback ", callback.ID, ":: ", err)
		return nil
	}
	return callback
}

func (b *Broker) EventCallback(key string, val string) *EventCallback {
	callback := &EventCallback{
		ID:   fmt.Sprintf("event:%d", len(b.cbIndex[E])),
		Key:  key,
		Val:  val,
		Chan: make(chan map[string]interface{}),
	}

	if err := b.RegisterCallback(callback); err != nil {
		Logger.Debug("error registering callback ", callback.ID, ":: ", err)
		return nil
	}
	return callback
}

func (b *Broker) TimerCallback(schedule string) *TimerCallback {
	callback := &TimerCallback{
		ID:       fmt.Sprintf("timer:%d", len(b.cbIndex[E])),
		Schedule: schedule,
		Chan:     make(chan time.Time),
	}
	if err := b.RegisterCallback(callback); err != nil {
		Logger.Debug("error registering callback ", callback.ID, ":: ", err)
		return nil
	}
	return callback
}

func (b *Broker) QuestionCallback(user string, prompt string) *QuestionCallback {
	callback := &QuestionCallback{
		ID:      fmt.Sprintf("question:%d", len(b.cbIndex[Q])),
		User:    user,
		Question: prompt,
		Answer:   make(chan string),
	}
	if err := b.RegisterCallback(callback); err != nil {
		Logger.Debug("error registering callback ", callback.ID, ":: ", err)
		return nil
	}
	return callback
}


// LinkCallback() def is in httpserver.go because it includes net/http (sorry)
