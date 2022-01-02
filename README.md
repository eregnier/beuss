# beuss

Simple message bus mainly for studying socket with golang.

But yay ! it works it is efficient fair enough for simple usages and is dead simple to use.

## What it does

What's more it is designed to have zero dependencies and being very lightweight. So the target binary is about 2.5Mb to run the whole thing and requires few memory (~1Mb)

Messages data live in memory as byte content from sent messsages, all remains in memory as long as it is not consumed. once consumed, it is removed from server memory.

By it's nature (pure golang) it shoud be very portable to the main OS (linux, raspberry pi, windows, macos)

This program send binary message, with by default 10Kb per message.

It is possible to use the cli to produce / consume messages and go client library to implement it in your own program (see tests for exemples how to do this)

## What it does not do

Scalability : it is designed to be simple. As it is seems to handle 10000 small messages I/O in a second with no problem. I did not tested more.

Auto reconnect : This is made to be dead simple. If you want reconnect support or just more resiciency, think about using it with pm2 or supervisord which should do all the hard things of resurect the program cleverly.


## How it works ?

You can build things yourself or get prebuilt binaries from release page.

Just run the server app

```
./beuss
```

then use the cli cmd tool to manage messages. issue the following command to get help :

```
./cli 
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

By default all listen to localhost:6552. it is possible to listen to custom host:port by running commands with an env variable like

```
BPORT=6222 BHOST=my.server.org ./beuss GET <queueName>
```

By default messages size are 10Kb. It is possible to override this size by setting a custom BMAXMESSAGESIZE variable in byte.

```
BMAXMESSAGESIZE=1000000 ./beuss GET <queueName>
```

let send binary content of ~1Mb. Then simply use cli as :

```
cat mypic.jpg | ./beuss PUT <queueName>
```

and then fetch it:

```
./beuss GET <queueName> /tmp/picfrombeuss.jpg
```


## What's next

I want to code clients for this app in python and javascript with the same simple api.
