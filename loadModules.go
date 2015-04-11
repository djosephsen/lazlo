package main

import(
	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/djosephsen/lazlo/modules"
)

func initModules(b *lazlo.Broker) error{
	b.Register(modules.Syn)
	b.Register(modules.RTMPing)
	b.Register(modules.LinkTest)
	return nil
}
