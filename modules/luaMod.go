package modules

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	lazlo "github.com/klaidliadon/lazlo/lib"
	luar "github.com/layeh/gopher-luar"
	lua "github.com/yuin/gopher-lua"
)

//luaMod implements a lua-parsing plugin for Lazlo.
//this enables lazlo to be scripted via lua instead of GO, which is
//preferable in some contexts (simpler(?), no recompiles for changes etc..).
var LuaMod = &lazlo.Module{
	Name:  `LuaMod`,
	Usage: `%HIDDEN% this module implements lua scripting of lazlo`,
	Run:   luaMain,
}

//Each LuaScript represents a single lua script/state machine
type LuaScript struct {
	Robot *Robot
	State *lua.LState
}

//Keep a local version of lazlo.Patternmatch so we can add methods to it
type LocalPatternMatch lazlo.PatternMatch

//A CBMap maps a specific callback Case to it's respective lua function and
//parent lua script
type CBMap struct {
	Func     lua.LValue
	Callback reflect.Value
	Script   *LuaScript
}

type Robot struct {
	ID int
}

func (r *Robot) GetID() int {
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
func luaMain(b *lazlo.Broker) {
	broker = b
	var luaDir *os.File
	luaDirName := "lua"
	if luaDirInfo, err := os.Stat(luaDirName); err == nil && luaDirInfo.IsDir() {
		luaDir, _ = os.Open(luaDirName)
	} else {
		lazlo.Logger.Error("Couldn't open the Lua Plugin dir: ", err)
	}
	luaFiles, _ := luaDir.Readdir(0)
	for _, f := range luaFiles {
		if f.IsDir() {
			continue
		}

		file := fmt.Sprintf("%s/%s", luaDirName, f.Name())

		//make a new script entry
		script := LuaScript{
			Robot: &Robot{
				ID: len(LuaScripts),
			},
			State: lua.NewState(),
		}
		defer script.State.Close()

		// register hear and respond inside this lua state
		script.State.SetGlobal("robot", luar.New(script.State, script.Robot))
		//script.State.SetGlobal("respond", luar.New(script.State, Respond))
		//script.State.SetGlobal("hear", luar.New(script.State, Hear))
		LuaScripts = append(LuaScripts, script)

		// the lua script will register callbacks to the Cases
		if err := script.State.DoFile(file); err != nil {
			panic(err)
		}
	}
	//block waiting on events from the broker
	for {
		index, value, _ := reflect.Select(Cases)
		handle(index, value.Interface())
	}
}

//pmTranslate makes a localized version of patternmatch so we can export it to lua
//with a few additional methods.
func pmTranslate(in lazlo.PatternMatch) LocalPatternMatch {
	return LocalPatternMatch{
		Event: in.Event,
		Match: in.Match,
	}
}

//handle takes the index and value of an event from lazlo,
//typifies the value and calls the right function to push the data back
//to whatever lua script asked for it.
func handle(index int, val interface{}) {
	switch val.(type) {
	case lazlo.PatternMatch:
		handleMessageCB(index, pmTranslate(val.(lazlo.PatternMatch)))
	case time.Time:
		handleTimerCB(index, val.(time.Time))
	case map[string]interface{}:
		handleEventCB(index, val.(map[string]interface{}))
	case *http.Request:
		handleLinkCB(index, val.(*http.Response))
	default:
		err := fmt.Errorf("luaMod handle:: unknown type: %T", val)
		lazlo.Logger.Error(err)
	}
}

//handleMessageCB brokers messages back to the lua script that asked for them
func handleMessageCB(index int, message LocalPatternMatch) {
	l := CBTable[index].Script.State
	lmsg := luar.New(l, message)

	if err := l.CallByParam(lua.P{
		Fn:      CBTable[index].Func,
		NRet:    0,
		Protect: false,
	}, lmsg); err != nil {
		panic(err)
	}
}

//handleTimerCB brokers timer alarms back to the lua script that asked for them
func handleTimerCB(index int, t time.Time) {
	return
}

//handleEventCB brokers slack rtm events back to the lua script that asked for them
func handleEventCB(index int, event map[string]interface{}) {
	return
}

//handleLinkCB brokers http GET requests back to the lua script that asked for them
func handleLinkCB(index int, resp *http.Response) {
	return
}

//creates a new message callback from robot.hear/respond
func newMsgCallback(RID int, pat string, lfunc lua.LValue, isResponse bool) {
	// cbtable and cases indexes have to match
	if len(CBTable) != len(Cases) {
		panic(`cbtable != cases`)
	}
	cb := broker.MessageCallback(pat, isResponse)
	cbEntry := CBMap{
		Func:     lfunc,
		Callback: reflect.ValueOf(cb),
		Script:   &LuaScripts[RID],
	}
	caseEntry := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(cb.Chan),
	}
	CBTable = append(CBTable, cbEntry)
	Cases = append(Cases, caseEntry)
}

//functions exported to the lua runtime below here

//lua function to overhear a message
func (r Robot) Hear(pat string, lfunc lua.LValue) {
	newMsgCallback(r.ID, pat, lfunc, false)
}

/*func Hear(id int, pat string, lfunc lua.LValue){
	newMsgCallback(id, pat, lfunc, false)
}*/

//lua function to process a command
func (r Robot) Respond(pat string, lfunc lua.LValue) {
	newMsgCallback(r.ID, pat, lfunc, true)
}

/*func Respond(id int, pat string, lfunc lua.LValue){
	newMsgCallback(id, pat, lfunc, true)
}*/

//lua function to reply to a message passed to a lua-side callback
func (pm LocalPatternMatch) Reply(words string) {
	pm.Event.Reply(words)
}
