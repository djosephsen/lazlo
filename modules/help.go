package modules

import (
	"fmt"
	lazlo "github.com/djosephsen/lazlo/lib"
)

var Help = lazlo.MessageHandler{
	Name:    `Help`,
	Usage:   `<botname> help: prints the usage information of every registered plugin`,
	Run: 		helpRun, 
}

func helpRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)(ping|syn)`, true)
		for {
			pm := <-cb.Chan
			pm.Event.Reply(randReply())
		}
}

		if len(e.Sbot.Broker.MessageHandlers) > 0 {
			line := fmt.Sprintf("######## Message Handlers ##########\n")
			for _, h := range e.Sbot.Broker.MessageHandlers {
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
