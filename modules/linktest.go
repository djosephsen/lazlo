package modules

import (
	"fmt"
	lazlo "github.com/djosephsen/lazlo/lib"
)

var LinkTest = &lazlo.Module{
	Name:  `LinkTest`,
	Usage: `"%BOTNAME% linkme foo" : creates a clickable link at servername/foo`,
	Run: func(b *lazlo.Broker) {
		clickChan := make(chan string)
		optionChan := make(chan string)
		command_cb := b.MessageCallback(`(?i)(link *me) (.*)`, true)
		command_cb1 := b.MessageCallback(`(?i)(link *test)`, true)
		command_cb2 := b.MessageCallback(`(?i)(link *choice)`, true)
		for {
			select {
			case msg := <-command_cb.Chan:
				msg.Event.Reply(newLink(b, msg.Match[2], clickChan))
			case msg := <-command_cb1.Chan:
				msg.Event.Reply(`<http://www.google.com|foo>`)
			case msg := <-command_cb2.Chan:
				msg.Event.Reply(newChoice(b, optionChan))
			case click := <-clickChan:
				b.Say(fmt.Sprintf("Somebody clicked on %s", click))
			case option := <-optionChan:
				if option == `THIS` {
					b.Say(fmt.Sprintf("I knew you'd get with this.. cause this is kinda phat"))
				} else {
					b.Say(fmt.Sprintf("Not a Blacksheep fan eh? bummer."))
				}
			}
		}
	},
}

func newLink(b *lazlo.Broker, path string, clickChan chan string) string {
	link_cb := b.LinkCallback(path)
	go func(link_cb *lazlo.LinkCallback, clickChan chan string) {
		for {
			<-link_cb.Chan
			clickChan <- link_cb.Path
		}
	}(link_cb, clickChan)
	return fmt.Sprintf("Ok, <%s|here> is a link on %s", link_cb.URL, path)
}

func newChoice(b *lazlo.Broker, clickChan chan string) string {
	opt1 := b.LinkCallback(`option1`)
	opt2 := b.LinkCallback(`option2`)
	go func(opt1 *lazlo.LinkCallback, opt2 *lazlo.LinkCallback, clickChan chan string) {
		for {
			select {
			case <-opt1.Chan:
				clickChan <- `THIS`
			case <-opt2.Chan:
				clickChan <- `THAT`
			}
		}
	}(opt1, opt2, clickChan)
	return fmt.Sprintf("you can get with <%s|THIS> or you can get with <%s|THAT>", opt1.URL, opt2.URL)
}
