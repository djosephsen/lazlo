package slackerlib

import (
   "fmt"
   "time"
   "github.com/gorhill/cronexpr"
)

type Chore struct {
   Name     string
   Usage    string
   Sched    string
   Run      func(b *Sbot)
   State    string
   Next     time.Time
   Timer    *time.Timer
}

//an abstraction for the benifit of time.AfterFunc's second argument
// (probably just me being dense)
type ChoreTrigger struct {
	Chore		*Chore
	Sbot		*Sbot
}

func (t *ChoreTrigger) Pull(){
   Logger.Debug("Chore Triggered: ",t.Chore.Name)
   t.Chore.State="running"
   go t.Chore.Run(t.Sbot)
   t.Chore.Start(t.Sbot)
}

// Schedule the chores
func (bot *Sbot) StartChores() error{
   for _, c := range bot.Chores {
      c.Start(bot)
   }
   return nil
}

func (c *Chore) Start(bot *Sbot) error{
   expr := cronexpr.MustParse(c.Sched)
   if expr.Next(time.Now()).IsZero(){
      Logger.Debug("invalid schedule",c.Sched)
      c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Sched)
   }else{
      c.Next = expr.Next(time.Now())
      dur := c.Next.Sub(time.Now())
         if dur>0{
            if c.Timer == nil{
   				Logger.Debug("Chore: ",c.Name, " Instantiated. Scheduled at: ",c.Next)
					trigger:=&ChoreTrigger{
						Chore: c,
						Sbot:   bot,
					}
               c.Timer = time.AfterFunc(dur, trigger.Pull) // auto go-routine'd
            }else{
               c.Timer.Reset(dur) // auto go-routine'd
   				Logger.Debug("Chore: ",c.Name, " rescheduled at: ",c.Next)
            }
         c.State=fmt.Sprintf("Scheduled: %s",c.Next.String())
         }else{
            Logger.Debug("invalid duration",dur)
            c.State=fmt.Sprintf("Halted. (invalid duration: %s)",dur)
         }
      }
   return nil
}

func GetChoreByName(name string, bot *Sbot) *Chore{
   for _, c := range bot.Chores {
      if c.Name == name{
         return c
      }else{
         Logger.Debug("chore not found: ",name)
      }
   }
   return nil
}

func StopChore(c *Chore) error{
   Logger.Debug(`Stopping: `,c.Name)
	if c.Timer == nil{
		Logger.Error(`hmm.. can't stop job: `, c.Name,`. Timer undefined!`)
		return fmt.Errorf(`I GOTS NO TIMER!`)	
	}
   c.State=`Halted by request`
   c.Timer.Stop()
   return nil
}
