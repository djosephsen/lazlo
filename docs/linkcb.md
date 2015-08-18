# The Link Callback

Link callbacks let you present your users with embedded URLs they can click on.
I use them to present my users with multiple-choice questions, 

```
cb := broker.LinkCallback(`option1`)
cb := broker.LinkCallback(`option2`)
```

The formal definition is [here]() in httpserver.go

Link callbacks are a lot different from the other types, and you have to know a
little bit more about how lazlo works to use them effectively. The first thing
you need to know, is that lazlo has a built-in HTTP server. This is implemented
in lib/httpserver.go.

The next thing you need to know is the built-in HTTP server uses two
environment variables to configure itself. The first of these is *PORT*, and
the second is *LAZLO_URL*

The *PORT* environment variable dictates the TCP port lazlo listens on. Paas
services like Heroku set this variable for you automatically when they run your
app in a dyno, but if you run Lazlo on localhost, you need to set this or lazlo
won't run.

The *LAZLO_URL* variable is an optional variable you need to set if you're
using link callbacks. It specifies the public address of the server you're
running Lazlo on. The linkCallback uses this value to create a clickable URL
that lazlo will be able to see. 

## Now you know what you need to know
Ok, lets say you're running lazlo on an AWS instance, and you set PORT to 5000,
and LAZLO_URL to 54.87.22.100, and then you create a link callback like so: 

```
cb := broker.LinkCallback(`option1`)
```

Lazlo is going to create a local API endpoint at /option1, and give you back a
struct of type lazlo.LinkCallback. It looks like this: 

```
type LinkCallback struct {
   ID      string
   Path    string // the computed URL
   URL     string
   Handler func(res http.ResponseWriter, req *http.Request)
   Chan    chan *http.Request
}
```

*ID* uniquely identifies your callback (as you're probably used to by now)

*PATH* is the URL Lazlo intends for you to present your user. In our current
example, PATH would equal this: ``` http://54.87.22.100:5000/option1 ```

*URL* is the URL you passed in to broker.LinkCallback()

*Handler* if you're familiar with [pat](http://github.com/bmizerany/pat), you
can pass in your own http handler function as an optional second argument to
broker.LinkCallback(). If that sounded like greek to you, no worries; ignore
it.

*Chan* here's the magic part. Whenver a user visits the link specified
by your callback's *Path* attribute, Lazlo will hand your module an
[http.Request]() that represents their click.

## Wait, what?

;tldr, if you properly set up the PORT and LAZLO_URL environment variables, and Lazlo is reachable from your user's browser, then you can do this:

```
option1 := broker.LinkCallback(`option1`)
option2 := broker.LinkCallback(`option2`)
```

And lazlo will wire up http://localhost/option1 and http://localhost/option2 to
your callback's Chan channel. 

You can present those links to your users by wrapping them in [slack compatible
markdown]() like this: 

```
options := fmt.Sprintf("choose <%s|option1> or <%s|option2>", option1.URL, option2.URL)
broker.Say(options)
```

Then you can block waiting for the user to click on an answer: 

```
for {
	select {
	case choice := <-option1.Chan:
		broker.Say('you chose option1')
	case choice := <-option2.Chan:
		broker.Say('you chose option2')
	}
}
```
