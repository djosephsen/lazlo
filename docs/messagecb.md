# The message callback

Message callbacks tell Lazlo that you're interested in being notified when
someone says something specific in chat. The function call looks like this: 

```
cb := broker.MessageCallback(`flip table`, true)
```

The formal definition is [here](https://github.com/klaidliadon/lazlo/blob/master/lib/callbacks.go#L113) in callbacks.go

As you can see, *MessageCallback* is a method on the lazlo.Broker type.
Whenever you write a lazlo module, you'll always be passed a pointer to the
broker ([more about that here](plugins.md))

MessageCallback requires two arguments, the first is a string, and the second a
bool. The string argument specifies a regular expression that matches the text
you want to be notified of.  The bool specifies weather or not this message is
a formal command addressed to the bot. When set to *true*, Lazlo will only
broker the message back to your module if it's prefaced with the name of the
your bot. In other words, if you named your bot lazlo, the above example will
only fire when a user says "lazlo flip table". If you set it to *false* it
would fire any time anyone said "flip table"

Lazlo will give you back an object of type lazlo.MessageCallback. It looks like
this: 

```
type MessageCallback struct {
   ID        string
   Pattern   string
   Respond   bool // if true, only respond if the bot is mentioned by name
   Chan      chan PatternMatch
   SlackChan string // if set filter message callbacks to this Slack channel
}
```

*ID* uniquely identifies this callback with the broker. It's set automatically
when you register the module. 

*Pattern* is the regex you specified 

*Respond* is the bool you specified 

*Chan* is the cool part, I'll talk about that next

*SlackChan* is an optional third argument to MessageCallback(). It's a string
that  specifies a slack channel on which you want this callback active (if, for
example, you want *bot deploy* to work in the ops channel but not the main
channel

## Waiting for something to happen
So about this cb.Chan thing; This is the Go channel that Lazlo will use to pass
you back messages that met your criteria. A common pattern is to block on it
waiting for someone to type something that matches your regex like this: 

```
pm := <-cb.Chan
pm.Event.Respond(`(╯°□°）╯︵ ┻━┻`)
``` 

If you do that though, your module will respond to the first event lazlo
brokers back to you, but never again.  So you probably want to block in a for
loop like this: 

```
for{
	pm := <-cb.Chan
	pm.Event.Respond(`(╯°□°）╯︵ ┻━┻`)
}
``` 
That way, your module will block waiting for a message from lazlo and respond
when it gets one, but then it'll loop back around and block again, waiting for
another message.

Remember you can register as many callbacks as you like. When you do that, a
common pattern is to select between them like so: 


```
cb1 := broker.MessageCallback(`flip table`, true)
cb2 := broker.MessageCallback(`I am [0-9]+ mads right now`, false)
cb3 := broker.MessageCallback(`who moved my cheese`, false)

for{
	select{
	case pm := <-cb1:
		go flipFunc(pm)
	case pm := <-cb2:
		go madsFunc(pm)
	case pm := <-cb3:
		go cheeseFunc(pm)
}
``` 

## lazlo.PatternMatch
You may have noticed that the callback chan was type lazlo.PatternMatch. Good
job! That was very observant of you. PatternMatch is a pretty simple type. It
looks like this: 

```
type PatternMatch struct {
   Event *Event
   Match []string
}
```

*Event* is a pointer to the *lazlo.Event* that matched your regex. It is a very
large and intricate type that gives you everything lazlo knows about the event
including who said the thing that matched your regex, what channel they said it
in, what time it was when they said it and so on and so forth. You can find the
Event type definition in lib/slacktypes. Event also comes with a couple
convienience functions that you can use to respond to the message: *Reply()*,
and *Respond()* (which I used in the example above. The distinction is that
*Reply()* echo's the users name. So if you pm.Event.Reply('I see you') to
something Tess said, Lazlo will literally say "Tess, I see you".

*Match* is an array of strings. It is exactly what you would get back if you
called [regex.FindAllStringSubmatch]() on pm.Event.Text with your regex
pattern, because that's literally what lazlo is doing for you internally. To
make a very long story short, this means you can use your regex to capture
strings like so: 

```
cb := b.MessageCallback(`(?i)(tell me how) (\w+) (I am)`, true)
pm := <- cb.Chan
```

Then if I said "Lazlo tell me how pretty I am" in channel, pm.Match[2] would
equal "pretty", so this...

```
pm.Event.Reply(fmt.Printf(`ZOMG you are 42 %s`, pm.Match[2]))
```

... would make lazlo respond: "Dave: ZOMG you are 42 pretty"


