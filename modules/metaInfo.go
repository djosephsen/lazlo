package handlers

import (
	lazlo "github.com/djosephsen/lazlo/lib"
	"fmt"
	"regexp"
	"encoding/json"
)

var MetaInfo = lazlo.Module{
	Name:    `MetaInfo`,
	Usage:   `%BOTNAME% (channel|user) (list|dump) <id>:: examine meta-info on channels and users`,
	Run:		metaRun,
}

func metaRun(b *lazlo.Broker){
	listCB := b.MessageCallback(`(?i)(channel|user)s* (list|dump) *(\w+)*`, true)
	newUserCB := b.EventCallback(`type`,`team_join`)
	newGroupCB := b.EventCallback(`type`,
	newChannelCB := b.EventCallback(`type`,
	replaceUserCB := b.EventCallback(`type`,
	replaceGroupCB := b.EventCallback(`type`,
	replaceChannelCB := b.EventCallback(`type`,

	for{
		select{
		case pm := <-listCB.Chan :
			listOrDump(b, cb)
		case evnt := <-newUserCB.Chan :
			modUser(b, evnt.[`user`])
		case evnt := <-newGroupCB.Chan :
			modGroup(b,cb)
		case evnt := <-newChannelCB.Chan :
			modChannel(b,cb)
		case evnt := <-replaceUserCB.Chan :
			modUser(b,cb)
		case evnt := <-replaceGroupCB.Chan :
			modGroup(b,cb)
		case evnt := <-replaceChannelCB.Chan :
			modChannel(b,cb)
		}
	}
}

//Figure out weather to list or dump a thingy and list or dump it
func listOrDump(b *Broker, pm *lazlo.PatternMatch){
	match:=pm.Match
	typeOfThing := match[1]
	cmd := match[2]
	id := match[3]
	var reply string
	if matches, _ := regexp.MatchString(`(?i)list`, cmd); matches {
		reply = listThing(b, typeOfThing)
	} else if matches, _ := regexp.MatchString(`(?i)dump`, cmd); matches {
		reply = dumpThing(b, typeOfThing, id)
	}
	if reply != `` {
		b.Reply(reply)
	}
}

//List a thingy
func listThing(b *lazlo.Broker, typeOfThing string) string {
	var reply string
	if matches, _ := regexp.MatchString(`(?i)channel`, typeOfThing); matches {
		reply = `Channels:`
		for _, c := range b.Meta.Channels {
			reply = fmt.Sprintf("%s\n%s (%s)", reply, c.ID, c.Name)
		}
	} else if matches, _ := regexp.MatchString(`(?i)user`, typeOfThing); matches {
		reply = `Users:`
		for _, u := range b.Meta.Users {
			reply = fmt.Sprintf("%s\n%s (%s)", reply, u.ID, u.Name)
		}
	}
	return reply
}

//Dump a thingy
func dumpThing(b *lazlo.Broker, typeOfThing string, id string) string {
	var reply string
	if matches, _ := regexp.MatchString(`(?i)channel`, typeOfThing); matches {
		channel := b.Meta.GetChannel(id)
		reply := fmt.Sprintf("Channel: %s", channel.Name)
		reply = fmt.Sprintf("%s\n%s", reply, channel)
	} else if matches, _ := regexp.MatchString(`(?i)user`, typeOfThing); matches {
		user := b.Meta.GetUser(id)
		reply := fmt.Sprintf("User: %s", user.Name)
		reply = fmt.Sprintf("%s\n%s", reply, user)
	}
	return reply
}

//Add a new user or replace it if it already exists
func modUser(b *lazlo.Broker, userThingy interface{}){
	user:= new(lazlo.User)
	juser,_:=json.Marshal(user)
	json.Unmarshal(juser,user)
   for k, exists := range b.Meta.Users {
      if exists.ID == user.ID {
         lazlo.Logger.Debug(`pre-existing user: `, user.Name)
			l:=len(b.Meta.Users)
			b.Meta.Users[k]=b.Meta.Users[l-1]
			b.Meta.Users=b.Meta.Users[:l-1]
         return
      }
   }
   lazlo.Logger.Debug(`adding user to Meta: `, user.Name)
   b.Meta.Users = append(b.Meta.Users, user)
}

//Add a new Channel or replace it if it already exists
func modChannel(b *lazlo.Broker, chanThingy interface{}) {
   channel := new(b.Channel)
   jthingy, _ := json.Marshal(chanThingy)
   json.Unmarshal(jthingy, channel)
   for _, exists := range b.Meta.Channels {
      if exists.ID == channel.ID {
			l:=len(b.Meta.Channels)
			b.Meta.Channels[k]=b.Meta.Channels[l-1]
			b.Meta.Channels=b.Meta.Channels[:l-1]
         lazlo.Logger.Debug(`pre-existing channel: `, channel.Name)
         return
      }
   }
   lazlo.Logger.Debug(`adding channel to Meta: `, channel.Name)
   b.Meta.Channels = append(b.Meta.Channels, *channel)
}

//Add a new group or replace it if it already exists
func modGroup(b *lazlo.Broker, groupThingy interface{}){
	group:= new(lazlo.Group)
	jgroup,_:=json.Marshal(group)
	json.Unmarshal(jgroup,group)
   for k, exists := range b.Meta.Groups {
      if exists.ID == group.ID {
         lazlo.Logger.Debug(`pre-existing group: `, group.Name)
			l:=len(b.Meta.Groups)
			b.Meta.Groups[k]=b.Meta.Groups[l-1]
			b.Meta.Groups=b.Meta.Groups[:l-1]
         return
      }
   }
   lazlo.Logger.Debug(`adding group to Meta: `, group.Name)
   b.Meta.Groups = append(b.Meta.Groups, group)
}
