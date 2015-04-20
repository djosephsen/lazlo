# Lazlo
## An event-driven chatops framework for Slack in Go. 

The prototypical IRC bot responds to text. Generally, the pattern is you
provide a regex to match on, and some code to run when someone says something
in chat that matches your regular expression. Your plugin runs when a pattern
match happens, and then returns.

Your Lazlo module, by comparison is started at runtime and stays resident in
memory. Outwardly, Lazlo *acts* like a bot, but internally Lazlo works as an
event broker.  Your module registers for callbacks -- you can tell Lazlo what
sorts of events your module finds interesting. For each callback your module
registers, Lazlo will hand back a *channel*. Your module can block on the
channel, waiting for something to happen, or it can register more callbacks (as
many as you have memory for), and select between them in a loop. Throughout its
lifetime, your Module can de-register the callbacks it doesn't need anymore, and
ask for new ones as circumstances demand.

Currently there are three different kinds of callbacks you can ask for.

* [Message callbacks](docs/messagecb.md) specify regex you want to listen for and respond to. 
* [Timer Callbacks](docs/timercb.md) start a (possibly reoccuring) timer (in cron syntax), and notify you when it runs down
* [Link Callbacks](docs/linkcb.md) create a URL that users can click on. When they do, their GET request is brokered back to your module. (Post and Put support coming soon)

Your module can register for all or none of these, as many times as it likes
during the lifetime of the bot. Lazlo makes it easier to write modules that
carry out common chat-ops patterns. For example, you can pretty easily write a
module that: 

1. registers for a message callback for `bot deploy (\w+)` 
2. blocks waiting for that command to be executed
3. when executed, registers for a message callback that matches the specific user that asked for the deploy with the regex: 'authenticate <password>'
4. DM's that user prompting for a password
5. registers a timer callback that expires in 3 minutes
6. Blocks waiting for either the password or the timer
7. Authenticates the user, and runs the CM tool of the week to perform the deploy
8. Captures output from that tool and presents it back to the user
9. de-registers the timer and password callbacks

That's an oversimplified example, but I think you probably get the idea. Check
out the Modules directory for working examples that use the various callbacks. 

[get started](docs/install.md)
