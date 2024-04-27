# Challenge #5b: Multi-Node Kafka-Style Log

In this challenge, you’ll need to take your [Single-Node Kafka system](https://fly.io/dist-sys/5a) and distribute it out to multiple nodes.

Your nodes can use the linearizable key/value store provided by Maelstrom to implement your distributed, replicated log. This challenge is about correctness and not efficiency. You only need to keep up with a reasonable request rate. It’s important to consider which components require [linearizability](https://jepsen.io/consistency/models/linearizable) versus [sequential consistency](https://jepsen.io/consistency/models/sequential).

## [](https://fly.io/dist-sys/5b//#specification)Specification

This challenge works the same as the single-node except that it’s now running with two nodes. All correctness checks in Maelstrom should pass.

### [](https://fly.io/dist-sys/5b//#service-lin-kv)Service: lin-kv

You’ve used the `seq-kv` service in the [Grow-only Counter challenge](https://fly.io/dist-sys/4), however, in this challenge you can use the linearizable version called [`lin-kv`](https://github.com/jepsen-io/maelstrom/blob/main/doc/services.md#lin-kv). The API is the same but they have different consistency guarantees.

You can instantiate the Go client in the library by using the [`NewLinKV()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#NewLinKV) function:

Unwrap text Copy to clipboard


<span class="n">node</span> <span class="o">:=</span> <span class="n">maelstrom</span><span class="o">.</span><span class="n">NewNode</span><span class="p">()</span>
<span class="n">kv</span> <span class="o">:=</span> <span class="n">maelstrom</span><span class="o">.</span><span class="n">NewLinKV</span><span class="p">(</span><span class="n">node</span><span class="p">)</span>


## [](https://fly.io/dist-sys/5b//#evaluation)Evaluation

Build your Go binary as `maelstrom-kafka` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard


./maelstrom test -w kafka --bin ~/go/bin/maelstrom-kafka --node-count 2 --concurrency 2n --time-limit 20 --rate 1000


This will run a two-node system for 20 seconds with 4 clients (`2n`). It will validate the system for correctness.

If you’re successful, that’s great! Continue on to the [Efficient Kafka challenge](https://fly.io/dist-sys/5c). If you’re having trouble, jump over to the [Fly.io Community forum](https://community.fly.io/) for help.
