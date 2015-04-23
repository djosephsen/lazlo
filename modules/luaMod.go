package modules

import(
	lazlo "github.com/djosephsen/lazlo/lib"
	lua "github.com/yuin/gopher-lua"
	luar "github.com/layeh/gopher-luar"
	"time"
	"reflect"
	"net/http"
)

//luaMod implements a lua-parsing plugin for Lazlo.
//this enables lazlo to be scripted via lua instead of GO, which is 
//preferable in some contexts (no recompiles for changes etc..).
var LuaMod = &lazlo.Module{
	Name:	`LuaMod`, 
	Usage: `%HIDDEN% this module implements lua scripting of lazlo`,
	Run:	 luaMain, 
}

//Each LuaScript represents a single lua script/state machine
type LuaScript struct{
	Robot			*Robot
	State			*lua.LState
}

//A CBMap maps a specific callback Case to it's respective lua function and
//parent lua script
type CBMap struct{
	Func			lua.LFunction
	Callback		*reflect.Value
	Script		*LuaScript
}

type Robot struct{
	ID		string
}
func (r *Robot) GetID() string{
	return r.ID
}

//LuaScripts allows us to Retrieve lua.LState by LuaScript.Robot.ID
var LuaScripts map[string]LuaScript

//CBtable allows us to Retrieve lua.LState by callback case index
var CBTable []CBMap

//Cases is used by reflect.Select to deliver events from lazlo
var Cases []reflect.SelectCase

//Broker is a global pointer back to our lazlo broker
var broker *lazlo.Broker

//luaMain reads in lua files, registers their callbacks with lazlo
//and brokers back events that occur in the chatroom
func luaMain (b *lazlo.Broker){
	 broker=b
    L := lua.NewState()
    defer L.Close()

	 //register functions here
    L.SetGlobal("hear", luar.New(L,Hear)) /* Original lua_setglobal uses stack... */

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
	l := LuaStates[index]
	lmsg := l.NewTable()
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
func (r *Robot)Hear(pat string, lfunc lua.LFunction){
	// it's important that the index is the same in both cbtable and cases
	if len(CBTable) != len(Cases){ panic(`cbtable != cases`) }
	cb:=lazlo.MessageCallback(pat,1)
	cbEntry:=CBMap{
		Func:			lfunc,
		Callback:	reflect.ValueOf(cb),
		Script		LuaScripts[r.ID],
	}
	caseEntry:=reflect.Case{
        Dir:		SelectRecv,
        Chan:		reflect.ValueOf(cb.Chan),
	}
	append(CBTable, cbEntry)
	append(Cases, caseEntry)
}

func (r *Robot)Respond(L *lua.LState) int {
    pat := L.ToString(-1)             /* get argument */

