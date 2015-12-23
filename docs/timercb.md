# The Timer Callback

Timer Callbacks let you arrange for Lazlo to wake your module up every so
often.  There are many different reasons you might want to implement a timer. I
often use them to periodically query an external service like an RSS feed. They also make it easy to create bot's that appear to "randomly" say things in channel.

```
cb := broker.TimerCallback(`0 */1 * * * * *`)
```

The formal definition is [here](https://github.com/klaidliadon/lazlo/blob/master/lib/callbacks.go#L147) in callbacks.go

As you can see, *TimerCallback()* is a method on the lazlo.Broker type.
Whenever you write a lazlo module, you'll always be passed a pointer to the
broker ([more about that here](plugins.md))

*TimerCallback()* requires one argument: the schedule, which is a string that
specifies a cron-syntax schedule.

Lazlo will give you back an object of type lazlo.TimerCallback. It looks like
this: 

```
type TimerCallback struct {
   ID       string
   Schedule string
   State    string
   Next     time.Time
   Chan     chan time.Time
}
```

*ID* uniquely identifies this callback with the broker. It's set automatically
when you register the module and you can use it to deregister the module later
if you want. 

*Schedule* is the cron schedule you specified 

*State* is a human-readable string that describes the current state of the
timer (it will have an error if something went wrong with the schedule you specified)

*Next* is a [time.Time]() value that specifies the next time this timer will
fire.  

*Chan* is a Channel of type [time.Time](). When the "alarm" goes off so to
speak, Lazlo notifies you using this channel. 

## Waiting for something to happen
A common pattern is to block on the callback's Chan attribute, waiting for the
timer to fire.

```
alarm := broker.TimerCallback(`0 0 14 * * * *`) //wake me up at 2pm every day
for{
	alarm := <-timer.Chan
	broker.Say(`welp it's 2pm, time to update the DB`)
	update_db_from_RSS_feed()
``` 

You can combine timers with other callbacks to achieve more advanced patterns,
like verifying dangerous chatops commands


```
shieldsUp := broker.MessageCallback(`(?i)shields up$`, true)

for{
	// block waiting for a shields up command
	scb := <-shieldsUp.Chan 

	// register new verify and timeout callbacks
	scb.Reply(`Verify shields up by pasting this verification code: 234567`)
	verify := broker.MessageCallback(`234567`)
	timeout := broker.TimerCallback(`30 * * * * * *`)

	// Give the user 30 seconds to verify it wants the shields up
	select{
	case v := <- verify.Chan:
		go putTheShieldsUP()
	case v := <-timeout.Chan:
		go nevermind()
	}
	broker.DeRegister(verify)
	broker.DeRegister(timeout)
}
``` 
