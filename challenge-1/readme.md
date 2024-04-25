# Challenge #1: Echo

Our first challenge is more of a “getting started” guide" to get the hang of working with [Maelstrom](https://github.com/jepsen-io/maelstrom) in Go. In Maelstrom, we create a _node_ which is a binary that receives JSON messages from `STDIN` and sends JSON messages to `STDOUT`. You can find a full [protocol specification](https://github.com/jepsen-io/maelstrom/blob/main/doc/protocol.md) on the Maelstrom project.

We’ve created a [Maelstrom Go library](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go) which provides `maelstrom.Node` that handles all this boilerplate for you. It lets you register handler functions for each message type—similar to how [`http.Handler`](https://pkg.go.dev/net/http#Handler) works in the standard library.

## [](https://fly.io/dist-sys/1//#specification)Specification

In this challenge, your node will receive an `"echo"` message from Maelstrom that looks like this:


```
{
  "src": "c1",
  "dest": "n1",
  "body": {
    "type": "echo",
    "msg_id": 1,
    "echo": "Please echo 35"
  }
}
```

Nodes & clients are sequentially numbered (e.g. `n1`, `n2`, etc). Nodes are prefixed with `"n"` and external clients are prefixed with `"c"`. Message IDs are unique per source node but that is handled automatically by the Go library.

Your job is to send a message with the same body back to the client but with a message type of `"echo_ok"`. It should also associate itself with the original message by setting the `"in_reply_to"` field to the original message ID. This reply field is handled automatically if you use the [`Node.Reply()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#Node.Reply) method.

It should look something like:


```
{
  "src": "n1",
  "dest": "c1",
  "body": {
    "type": "echo_ok",
    "msg_id": 1,
    "in_reply_to": 1,
    "echo": "Please echo 35"
  }
}
```

## [](https://fly.io/dist-sys/1//#implementing-a-node)Implementing a node

For this first challenge, we’ll walk you through how to implement the echo program. First, create a directory for your binary called `maelstrom-echo` and initialize a Go module for it:


```
$ mkdir maelstrom-echo
$ cd maelstrom-echo
$ go mod init maelstrom-echo
$ go mod tidy
```

Then create your `main.go` file in this directory. This should start with the `main` package declaration and a few imports:


```
package main

import (
    "encoding/json"
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)
```

Inside our `main()` function, we’ll start by instantiating our node type:


```
n := maelstrom.NewNode()
```

From here, we can register a handler callback function for our “echo” message. This function accepts a [`maelstrom.Message`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#Message) which contains the source and destination nodes for the message as well as the body content.


```
n.Handle("echo", func(msg maelstrom.Message) error {
    // Unmarshal the message body as an loosely-typed map.
    var body map[string]any
    if err := json.Unmarshal(msg.Body, &body); err != nil {
        return err
    }

    // Update the message type to return back.
    body["type"] = "echo_ok"

    // Echo the original message back with the updated message type.
    return n.Reply(msg, body)
})
```

In this handler, we’re unmarshaling to a generic map since we simply want to echo back the same message we received. The [`Reply()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#Node.Reply) method will automatically set the source and destination fields in the return message and it will associate the message as a reply to the original one received.

Finally, we’ll delegate execution to the `Node` by calling its [`Run()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#Node.Run) method. This method continuously reads messages from `STDIN` and fires off a goroutine for each one to the associated handler. If no handler exists for a message type, `Run()` will return an error.


```

```

You can find a full implementation of this [`maelstrom-echo` program](https://github.com/jepsen-io/maelstrom/blob/main/demo/go/cmd/maelstrom-echo/main.go) in the Maelstrom codebase.

To compile our program, fetch the Maelstrom library and install:


```
if err := n.Run(); err != nil {
    log.Fatal(err)
}
```

This will build the `maelstrom-echo` binary and place it in your `$GOBIN` path which is typically `~/go/bin`.

## [](https://fly.io/dist-sys/1//#installing-maelstrom)Installing Maelstrom

Maelstrom is built in [Clojure](https://clojure.org/) so you’ll need to install [OpenJDK](https://openjdk.org/). It also provides some plotting and graphing utilities which rely on [Graphviz](https://graphviz.org/) & [gnuplot](http://www.gnuplot.info/). If you’re using Homebrew, you can install these with this command:


```
go get github.com/jepsen-io/maelstrom/demo/go
go install .
```

You can find more details on the [Prerequisites](https://github.com/jepsen-io/maelstrom/blob/main/doc/01-getting-ready/index.md#prerequisites) section on the Maelstrom docs.

Next, you’ll need to download Maelstrom itself. These challenges have been tested against the [Maelstrom 0.2.3](https://github.com/jepsen-io/maelstrom/releases/tag/v0.2.3). Download the tarball & unpack it. You can run the `maelstrom` binary from inside this directory.

## [](https://fly.io/dist-sys/1//#running-our-node-in-maelstrom)Running our node in Maelstrom

We can now start up Maelstrom and pass it the full path to our binary:


```
brew install openjdk graphviz gnuplot
```

This command instructs `maelstrom` to run the `"echo"` workload against our binary. It runs a single node and it will send `"echo"` commands for 10 seconds.

Maelstrom will only inject network failures and it will not intentionally crash your node process so you don’t need to worry about persistence. You can use in-memory data structures for these challenges.

If everything ran correctly, you should see a bunch of log messages and stats and then finally a pleasent message from Maelstrom:

Unwrap text Copy to clipboard

```
Everything looks good! ヽ(‘ー`)ノ
```

Success! If everything is working, move on to the [Unique ID Generation challenge](https://fly.io/dist-sys/2) to build and test a distributed unique ID generator on your own.

If you’re not seeing this success message, head over to our [Fly.io Community forum](https://community.fly.io/) for some help.

## [](https://fly.io/dist-sys/1//#debugging-maelstrom)Debugging maelstrom

If your test fail, you can run the Maelstrom web server to view your results in more depth:

Unwrap text Copy to clipboard

```
./maelstrom serve
```

You can then open a browser to [http://localhost:8080](http://localhost:8080) to see results. Consult the Maelstrom [result documentation](https://github.com/jepsen-io/maelstrom/blob/main/doc/results.md) for further details.
