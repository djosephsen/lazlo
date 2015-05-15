package handlers

import (
	"fmt"
	lazlo "github.com/djosephsen/lazlo/lib"
	"regexp"
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
			newUser(b, evnt.[`user`].(lazlo.User))
		case evnt := <-newGroupCB.Chan :
			newGroup(b,cb)
		case evnt := <-newChannelCB.Chan :
			newChannel(b,cb)
		case evnt := <-replaceUserCB.Chan :
			replaceUser(b,cb)
		case evnt := <-replaceGroupCB.Chan :
			replaceGroup(b,cb)
		case evnt := <-replaceChannelCB.Chan :
			replaceChannel(b,cb)
		}
	}
}
	
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

func newUser(b *lazlo.Broker, user lazlo.User){
   for _, exists := range b.Meta.Users {
      if exists.ID == user.ID {
         lazlo.Logger.Debug(`pre-existing user: `, user.Name)
         return
      }
   }
   lazlo.Logger.Debug(`adding user to Meta: `, user.Name)
   b.Meta.Users = append(b.Meta.Users, user)
}

func newChannel(b *lazlo.Broker, chanThingy interface{}) {
   channel := new(sl.Channel)
   jthingy, _ := json.Marshal(chanThingy)
   json.Unmarshal(jthingy, channel)
   for _, exists := range bot.Meta.Channels {
      if exists.ID == channel.ID {
         sl.Logger.Debug(`pre-existing channel: `, channel.Name)
         return
      }
   }
   sl.Logger.Debug(`adding channel to Meta: `, channel.Name)
   bot.Meta.Channels = append(bot.Meta.Channels, *channel)
}
