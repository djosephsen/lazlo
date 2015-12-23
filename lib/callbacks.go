package lib

import (
	"fmt"
	"time"
)

type MessageCallback struct {
	ID        string
	Pattern   string
	Respond   bool // if true, only respond if the bot is mentioned by name
	Chan      chan PatternMatch
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

type QuestionCallback struct {
	ID       string
	User     string
	DMChan   string
	Question string
	Answer   chan string
	asked    bool
}

type QuestionQueue struct {
	in chan *QuestionCallback
}

func (b *Broker) RegisterCallback(callback interface{}) error {
	switch callback.(type) {
	case *MessageCallback:
		m := callback.(*MessageCallback)
		b.callbacks[M].Lock()
		defer b.callbacks[M].Unlock()
		b.callbacks[M].Index[m.ID] = callback
		Logger.Debug("New Callback Registered, id:", m.ID, " *MessageCallback")
	case *EventCallback:
		e := callback.(*EventCallback)
		b.callbacks[E].Lock()
		defer b.callbacks[E].Unlock()
		b.callbacks[E].Index[e.ID] = callback
		Logger.Debug("New Callback Registered, id:", e.ID, " *EventCallback")
	case *TimerCallback:
		t := callback.(*TimerCallback)
		t.Start()
		b.callbacks[T].Lock()
		defer b.callbacks[T].Unlock()
		b.callbacks[T].Index[t.ID] = callback
		Logger.Debug("New Callback Registered, id:", t.ID, ":= callback.( *TimerCallback")
	case *LinkCallback:
		l := callback.(*LinkCallback)
		b.callbacks[L].Lock()
		defer b.callbacks[L].Unlock()
		b.callbacks[L].Index[l.ID] = callback
		Logger.Debug("New Callback Registered, id:", l.ID, " *LinkCallback")
	case *QuestionCallback:
		q := callback.(*QuestionCallback)
		b.callbacks[Q].Lock()
		defer b.callbacks[Q].Unlock()
		b.callbacks[Q].Index[q.ID] = callback
		Logger.Debug("New Callback Registered, id:", q.ID, " *QuestionCallback")
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
		b.callbacks[M].Lock()
		defer b.callbacks[M].Unlock()
		delete(b.callbacks[M].Index, m.ID)
		Logger.Debug("De-Registered callback, id: ", m.ID)
	case *EventCallback:
		e := callback.(*EventCallback)
		b.callbacks[E].Lock()
		defer b.callbacks[E].Unlock()
		delete(b.callbacks[E].Index, e.ID)
		Logger.Debug("De-Registered callback, id: ", e.ID)
	case *TimerCallback:
		t := callback.(*TimerCallback)
		t.Stop() // dont leak timers
		b.callbacks[T].Lock()
		defer b.callbacks[T].Unlock()
		delete(b.callbacks[T].Index, t.ID)
		Logger.Debug("De-Registered callback, id: ", t.ID)
	case *LinkCallback:
		l := callback.(*LinkCallback)
		l.Delete() //dont leak httproutes
		b.callbacks[L].Lock()
		defer b.callbacks[L].Unlock()
		delete(b.callbacks[L].Index, l.ID)
		Logger.Debug("De-Registered callback, id: ", l.ID)
	case *QuestionCallback:
		q := callback.(*QuestionCallback)
		delete(b.callbacks[Q].Index, q.ID)
		Logger.Debug("De-Registered callback, id: ", q.ID)
	default:
		err := fmt.Errorf("unknown type in de-register callback: %T", callback)
		Logger.Error(err)
		return err
	}
	return nil
}

func (b *Broker) MessageCallback(pattern string, respond bool, channel ...string) *MessageCallback {
	b.callbacks[M].Lock()
	l := len(b.callbacks[M].Index)
	b.callbacks[M].Unlock()
	callback := &MessageCallback{
		ID:      fmt.Sprintf("message:%d", l),
		Pattern: pattern,
		Respond: respond,
		Chan:    make(chan PatternMatch),
	}

	if channel != nil {
		callback.SlackChan = channel[0] // todo: support an array of channels
	}

	if err := b.RegisterCallback(callback); err != nil {
		Logger.Debug("error registering callback ", callback.ID, ":: ", err)
		return nil
	}
	return callback
}

func (b *Broker) EventCallback(key string, val string) *EventCallback {
	b.callbacks[E].Lock()
	l := len(b.callbacks[E].Index)
	b.callbacks[E].Unlock()
	callback := &EventCallback{
		ID:   fmt.Sprintf("event:%d", l),
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
	b.callbacks[E].Lock()
	l := len(b.callbacks[E].Index)
	b.callbacks[E].Unlock()
	callback := &TimerCallback{
		ID:       fmt.Sprintf("timer:%d", l),
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
	b.callbacks[Q].Lock()
	l := len(b.callbacks[Q].Index)
	b.callbacks[Q].Unlock()
	callback := &QuestionCallback{
		ID:       fmt.Sprintf("question:%d", l),
		User:     user,
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
