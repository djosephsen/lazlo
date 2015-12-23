package main

import (
	lazlo "github.com/klaidliadon/lazlo/lib"
	"github.com/klaidliadon/lazlo/modules"
)

func initModules(b *lazlo.Broker) error {
	b.Register(modules.Syn)
	//	b.Register(modules.RTMPing)
	b.Register(modules.LinkTest)
	b.Register(modules.BrainTest)
	b.Register(modules.Help)
	b.Register(modules.LuaMod)
	b.Register(modules.QuestionTest)
	return nil
}
