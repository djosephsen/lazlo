# The event callback

event callbacks tell Lazlo that you're interested in being notified of API
events other than messages. These are things like new users entering a room, or
someone uploading a snippet or image. The Function call looks like this:

```
cb := broker.EventCallback(`type`, `foo`)
```

The formal definition is [here](https://github.com/klaidliadon/lazlo/blob/master/lib/callbacks.go#L132) in callbacks.go

As you can see, *EventCallback* is a method on the lazlo.Broker type.
Whenever you write a lazlo module, you'll always be passed a pointer to the
broker ([more about that here](plugins.md))

EventCallbacks require two string arguments, the first is a key, and the second
a value. Together they give you a way to say "Hey lazlo, tell me whenever you
see an event that has this key with this value". Check out the [slack API
documentation]() for a list of all possible events.

Lazlo will give you back an object of type lazlo.EventCallback. It looks like
this: 

```
type EventCallback struct {
   ID   string
   Key  string
   Val  string
   Chan chan map[string]interface{}
}
```

*ID* uniquely identifies this callback with the broker. It's set automatically
when you register the module, and you can use it to deregister the module later
if you want.

*key* is a regex that matches a json key value in the event type that you're interested in.

*val* is a regex that matches the value of the json key in the event type that
you're interested in. 

*Chan* is the cool part, I'll talk more about that in a minute. First, here are a couple examples:


"Lazlo, tell me about any non-message event in the 'ops' channel"
```
opschan := broker.SlackMeta.GetChannelByName(`ops`)
cb := broker.EventCallback(`Channel`,opschan.ID)
thingy := <- cb.Chan
```

"Lazlo, tell me whenever Tess does anything (except chatting)"
```
tess := broker.SlackMeta.GetUserByName(`Tess`)
cb := broker.EventCallback(`User`,tess.ID)
thingy := <- cb.Chan
```

## Waiting for something to happen
cb.Chan is the Go channel that Lazlo will use to pass you back events that
met your criteria. A common pattern is to block on it waiting for something to
happen that matches your regex like this: 

```
thingy := <-cb.Chan
if thingy[`type`] == `user_typing`{
	b.say(`I see you typing Tess`)
}
``` 

If you do that though, your module will respond to the first event lazlo
brokers back to you, but never again.  So you probably want to block in a for
loop like this: 

```
for{
	thingy := <-cb.Chan
	if thingy[`type`] == `user_typing`{
		b.say(`I see you typing Tess`)
	}
}
```

That way, your module will block waiting for an event from lazlo and respond
when it gets one, but then it'll loop back around and block again, waiting for
the next matching event.

Remember you can register as many callbacks as you like. When you do that, a
common pattern is to select between them like so: 


```
helpChannel := broker.SlackMeta.GetChannelByName(`help`)
cb1 := broker.EventCallback(`type`, `user_entered`)
cb2 := broker.EventCallback(`type`, `new_channel`)
cb3 := broker.EventCallback(`channel`, helpChannel)

for{
	select{
	case thingy := <-cb1.Chan:
		go handle_new_users(thingy)
	case thingy := <-cb2.Chan:
		go handle_new_channels(thingy)
	case thingy := <-cb3.Chan:
		go handle_somebody_needs_help(thingy)
}
``` 

## Thingy
The thingies passed to you by lazlo.EventCallback.Chan are type
map[string]interface{}. That's why I refer to them as thingies; they're just
JSON blobs from the slack API, wrapped up in a map of interfaces. You can
implement the types yourself if you want, or just use them directly. Refer to
the [slack API]() documentation for the response structure of the object you're
asking for. 


## In the Future
I'll probably add an additional EventCallback() method that works in a list
context.. accepting more than one key/value pair, so you can specify multiple
criteria to match. 
