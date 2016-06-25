package modules

import (
	"regexp"

	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/peterhellberg/giphy"
)

var Giphy = &lazlo.Module{
	Name:  "Giphy",
	Usage: `"gif [search|translate] <term>": Get a gif from giphy.com. Need to set ENV variables documented @github.com/peterhellberg/giphy`,
	Run:   getGIF,
}

func getGIF(b *lazlo.Broker) {
	g := giphy.DefaultClient
	str := []string{}
	cb := b.MessageCallback(`(?i)gif ((?i)search|translate) (\w+)`, false)
	for {
		pm := <-cb.Chan
		cmd := pm.Match[1]
		lazlo.Logger.Debug(cmd)
		if matched, _ := regexp.MatchString(`(?i)search`, cmd); matched {
			str := append(str, pm.Match[2])
			lazlo.Logger.Debug(str)
			resp, err := g.Search(str)
			if err != nil {
				lazlo.Logger.Error(err)
			}
			pm.Event.Reply(resp.Data[0].URL)
		} else if matched, _ := regexp.MatchString(`(?i)translate`, cmd); matched {
			str := append(str, pm.Match[2])
			lazlo.Logger.Debug(str)
			resp, err := g.Translate(str)
			if err != nil {
				lazlo.Logger.Error(err)
			}
			pm.Event.Reply(resp.Data.URL)
		}
	}
}
