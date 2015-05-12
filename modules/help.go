package modules

import (
	"fmt"
	"string"
	lazlo "github.com/djosephsen/lazlo/lib"
)

var Help = lazlo.MessageHandler{
	Name:    `Help`,
	Usage:   `<botname> help: prints the usage information of every registered plugin`,
	Run: 		helpRun, 
}

func helpRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)help`, true)
	for {
		pm := <-cb.Chan
		go getHelp(b, pm)
	}
}

func getHelp(b *lazlo.Broker, pm *lazlo.PatternMatch){
	dmChan := b.GetDM(pm.User)
	reply:=``
	for _, m := range b.Modules {
		if strings.Contains(m.Usage,`%HIDDEN%`){continue}
		usage:=strings.Replace(m.Usage,`%BOTNAME%`,b.Config.Name,-1)


	}
	pm.Event.Reply(randReply())
}

		if len(e.Sbot.Broker.MessageHandlers) > 0 {
			line := fmt.Sprintf("######## Message Handlers ##########\n")
				line += fmt.Sprintf("*%s*:: %s\n", h.Name, h.Usage)
			}
			e.Respond(line)
		}
		if len(e.Sbot.Broker.EventHandlers) > 0 {
			line := fmt.Sprintf("######## Event Handlers ##########\n")
			for _, h := range e.Sbot.Broker.EventHandlers {
				line += fmt.Sprintf("*%s*:: %s\n", h.Name, h.Usage)
				e.Respond(line)
			}
		}
		if len(e.Sbot.Chores) > 0 {
			line := fmt.Sprintf("######## Chores ##########\n")
			for _, h := range e.Sbot.Chores {
				line += fmt.Sprintf("*%s* (%s):: %s\n", h.Name, h.Sched, h.Usage)
				e.Respond(line)
			}
		}
		if len(e.Sbot.StartupHooks) > 0 {
			line := fmt.Sprintf("######## Startup Hooks ##########\n")
			for _, h := range e.Sbot.StartupHooks {
				line += fmt.Sprintf("*%s*:: %s\n", h.Name, h.Usage)
				e.Respond(line)
			}
		}
		if len(e.Sbot.ShutdownHooks) > 0 {
			line := fmt.Sprintf("######## Shutdown Hooks ##########\n")
			for _, h := range e.Sbot.ShutdownHooks {
				line += fmt.Sprintf("*%s*:: %s\n", h.Name, h.Usage)
				e.Respond(line)
			}
		}
		if len(e.Sbot.Broker.PreFilters) > 0 {
			line := fmt.Sprintf("######## Input Filters ##########\n")
			for _, h := range e.Sbot.Broker.PreFilters {
				line += fmt.Sprintf("*%s*:: %s\n", h.Name, h.Usage)
				e.Respond(line)
			}
		}
		if len(e.Sbot.WriteThread.OutputFilters) > 0 {
			line := fmt.Sprintf("######## Output Filters ##########\n")
			for _, h := range e.Sbot.WriteThread.OutputFilters {
				line += fmt.Sprintf("*%s*:: %s\n", h.Name, h.Usage)
				e.Respond(line)
			}
		}
	},
}
