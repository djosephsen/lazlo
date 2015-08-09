package modules

import (
	lazlo "github.com/djosephsen/lazlo/lib"
	"math/rand"
	"time"
)

var Syn = &lazlo.Module{
	Name:  `Karma`,
	Usage: `"%BOTNAME% (+|-),<n> <user>" : add or subtract karma points from a user`,
	Run:   karmaRun,
}

func karmaRun(b *lazlo.Broker) {
	plusCB := b.MessageCallback(`(?i)\+([0-9]+) (\w+)`, true
	minusCB := b.MessageCallback(`(?i)\-([0-9]+) (\w+)`, true)

	for {
		select{
			case pm := <-plusCB.Chan:
				addKarma(pm.Match[1], pm.Match[2], pm.Event.User, b)
			case pm := <-minusCB.Chan:
				subtractKarma(pm.Match[1], pm.Match[2], pm.Event.User, b)
		}
	}
}

addKarma(points string, name string, requestor, string, b lazlo.Broker){
	if requestor == name{
		b.Say(fmt.Sprintf("sorry %s, ones karma must reflect ones actions", name)
	}
	brain:=b.Brain



subtractKarma(points string, name string, requestor, string, b lazlo.Broker){
	if requestor == name{
		b.Say(fmt.Sprintf("sorry %s, ones karma must reflect ones actions", name)
	}
