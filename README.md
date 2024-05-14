# Broadsocket
broadcast websocket messages to every client

## About
Broadsocket is a websocket testing server, where every message is broadcasted to **all connected clients**.
This is useful if you want to write a service that depends on a websocket server that does not yet exists.

Message broadcasts can be scoped in topics, i.e. only clients connected to the same topic will receive messages from another.

## Features
* Simple web-based UI
* Scoping messages into topics
* Single-file, zero-dependency executable
* Ephemeral messaging, no database

## Usage
Broadsocket is written to work as a single executable with no runtime dependencies. Simply run it and optionally specifiy the bind address:

```console
$ ./broadsocket [-b <host>:<port>]
$
```

By default broadsocket is served at `localhost:8888`.

The websockets are now available at `localhost:8888/.ws/<topic>`. To read and send messages via a basic UI simply discard the `.ws` part of the URL (`localhost:8888/<topic>`). The topic can basically be any path appended to the root path (everaything after `#` and `? ` is ignored because of the way URLs work). New topics are automatically created when they are first called and automatically deleted if there are no clients connected anymore.

## Building
Right now broadsocket does not provide any prebuilt binaries. But don't worry, compiling your own is pretty easy. THe only requirement is a go compiler.

First, clone the repository and enter the it:
```console
$ git clone https://github.com/bliepp/broadsocket.git
$ cd broadsocket
$
```

### For production
Simply run
```console
$ go build .
$
```
You should find an executable called `broadsocket` in the project's directory (`broadsocket.exe`on windows). You can directly execute it using `./broadsocket` or add it to your `$PATH` variable.

If you need a true zero-dependency build, make sure to disable `CGO`:
```console
$ CGO_ENABLED=0 go build .
$
```

### For development
When you want to modify the source code of broadsocket I'd recommend [cosmtrek's air](https://github.com/cosmtrek/air). With that, simply run:
```console
$ air

  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.51.0, built with Go go1.22.0

watching .
watching broadcast
!exclude tmp
building...
running...
2024/05/14 09:05:53 Starting server at localhost:8888
```