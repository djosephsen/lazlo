## Up and running in 5 minutes

0: Have [Golang](https://golang.org/doc/install) installed

1: Select *Configure Integrations* from your team menu in slack

2: Add a new *Bots* integration, give your bot a clever name, and take note of your Token

![integration](docs/screenshots/add_bot_integration.png)

3: 
```
	go get github.com/djosephsen/lazlo
```

4: 
```
	export LAZLO_NAME=<whatever you named your bot in the Slack UI>
	export LAZLO_TOKEN=<your token>
	export PORT=5000
	export LAZLO_LOG_LEVEL=DEBUG  # (optional if you'd like to see verbose console messages)
```

5: 
```
lazlo
```

Now you can ping Lazlo to make sure he's alive: 

![](docs/screenshots/ping_lazlo.png)

Lazlo comes with a variety of simple plugins to get you started and give you
examples to work from, and it's pretty easy to add your own. [Making and
managing your own plugins](docs/plugins.md) is pretty much why you're here in
the first place after all.

## Deploy Lazlo to Heroku and be all #legit in 10 minutes

0: Have a github account, a Heroku account, Heroku Toolbelt installed, and upload your ssh key to Github and Heroku

1: Select *Configure Integrations* from your team menu in slack

2: Add a new *Bots* integration, give your bot a clever name, and take note of your Token

3: 
```
go get github.com/kr/godep
```

4: Go to https://github.com/djosephsen/lazlo/fork to fork this repository (or click the fork button up there ^^) 

5 through like 27:  
```
mkdir -p $GOPATH/github.com/<yourgithubname>
cd $GOPATH/github.com/<yourgithubname>
git clone git@github.com:<yourgithubname>/lazlo.git
cd lazlo
git remote add upstream https://github.com/djosephsen/lazlo.git
chmod 755 ./importfix.sh && ./importfix.sh
go get
godep save
heroku create -b https://github.com/kr/heroku-buildpack-go.git
heroku config:set LAZLO_NAME=<whatever you named your bot in the Slack UI>
heroku config:set LAZLO_TOKEN=<your token>
heroku config:set LAZLO_LOG_LEVEL=DEBUG
git add --all .
git commit -am 'lets DO THIS'
git push
git push heroku master
```

At this point you can ping Lazlo to make sure he's alive.

![hi](docs/screenshots/ping_lazlo.png)

### kind of done mostly
When you make changes or add plugins in the future, you can push them to heroku with: 

```
godep save
git add --all .
git commit -am 'snarky commit message'
git push && get push heroku
```

## Use docker to run lazlo and be one of the cool kids in like 42 seconds
(sorry this isn't actually a thing yet)

## What now?
Find out [what lazlo can do](docs/included_plugins.md) out of the box
Get started [adding, removing, and creating plugins](docs/plugins.md)
Learn more about [configuring](docs/configuration.md) Lazlo (there's not much to it)

