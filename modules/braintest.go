package modules

import (
	"fmt"
	"regexp"

	lazlo "github.com/klaidliadon/lazlo/lib"
)

var BrainTest = &lazlo.Module{
	Name:  `BrainTest`,
	Usage: `"%BOTNAME% brain [set|get] <key> <value>": tests lazlo's persistent storage (aka the brain)`,
	Run: func(b *lazlo.Broker) {
		callback := b.MessageCallback(`(?i:brain) ((?i)set|get) (\w+) *(\w*)$`, true)
		for {
			msg := <-callback.Chan
			brain := b.Brain
			cmd := msg.Match[1]
			key := msg.Match[2]
			if matched, _ := regexp.MatchString(`(?i)set`, cmd); matched {
				val := msg.Match[3]
				if err := brain.Set(key, []byte(val)); err != nil {
					msg.Event.Reply(fmt.Sprintf("Sorry, something went wrong: %s", err))
					lazlo.Logger.Error(err)
				} else {
					msg.Event.Reply(fmt.Sprintf("Ok, %s set to %s", key, val))
				}
			} else {
				val, err := brain.Get(key)
				if err != nil {
					msg.Event.Reply(fmt.Sprintf("Sorry, something went wrong: %s", err))
				} else {
					msg.Event.Reply(string(val))
				}
			}
		}
	},
}
