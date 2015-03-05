# Slacker 
## A neck-beardy chatops bot for the discriminating Ops professional

Written in Go for SlackHQ (using the RTM API), Slacker provides features that
make it easy to model and abstract many different types of interaction between
engineers and production environments. 

 * Define Event and Message [Handlers](docs/handlers.md) so Slacker will parse and respond to chat conversation and commands
 * Specify Cron-Syntax Schedules for Slacker to run periodic [*Chores*](docs/chores.md)
 * In-Memory or Redis-backed [persistent storage](docs/brain.md) built-in (aka hubot.brain)
 * Set up inbound and outbound [filters](docs/filters.md) to modify events enroute to or from SlackHQ
 * Correctly handle and respond to system-level events like SigTerms by configuring runtime [hooks](docs/hooks.md)
 * [Configures](docs/config.md) via env-vars and logs to stdout like a well behaved 12-factor app
 * Efficient, parallel execution of handlers, filters, and chores
 * Simple semantics for creating and linking-in your own plugins
 * (Pretty) complete [support for the Slack Web API](docs/api-calls.md)

## Check out Slacker from your workstation in 5 minutes

1: Select *Configure Integrations* from your team menu in slack

2: Add a new *Bots* integration, give your bot a clever name, and take note of your Token

![integration](docs/screenshots/add_bot_integration.png)

3: 
```
	go get github.com/djosephsen/slacker
```

4: 
```
	export SLACKER_NAME=<whatever you named your bot in the Slack UI>
	export SLACKER_TOKEN=<your token>
	export SLACKER_LOG_LEVEL=DEBUG  # (optional if you'd like to see verbose console messages)
```

5: 
```
slacker
```

At this point you should see slacker join your default channel and say hi. 

![hi](docs/screenshots/hi.png)

If you ctl-C him in the console window, you should also see him say goodbye
before he leaves the channel. 

![bye](docs/screenshots/bye.png)

These behaviors are both implemented as inithooks (the first as a startup hook
and the second as a shutdown hook). If you find this annoying you can comment
them out from the [yourPluginsGoHere.go](yourPluginsGoHere.go) file, which is the
file that controls all of the plugins Slacker uses to interact with you and
your production environment.  Slacker comes with a variety of simple plugins to
get you started and give you examples to work from, and it's pretty easy to add
your own. [Making and managing your own plugins](docs/plugins.md) is pretty
much why you're here in the first place after all.

## Deploy Slacker to Heroku and be all #legit in 10 minutes

0: Have a github account, a Heroku account, Heroku Toolbelt installed, and upload your ssh key to Github and Heroku

1: Select *Configure Integrations* from your team menu in slack

2: Add a new *Bots* integration, give your bot a clever name, and take note of your Token

3: 
```
go get github.com/kr/godep
```

4: Go to https://github.com/djosephsen/slacker/fork to fork this repository (or click the fork button up there ^^) 

5 through like 27:  
```
mkdir -p $GOPATH/github.com/<yourgithubname>
cd $GOPATH/github.com/<yourgithubname>
git clone git@github.com:<yourgithubname>/slacker.git
cd slacker
git remote add upstream https://github.com/djosephsen/slacker.git
chmod 755 ./importfix.sh && ./importfix.sh
go get
godep save
heroku create -b https://github.com/kr/heroku-buildpack-go.git
heroku config:set SLACKER_NAME=<whatever you named your bot in the Slack UI>
heroku config:set SLACKER_TOKEN=<your token>
heroku config:set SLACKER_LOG_LEVEL=DEBUG
git add --all .
git commit -am 'lets DO THIS'
git push
git push heroku master
```

At this point you should see slacker join your channel.

![hi](docs/screenshots/hi.png)

When you make changes or add plugins in the future, you can push them to heroku with: 

```
godep save
git add --all .
git commit -am 'snarky commit message'
git push && get push heroku
```

## What now?
Find out [what slacker can do](docs/builtins.md) out of the box
Get started [adding, removing, and creating plugins](docs/plugins.md)
Learn more about [configuring](docs/configuration.md) Slacker (there's not much to it)

## Why Slacker? 

Slacker wants to be a quality, featureful Slack-specific chatbot with a focus
on operations-abstractions. It borrows heavily in design and implementation
from [Hal](https://github.com/danryan/hal) and
[gopherbot](https://github.com/daph/gopherbot), but focuses less on giving you
a library to write a chatbot, and more on giving you a chatbot.

I made Slacker to be a tool. Through it, I'm hoping to provide a flexible
framework that can solve the operational needs of engineers, like deploying
code to production, providing timely metric data, running automated break-fix
and defensive network re-configurations and stuff like that.  

I also made Slacker to be an anthropomorphic bot whose personality you can
shape, because people don't use what they don't like. But things that are fun
are interesting, and things that are interesting get adopted, improved, used,
and cherished.

If you're in a chatops shop today and that sounds interesting to you I'd love
your comments and help. Actually, even if you're not in a chatops shop today
and this looks interesting to you I'd love your comments and help. 

## Current Status

Slacker is basically working and basically documented. All plugin types are
implemented and functional and there are several included plugins of varying
degress of complexity that should help get you started writing your own
plugins. 

### Todo's in order of when I'll probably get to them: 

* I'm considering some drastic changes to the Broker code to support the bullet point above
* Integrated statsd support for emitting metrics
* Transparent support for custom [slash-commands](https://dbgone.slack.com/services/new/slash-commands)
* Other loftier stuff like redundancy, failover, and Lua-Plugins that I'm too
busy to seriously think about yet.
