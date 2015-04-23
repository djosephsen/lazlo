package modules

import(
	lazlo "github.com/djosephsen/lazlo/lib"
	lua "github.com/yuin/gopher-lua"
	"math/rand"
	"time"
)

//luaMod implements a lua-parsing plugin for Lazlo.
//this enables lazlo to be scripted via lua instead of GO, which is 
//preferable in some contexts (no recompiles for changes etc..).
var LuaMod = &lazlo.Module{
	Name:	`LuaMod`, 
	Usage: `%HIDDEN% this module implements lua scripting of lazlo`
	Run:	 luaMain, 
}

//LuaStates is a lookup table for each of the lua scripts
var LuaStates []*lua.LState
//CBtable is a lookup table for callbacks we register with lazlo
var CBTable []
//Cases is a lookup table that allows us to select on a dynamic number of callback Channels
var Cases []SelectCase

//luaMain reads in lua files, registers their callbacks with lazlo
//and brokers back events that occur in the chatroom
func luaMain (b *lazlo.Broker){
    L := lua.NewState()
    defer L.Close()

	 //register functions here
    L.SetGlobal("double", L.NewFunction(Double)) /* Original lua_setglobal uses stack... */

	//block waiting on events from the broker
	for{
		index, value, _ := reflect.Select(Cases)
		handle(index, value.Interface())	
	}
    L.Push(lua.LNumber(lv * 2)) /* push result */
}

//handle takes takes the index and value of an event from lazlo, 
//typifies the value and calls the right function to push the data back
//to whatever lua script asked for it.
func handle(index int, val interface{}){
	switch val.(type){
	case lazlo.PatternMatch:
		handleMessageCB(index, val.(lazlo.PatternMatch))	
	case time.Time:
		handleTimerCB(index, val.(time.Time))	
	case map[string]interface{}:
		handleEventCB(index, val.(map[string]interface{}))	
	case *http.Request:
		handleLinkCB(index, val.(*http.Request))	
	default
		err:=fmt.Sprintf("luaMod handle:: unknown type: %T",val)
		lazlo.Logger.Error(err)
	}
}

//handleMessageCB brokers messages back to the lua script that asked for them
func handleMessageCB(index int, message lazlo.PatternMatch){
	l:=LuaStates[index]
	lmsg:=l.NewTable()
	return
}

//handleTimerCB brokers timer alarms back to the lua script that asked for them
func handleTimerCB(index int, t time.TIME){
	return
}

//handleEventCB brokers slack rtm events back to the lua script that asked for them
func handleEventCB(index int, event map[string]interface{}){
	return
}

//handleLinkCB brokers http GET requests back to the lua script that asked for them
func handleLinkCB(index int, resp *http.Response())
	return
}

//functions exported to the lua runtime below here

func Hear(b* lazlo.Broker, L *lua.LState){
    pat := L.ToString(-1)             /* get argument */
	 cb := b.MessageCallback(pat,1)
    return 1                     /* number of results */
}

func Respond(L *lua.LState) int {
    pat := L.ToString(-1)             /* get argument */

