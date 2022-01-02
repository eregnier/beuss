# beuss

![beuss image](./bus.png)

Simple message bus mainly for studying socket and bus messages system with golang.

But yay ! it works it is efficient fair enough for many usages and is dead simple to use.

## What it does

It is designed to have zero dependencies and being very lightweight. So the target binary is about 2.5Mb to run the whole thing and requires few memory (~1Mb)

It is fifo simple queue message. queues are created on the fly if needed and message are distributed in the order they were added.

Messages data live in memory as byte content from sent messsages, all remains in memory as long as it is not consumed. once consumed, it is removed from server memory.

By it's nature (pure golang) it shoud be very portable to the main OS (linux, raspberry pi, windows, macos)

This program send binary message, with by default 10Kb per message.

It is possible to use the cli to produce / consume messages and go client library to implement it in your own program (see tests for exemples how to do this)

Server code is quite stupid. It add and deliver messages on demand only. This is on purpose. Message consumption event is handled client side by asking often to server if  there is new messages available on a queue.

Messages try to be always processed. It their size is > max message size, then the message is truncated to the max message length. and sent back "as it" to clients.

## What it does not do

Scalability : it is designed to be simple. As it is seems to handle 10000 small messages I/O in a second with no problem. I did not tested more. However this is not this project goal to go into distributed system. instead, I want to keep it simple.

Auto reconnect : This is made to be dead simple. If you want reconnect support or just more resiliency, think about using it with `pm2` or `supervisord` which should do all the hard things of resurect the program cleverly.

Fan system / no index : if many consumer fetch message to the same queue they will get sequential messages back. also message once read are removed from memory. nothing left and there is no plan for persistance at the moment even if it should be possible to add a persistance layer someday.

Thread safe: if there is heavy message consumption from many client, there is no guaranty message will be distributed in the right order.

## How it works ?

You can build things yourself or get prebuilt binaries from release page.

Just run the server app

```
./beuss-server
```

then use the cli cmd tool to manage messages. issue the following command to get help :

```
./beuss
```

run this one to send message to the server

```
echo "my-message-content" | ./beuss PUT <queueName>
```

run this one to retreieve messages one by one from the target queue on the server

```
./beuss GET <queueName> <output.ext>
```

This command listen to a queue and produce a string output on each message

```
./beuss ON <queueName>
```

## Configuration by environment variables

By default all listen to `localhost:6552`. it is possible to listen to custom host:port by running commands with an env variable like

```bash
BPORT=6222 BHOST=my.server.org ./beuss GET <queueName>
```

By default messages size are 10Kb. It is possible to override this size by setting a custom BMAXMESSAGESIZE variable in byte.

```bash
BMAXMESSAGESIZE=1000000 ./beuss-server
```

let send binary content of ~1Mb. Then simply use cli as :

```bash
cat mypic.jpg | ./beuss PUT <queueName>
```

and then fetch it:

```bash
./beuss GET <queueName> /tmp/picfrombeuss.jpg
```


## What's next (ordered timeline)

* I want to code clients for this app in python and javascript with the same simple api in mind.

* I also want to add a fan system that let many clients connect to get the same message for all. this would allow simple distributed messaging.

* Clean queue from memory if it becomes empty for a certain time. For now, if a message is sent on a new queue, it create a queue on the fly and it is never removed event if it is empty for the whole server program life cycle.

* In the future, in case this tool starts to be used somewhere : a scheduler system that would let trigger message creation on specific date or periodically (crontab like recurrence description I guess) that let generate messages alone over time.