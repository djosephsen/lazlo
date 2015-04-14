package modules

import(
	lazlo "github.com/djosephsen/lazlo/lib"
	"math/rand"
	"time"
)

var Syn = &lazlo.Module{
	Name:	`Ping`,
	Usage: `"%BOTNAME% (ping|syn)" : Test that the bot is currently running`,
	Run:	 func (b *lazlo.Broker){
		cb := b.MessageCallback(`(?i)(ping|syn)`, true)
		for {
			pm := <- cb.Chan 
			pm.Event.Reply(randReply())
		}
	},
}

func randReply() string{
	now:=time.Now()
	rand.Seed(int64(now.Unix()))
	replies := []string{
		"yeah um.. pong?",
		"WHAT?! jeeze.",
		"what? oh, um SYNACKSYN? ENOSPEAKTCP.",
		"RST (lulz)",
		"64 bytes from go.away.your.annoying icmp_seq=0 ttl=42 time=42.596 ms",
		"hmm?",
		"ack. what?",
		"pong. what?",
		"yup. still here.",
		"super busy just now.. Can I get back to you in like 5min?",
	}
	return replies[rand.Intn(len(replies)-1)]
}
