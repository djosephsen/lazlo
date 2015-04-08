package lib

import (
   "fmt"
   "time"
   "github.com/gorhill/cronexpr"
)

// verify the schedule and start the timer
func (t *TimerCallback) Start() error{
   expr := cronexpr.MustParse(t.Schedule)
   if expr.Next(time.Now()).IsZero(){
      Logger.Debug("invalid schedule",t.Schedule)
      t.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",t.Schedule)
		return fmt.Errorf("invalid schedule",t.Schedule)
	}
   t.Next = expr.Next(time.Now())
   dur := t.Next.Sub(time.Now())
   if dur>0{
		go t.Run(dur)
	}
   return nil
}

// wait for the timer to expire, callback to the module, and reschedule
func (t *TimerCallback) Run(dur time.Duration){
	Logger.Debug(`scheduling timer `, t.ID, ` for: `,t.Next)
	timer := time.NewTimer(dur)
	alarm := <- timer.C //blocks waiting for the timer
   t.Chan <- alarm //signals the module
	t.Start() // (potentially) reschedule
}

