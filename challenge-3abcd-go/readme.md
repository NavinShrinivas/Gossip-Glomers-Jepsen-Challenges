# Challenge #3a: Single-Node Broadcast

In this challenge, you’ll need to implement a broadcast system that gossips messages between all nodes in the cluster. Gossiping is a common way to propagate information across a cluster when you don’t need strong consistency guarantees.

This challenge is broken up in multiple sections so that you can build out your system incrementally. First, we’ll start out with a single-node broadcast system. That may sound like an oxymoron but this lets us get our message handlers working correctly in isolation before trying to share messages between nodes.

## [](https://fly.io/dist-sys/3a//#specification)Specification

Your node will need to handle the [`"broadcast"`](https://github.com/jepsen-io/maelstrom/blob/main/doc/workloads.md#workload-broadcast) workload which has 3 RPC message types: `broadcast`, `read`, & `topology`. Your node will need to store the set of integer values that it sees from `broadcast` messages so that they can be returned later via the `read` message RPC.

The Go library has two methods for sending messages:

1.  `Send()` sends a fire-and-forget message and doesn’t expect a response. As such, it does not attach a message ID.
    
2.  `RPC()` sends a message and accepts a response handler. The message will be decorated with a message ID so the handler can be invoked when a response message is received.
    

Data can be stored in-memory as node processes are not killed by Maelstrom.

### [](https://fly.io/dist-sys/3a//#rpc-broadcast)RPC: `broadcast`

This message requests that a value be broadcast out to all nodes in the cluster. The value is always an integer and it is unique for each message from Maelstrom.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"broadcast"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"message"</span><span class="p">:</span><span class="w"> </span><span class="mi">1000</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

It should store the `"message"` value locally so it can be read later. In response, it should send an acknowledge with a `broadcast_ok` message:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"broadcast_ok"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

### [](https://fly.io/dist-sys/3a//#rpc-read)RPC: `read`

This message requests that a node return all values that it has seen.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"read"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

In response, it should return a `read_ok` message with a list of values it has seen:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"read_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"messages"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="mi">1</span><span class="p">,</span><span class="w"> </span><span class="mi">8</span><span class="p">,</span><span class="w"> </span><span class="mi">72</span><span class="p">,</span><span class="w"> </span><span class="mi">25</span><span class="p">]</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

The order of the returned values does not matter.

### [](https://fly.io/dist-sys/3a//#rpc-topology)RPC: `topology`

This message informs the node of who its neighboring nodes are. Maelstrom has multiple topologies available or you can ignore this message and make your own topology from the list of nodes in the [`Node.NodeIDs()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#Node.NodeIDs) method. All nodes can communicate with each other regardless of the topology passed in.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"topology"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"topology"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nl">"n1"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"n2"</span><span class="p">,</span><span class="w"> </span><span class="s2">"n3"</span><span class="p">],</span><span class="w">
    </span><span class="nl">"n2"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"n1"</span><span class="p">],</span><span class="w">
    </span><span class="nl">"n3"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"n1"</span><span class="p">]</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

In response, your node should return a `topology_ok` message body:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"topology_ok"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

## [](https://fly.io/dist-sys/3a//#evaluation)Evaluation

Build your Go binary as `maelstrom-broadcast` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard

```
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10
```

This will run a single node for 20 seconds and send messages at the rate of 10 messages per second. It will validate that all values sent by broadcasts are returned via read.

If you’re successful, that’s awesome! Continue on to the [Multi-Node Broadcast challenge](https://fly.io/dist-sys/3b). If you’re having trouble, jump over to the [Fly.io Community forum](https://community.fly.io/) for help.

# Challenge #3a: Single-Node Broadcast

In this challenge, you’ll need to implement a broadcast system that gossips messages between all nodes in the cluster. Gossiping is a common way to propagate information across a cluster when you don’t need strong consistency guarantees.

This challenge is broken up in multiple sections so that you can build out your system incrementally. First, we’ll start out with a single-node broadcast system. That may sound like an oxymoron but this lets us get our message handlers working correctly in isolation before trying to share messages between nodes.

## [](https://fly.io/dist-sys/3a//#specification)Specification

Your node will need to handle the [`"broadcast"`](https://github.com/jepsen-io/maelstrom/blob/main/doc/workloads.md#workload-broadcast) workload which has 3 RPC message types: `broadcast`, `read`, & `topology`. Your node will need to store the set of integer values that it sees from `broadcast` messages so that they can be returned later via the `read` message RPC.

The Go library has two methods for sending messages:

1.  `Send()` sends a fire-and-forget message and doesn’t expect a response. As such, it does not attach a message ID.
    
2.  `RPC()` sends a message and accepts a response handler. The message will be decorated with a message ID so the handler can be invoked when a response message is received.
    

Data can be stored in-memory as node processes are not killed by Maelstrom.

### [](https://fly.io/dist-sys/3a//#rpc-broadcast)RPC: `broadcast`

This message requests that a value be broadcast out to all nodes in the cluster. The value is always an integer and it is unique for each message from Maelstrom.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"broadcast"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"message"</span><span class="p">:</span><span class="w"> </span><span class="mi">1000</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

It should store the `"message"` value locally so it can be read later. In response, it should send an acknowledge with a `broadcast_ok` message:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"broadcast_ok"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

### [](https://fly.io/dist-sys/3a//#rpc-read)RPC: `read`

This message requests that a node return all values that it has seen.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"read"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

In response, it should return a `read_ok` message with a list of values it has seen:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"read_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"messages"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="mi">1</span><span class="p">,</span><span class="w"> </span><span class="mi">8</span><span class="p">,</span><span class="w"> </span><span class="mi">72</span><span class="p">,</span><span class="w"> </span><span class="mi">25</span><span class="p">]</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

The order of the returned values does not matter.

### [](https://fly.io/dist-sys/3a//#rpc-topology)RPC: `topology`

This message informs the node of who its neighboring nodes are. Maelstrom has multiple topologies available or you can ignore this message and make your own topology from the list of nodes in the [`Node.NodeIDs()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#Node.NodeIDs) method. All nodes can communicate with each other regardless of the topology passed in.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"topology"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"topology"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nl">"n1"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"n2"</span><span class="p">,</span><span class="w"> </span><span class="s2">"n3"</span><span class="p">],</span><span class="w">
    </span><span class="nl">"n2"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"n1"</span><span class="p">],</span><span class="w">
    </span><span class="nl">"n3"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"n1"</span><span class="p">]</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

In response, your node should return a `topology_ok` message body:

Unwrap text Copy to clipboard

<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"topology_ok"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>

## [](https://fly.io/dist-sys/3a//#evaluation)Evaluation

Build your Go binary as `maelstrom-broadcast` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard

```
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10
```

This will run a single node for 20 seconds and send messages at the rate of 10 messages per second. It will validate that all values sent by broadcasts are returned via read.

If you’re successful, that’s awesome! Continue on to the [Multi-Node Broadcast challenge](https://fly.io/dist-sys/3b). If you’re having trouble, jump over to the [Fly.io Community forum](https://community.fly.io/) for help.

# Challenge #3c: Fault Tolerant Broadcast

In this challenge, we’ll build on our [Multi-Node Broadcast](https://fly.io/dist-sys/3b) implementation, however, this time we’ll introduce network partitions between nodes so they will not be able to communicate for periods of time.

## [](https://fly.io/dist-sys/3c//#specification)Specification

Your node should propagate values it sees from `broadcast` messages to the other nodes in the cluster—even in the face of network partitions! Values should propagate to all other nodes by the end of the test. Nodes should only return copies of their own local values.

## [](https://fly.io/dist-sys/3c//#evaluation)Evaluation

Build your Go binary as `maelstrom-broadcast` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard

```
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition
```

This will run a 5-node cluster like before, but this time with a failing network! Fun!

On success, continue on to [Part One of the Broadcast Efficiency challenge](https://fly.io/dist-sys/3d). If you’re having trouble, head to the [Fly.io Community forum](https://community.fly.io/).

# Challenge #3d: Efficient Broadcast, Part I

In this challenge, we’ll improve on our [Fault-Tolerant, Multi-Node Broadcast](https://fly.io/dist-sys/3c) implementation. Distributed systems have different metrics for success. Not only do they need to be correct but they also need to be fast.

The neighbors Maelstrom suggests are, by default, arranged in a two-dimensional grid. This means that messages are often duplicated en route to other nodes, and latencies are on the order of `2 * sqrt(n)` network delays.

## [](https://fly.io/dist-sys/3d//#specification)Specification

We will increase our node count to `25` and add a delay of `100ms` to each message to simulate a slow network. This could be geographic latencies (such as US to Europe) or it could simply be a busy network.

Your challenge is to achieve the following:

-   Messages-per-operation is below `30`
-   Median latency is below `400ms`
-   Maximum latency is below `600ms`

Feel free to ignore the topology you’re given by Maelstrom and use your own; it’s only a suggestion. Don’t compromise safety under faults. Double-check that your solution is still correct (even though it will be much slower) with `--nemesis partition`

### [](https://fly.io/dist-sys/3d//#messages-per-operation)Messages-per-operation

In the `results.edn` file produced by Maelstrom, you’ll find a `:net` key with information about the number of network messages. The `:servers` key shows just messages between server nodes, and `:msgs-per-op` shows the number of messages exchanged per logical operation. Almost all our operations are `broadcast` or `read`, in a 50/50 mix.

Unwrap text Copy to clipboard

```
:net {:all {:send-count 129592,
            :recv-count 129592,
            :msg-count 129592,
            :msgs-per-op 65.121605},
      :clients {:send-count 4080, :recv-count 4080, :msg-count 4080},
      :servers {:send-count 125512,
                :recv-count 125512,
                :msg-count 125512,
                :msgs-per-op 63.071358}
```

In this example we exchanged 63 messages per operation. Half of those are reads, which require no inter-server messages. That means we sent on average 126 messages per broadcast, between 25 nodes: roughly five messages per node.

### [](https://fly.io/dist-sys/3d//#stable-latencies)Stable latencies

Under `:workload` you’ll find a map of `:stable-latencies`. These are quantiles which show the broadcast latency for the minimum, median, 95th, 99th, and maximum latency request. These latencies are measured from the time a broadcast request was acknowledged to when it was last missing from a read on any node. For example, here’s a system whose median latency was 452 milliseconds:

Unwrap text Copy to clipboard

```
:stable-latencies {0 0,
                   0.5 452,
                   0.95 674,
                   0.99 731,
                   1 794},
```

## [](https://fly.io/dist-sys/3d//#evaluation)Evaluation

Build your Go binary as `maelstrom-broadcast` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard

```
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100
```

You can run `maelstrom serve` to view results or you can locate your most recent run in the `./store` directory.

On success, continue on to [Part Two of the Broadcast Efficiency challenge](https://fly.io/dist-sys/3e). If you’re having trouble, head to the [Fly.io Community forum](https://community.fly.io/).
