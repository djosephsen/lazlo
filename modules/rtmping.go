package modules

import (
	lazlo "github.com/klaidliadon/lazlo/lib"
)

var RTMPing = &lazlo.Module{
	Name:  `RTMPing`,
	Usage: `Automatially sends an RTM Ping to SlackHQ every 20 seconds`,
	Run:   rtmrun,
}

func rtmrun(b *lazlo.Broker) {
	for {
		// get a timer callback
		timer := b.TimerCallback(`*/20 * * * * * *`)

		// block waiting for an alarm from the timer
		<-timer.Chan

		//send a ping
		b.Send(&lazlo.Event{
			Type: `ping`,
			Text: `just pingin`,
		})
	}
}
