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
//preferable in some contexts (simpler(?), no recompiles for changes etc..).
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
	ID		int
}
func (r *Robot) GetID() int{
	return r.ID
}

//LuaScripts allows us to Retrieve lua.LState by LuaScript.Robot.ID
var LuaScripts []LuaScript

//CBtable allows us to Retrieve lua.LState by callback case index
var CBTable []CBMap

//Cases is used by reflect.Select to deliver events from lazlo
var Cases []reflect.SelectCase

//Broker is a global pointer back to our lazlo broker
var broker *lazlo.Broker

//luaMain creates a new lua state for each file in ./lua
//and hands them the globals they need to interact with lazlo
func luaMain (b *lazlo.Broker){
	 broker=b
	 for <some array of lua files to read>{

		//make a new script entry
		script := &luaScript{
			Robot:	&Robot{
				ID:	len(LuaScripts),
			},
			State:	lua.NewState(),
		}
    	defer script.State.Close()

      // register hear and respond
    	script..SetGlobal("robot", luar.New(script.State,script.Robot))
		LuaScripts[script.Robot.ID]=script

		// the lua script will register callbacks to the Cases
		if err := script.State.DoFile(file); err != nil {
			panic(err)
		}

	//block waiting on events from the broker
	for{
		index, value, _ := reflect.Select(Cases)
		handle(index, value.Interface())	
	}
    L.Push(lua.LNumber(lv * 2)) /* push result */
}

//handle takes the index and value of an event from lazlo, 
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
	l := CBTable[index].Script.State
	lmsg := luar.New(message)

	if err := l.CallByParam(lua.P{
    Fn: CBTable[index].Func,
    NRet: 0,
    Protect: true,
    }, lmsg ); err != nil {
    panic(err)
	}
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

//creates a new message callback from robot.hear/respond
func newMsgCallback(pat string, lfunc lua.LFunction, isResponse bool){
	// cbtable and cases indexes have to match 
	if len(CBTable) != len(Cases){ panic(`cbtable != cases`) }
	cb:=lazlo.MessageCallback(pat,isResponse)
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

//functions exported to the lua runtime below here

//lua function to overhear a message
func (r *Robot)Hear(pat string, lfunc lua.LFunction){
	newMesgCallback(pat, lfunc, false)
}

//lua function to process a command 
func (r *Robot)Respond(pat string, lfunc lua.LFunction){
	newMesgCallback(pat,lfunc,true)
}

//lua function to reply to a message passed to a lua-side callback
func (pm *lazlo.PatternMatch)Reply(words string){
	pm.Event.Reply(words)
}
