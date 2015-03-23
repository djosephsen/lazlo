package main

import(
	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/djosephsen/lazlo/modules"
)

func initModules(b *lib.Broker) error{
	b.Register(modules.Ping)
	b.Register(modules.RTMPing)
	return nil
}
