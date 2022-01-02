# beuss
simple message bus mainly for studying socket with golang. But yay!  it works is quite efficient and is dead simple to use.

But yay!  it works is quite efficient and is dead simple to use. What's more it is designed to have 0 dependencies and being very lightweight. So the target binary is about 3Mb to run the whole thing and requires few memory (~1Mb)

By it's nature (pure golang) it shoud be very portable to the main OS (linux, raspberry pi, windows, macos)

This program should be able to send binary message, in the limit (hardcoded for now) of 10Kb per message

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
./beuss GET <queueName>
```

This command listen to a queue and produce a string output on each message

```
./beuss ON <queueName>
```

## remote host and custom port

By default all listen to localhost:6552. it is possible to listen to custom host:port by running commands with an env variable like

```
BPORT=6222 BHOST=my.server.org ./beuss GET <queueName>
```

