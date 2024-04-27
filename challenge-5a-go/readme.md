# Challenge #5a: Single-Node Kafka-Style Log

In this challenge, you’ll need to implement a replicated log service similar to [Kafka](https://kafka.apache.org/). Replicated logs are often used as a message bus or an event stream.

This challenge is broken up in multiple sections so that you can build out your system incrementally. First, we’ll start out with a single-node log system and then we’ll distribute it in later challenges.

## [](https://fly.io/dist-sys/5a//#specification)Specification

Your nodes will need to store an append-only log in order to handle the [`"kafka"`](https://github.com/jepsen-io/maelstrom/blob/main/doc/workloads.md#workload-kafka) workload. Each log is identified by a string key (e.g. `"k1"`) and these logs contain a series of messages which are identified by an integer offset. These offsets can be sparse in that not every offset must contain a message.

Maelstrom will check to make sure several anomalies do not occur:

-   _Lost writes:_ for example, a client sees offset 10 but not offset 5.
-   _Monotonic increasing offsets:_ an offset for a log should always be increasing.

There are no recency requirements so acknowledged `send` messages do not need to return in `poll` messages immediately.

### [](https://fly.io/dist-sys/5a//#rpc-send)RPC: `send`

This message requests that a `"msg"` value be appended to a log identified by `"key"`. Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"send"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"key"</span><span class="p">:</span><span class="w"> </span><span class="s2">"k1"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"msg"</span><span class="p">:</span><span class="w"> </span><span class="mi">123</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


In response, it should send an acknowledge with a `send_ok` message that contains the unique offset for the message in the log:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"send_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"offset"</span><span class="p">:</span><span class="w"> </span><span class="mi">1000</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


### [](https://fly.io/dist-sys/5a//#rpc-poll)RPC: `poll`

This message requests that a node return messages from a set of logs starting from the given offset in each log. Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"poll"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"offsets"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nl">"k1"</span><span class="p">:</span><span class="w"> </span><span class="mi">1000</span><span class="p">,</span><span class="w">
    </span><span class="nl">"k2"</span><span class="p">:</span><span class="w"> </span><span class="mi">2000</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


In response, it should return a `poll_ok` message with messages starting from the given offset for each log. Your server can choose to return as many messages for each log as it chooses:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"poll_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"msgs"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nl">"k1"</span><span class="p">:</span><span class="w"> </span><span class="p">[[</span><span class="mi">1000</span><span class="p">,</span><span class="w"> </span><span class="mi">9</span><span class="p">],</span><span class="w"> </span><span class="p">[</span><span class="mi">1001</span><span class="p">,</span><span class="w"> </span><span class="mi">5</span><span class="p">],</span><span class="w"> </span><span class="p">[</span><span class="mi">1002</span><span class="p">,</span><span class="w"> </span><span class="mi">15</span><span class="p">]],</span><span class="w">
    </span><span class="nl">"k2"</span><span class="p">:</span><span class="w"> </span><span class="p">[[</span><span class="mi">2000</span><span class="p">,</span><span class="w"> </span><span class="mi">7</span><span class="p">],</span><span class="w"> </span><span class="p">[</span><span class="mi">2001</span><span class="p">,</span><span class="w"> </span><span class="mi">2</span><span class="p">]]</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


### [](https://fly.io/dist-sys/5a//#rpc-commit_offsets)RPC: `commit_offsets`

This message informs the node that messages have been successfully processed up to and including the given offset. Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"commit_offsets"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"offsets"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nl">"k1"</span><span class="p">:</span><span class="w"> </span><span class="mi">1000</span><span class="p">,</span><span class="w">
    </span><span class="nl">"k2"</span><span class="p">:</span><span class="w"> </span><span class="mi">2000</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


In this example, the messages have been processed up to and including offset `1000` for log `k1` and all messages up to and including offset `2000` for `k2`.

In response, your node should return a `commit_offsets_ok` message body to acknowledge the request:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"commit_offsets_ok"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


### [](https://fly.io/dist-sys/5a//#rpc-list_committed_offsets)RPC: `list_committed_offsets`

This message returns a map of committed offsets for a given set of logs. Clients use this to figure out where to start consuming from in a given log.

Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"list_committed_offsets"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"keys"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="s2">"k1"</span><span class="p">,</span><span class="w"> </span><span class="s2">"k2"</span><span class="p">]</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


In response, your node should return a `list_committed_offsets_ok` message body containing a map of offsets for each requested key. Keys that do not exist on the node can be omitted.

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"list_committed_offsets_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"offsets"</span><span class="p">:</span><span class="w"> </span><span class="p">{</span><span class="w">
    </span><span class="nl">"k1"</span><span class="p">:</span><span class="w"> </span><span class="mi">1000</span><span class="p">,</span><span class="w">
    </span><span class="nl">"k2"</span><span class="p">:</span><span class="w"> </span><span class="mi">2000</span><span class="w">
  </span><span class="p">}</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


## [](https://fly.io/dist-sys/5a//#evaluation)Evaluation

Build your Go binary as `maelstrom-kafka` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard


./maelstrom test -w kafka --bin ~/go/bin/maelstrom-kafka --node-count 1 --concurrency 2n --time-limit 20 --rate 1000


This will run a single node for 20 seconds with two clients. It will validate that messages are queued and committed properly.

If you’re successful, wahoo! Continue on to the [Multi-Node Kafka challenge](https://fly.io/dist-sys/5b). If you’re having trouble, jump over to the [Fly.io Community forum](https://community.fly.io/) for help.
