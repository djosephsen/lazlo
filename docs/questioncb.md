# The question callback

Imagine for a moment that you want to ask a user a question in chat. The
problem is, there's no way for the bot to discern between a user-response, and
users just talking. We can DM the user, and assume that the next thing they say
is the reply. In practice this works pretty well, but think about what you need
to do programatically to make that happen.. you have to: 

1. Grab a DM channel ID with the user
2. Send the user a question
3. Setup an event-callback to listen on the DM channel for the user's response
4. Parse the user's response out of the callback 

The question callback does all this for you. You provide a UserID and the question (both in the form of strings), like this: 

```
cb := broker.QuestionCallback(user.ID, `why is the sky blue?`)
```


You'll get back an object of type lazlo.QuestionCallback, which looks like this:

```
type QuestionCallback struct {
   ID       string
   User     string
   DMChan   string
   Question string
   Answer   chan string
}
```

*ID* uniquely identifies this callback with the broker. Currently, because of
how questions work internally, you can't deregister them using their ID (well
you can, but the question will live on, in the serialization queue (more on
that below). I'm working on this sorry.

*User* is the User ID you provided when you asked for the callback

*DMChan* is the ID of the DM channel lazlo setup with the user. 

*Question* is the string you provided as the question. 

*Answer* is a channel that you can listen on for the user's response.

## Serialization
Because the question callback is available to any plugin that wants to use it,
and more than one plugin might decide to ask the same user a question at
roughly the same time, the Question-Callback subsystem works as a serializer
service; it automatically serializes questions in a first-in, first-out manner,
making sure that each question only gets asked once, and a user only has one
question to answer at a time.

A caveat of this behavior is that it breaks the broker.Deregister() function
(because questions are tracked in a separate serialization queue internally.
This might be a problem if you want to time-out a question, if, for example a
user goes on vacation or otherwise just isn't around to answer. I'm still
thinking about the cleanest way to implement question cancelations. Feel free
to throw me an issue if you really need it meow.

## Ask and you shall receive
Question callbacks are actually pretty neat to work with. You can, for example implement decision tree type stuff: 

```
q1 := broker.QuestionCallback(user.ID,`WHAT is your name?`
ans := <- q1.Answer
if ans == user.Name {
	q2:= broker.QuestionCallback(user.ID,`WHAT is your quest?`)
	ans = <- q2.Answer
	if ans == `I seek the grail`{
		q3:= broker.QuestionCallback(user.ID, randQuestion)
		ans = <- q3.Answer
		if checkAnswer(ans){
			broker.Say(`You may pass`)
			return 0
		}
	}
}
killUserWithLightening(user.ID)

``` 

Or you can rely on the nature of the serialization service to front-load all of
your questions: 

```
q1:=broker.QuestionCallback(user.ID,`WHAT is your name?`)
q2:=broker.QuestionCallback(user.ID,`WHAT is your quest?`)
q3:=broker.QuestionCallback(user.ID,`WHAT is the airspeed of a swallow?`)

name:=<-q1.Answer
quest:=<-q2.Answer
airSpeedOfASwallow:=<-q3.Answer
```

## Beware the limit
The serialization queue is itself a buffered channel with a hard-coded limit of
100. In other words, there can only be 100 open questions at a time in the
queue. In the future I might add an environment variable for this, but it's pretty easy to change in broker.go (hint: grep for "LIMIT ALERT"). 

## In the Future
Currently there's no way to de-register a question callback, which is a problem
if, for example you ask a question from someone who went on vacation for
example, or you want to timeout a question for any other reason. So yeah, I
need to fix that.
