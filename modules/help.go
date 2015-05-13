package modules

import (
	"fmt"
	"strings"
	lazlo "github.com/djosephsen/lazlo/lib"
)

var Help = lazlo.Module{
	Name:    `Help`,
	Usage:   `%BOTNAME% help: prints the usage information of every registered plugin`,
	Run: 		helpRun, 
}

func helpRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)help`, true)
	for {
		pm := <-cb.Chan
		go getHelp(b, &pm)
	}
}

func getHelp(b *lazlo.Broker, pm *lazlo.PatternMatch){
	dmChan := b.GetDM(pm.Event.User)
	reply:=`########## Modules In use: `
	for _, m := range b.Modules {
		if strings.Contains(m.Usage,`%HIDDEN%`){continue}
		usage := strings.Replace(m.Usage,`%BOTNAME%`,b.Config.Name,-1)
		reply = fmt.Sprintf("%s\n%s",reply,usage)
	}
	b.Say(reply,dmChan)
}
