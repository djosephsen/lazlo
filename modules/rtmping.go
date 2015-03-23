package modules

import(
	lazlo "github.com/djosephsen/lazlo/lib"
	"time"
	"rand"
)

RTMPing := &lazlo.Module{
	Name:	`RTMPing`,
	Usage: "Automatially sends an RTM Ping to SlackHQ every 20 seconds",
	Run:	 run,
	SyncChan: make(chan bool),
}

func run(b *lazlo.Broker){
	timer := b.TimerCallback(`*/20 * * * * * *`)
	stop := false
	for !stop{
		select{
		case event := <- cb.Chan	
			b.Send(&lazlo.Event{
         	Type: `ping`,
         	Text: `just pingin`,
      	})
		case stop = <- cb.SyncChan
			stop = true
		}
	}
	b.SyncChan <- true
}

