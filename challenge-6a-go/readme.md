# Challenge #6a: Single-Node, Totally-Available Transactions

In this challenge, you’ll need to implement a key/value store which implements transactions. These transactions contain micro-operations (read & write) and the results of those operations depends on the consistency guarantees of the challenge. Your goal is to support weak consistency while also being totally available. We begin with a single-node service and then write a multi-node version.

## [](https://fly.io/dist-sys/6a//#specification)Specification

Your node will support the [`txn-rw-register`](https://github.com/jepsen-io/maelstrom/blob/main/doc/workloads.md#workload-txn-rw-register) workload by implementing a key/value store that accepts only one message. How easy, right?? This message is the `txn` message which passes in a list of operations to perform.

Writes in this workload are unique per-key so key `100` would only ever see a write of `1` once, a write of `2` once, etc. This helps Maelstrom to verify correctness.

### [](https://fly.io/dist-sys/6a//#rpc-txn)RPC: `txn`

This message passes in an array operations in the `"txn"` key. Each operation is represented by a 3-element array containing the operation name, the integer key to operate on, and a possibly-null integer value.

For example, your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"txn"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"msg_id"</span><span class="p">:</span><span class="w"> </span><span class="mi">3</span><span class="p">,</span><span class="w">
  </span><span class="nl">"txn"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="w">
    </span><span class="p">[</span><span class="s2">"r"</span><span class="p">,</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w"> </span><span class="kc">null</span><span class="p">],</span><span class="w">
    </span><span class="p">[</span><span class="s2">"w"</span><span class="p">,</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w"> </span><span class="mi">6</span><span class="p">],</span><span class="w">
    </span><span class="p">[</span><span class="s2">"w"</span><span class="p">,</span><span class="w"> </span><span class="mi">2</span><span class="p">,</span><span class="w"> </span><span class="mi">9</span><span class="p">]</span><span class="w">
  </span><span class="p">]</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


This represents three operations:

-   Read from key `1`
-   Write the value of `6` to key `1`
-   Write the value of `9` to key `2`

In response, it should send a `txn_ok` message that contains the same operation list, however, read (`"r"`) operations should have their value filled in with the current value. For example, if the value of key `1` was `3` before this transaction, it should be returned in the read operation. Non-existent keys should be returned as `null`.

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"txn_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"msg_id"</span><span class="p">:</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w">
  </span><span class="nl">"in_reply_to"</span><span class="p">:</span><span class="w"> </span><span class="mi">3</span><span class="p">,</span><span class="w">
  </span><span class="nl">"txn"</span><span class="p">:</span><span class="w"> </span><span class="p">[</span><span class="w">
    </span><span class="p">[</span><span class="s2">"r"</span><span class="p">,</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w"> </span><span class="mi">3</span><span class="p">],</span><span class="w">
    </span><span class="p">[</span><span class="s2">"w"</span><span class="p">,</span><span class="w"> </span><span class="mi">1</span><span class="p">,</span><span class="w"> </span><span class="mi">6</span><span class="p">],</span><span class="w">
    </span><span class="p">[</span><span class="s2">"w"</span><span class="p">,</span><span class="w"> </span><span class="mi">2</span><span class="p">,</span><span class="w"> </span><span class="mi">9</span><span class="p">]</span><span class="w">
  </span><span class="p">]</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


## [](https://fly.io/dist-sys/6a//#evaluation)Evaluation

Build your Go binary as `maelstrom-txn` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard


./maelstrom test -w txn-rw-register --bin ~/go/bin/maelstrom-txn --node-count 1 --time-limit 20 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total


This will verify your single-node system works before we move on to distributing our writes across nodes.

If you’re successful, continue on to the [Totally-Available, Read Uncommitted Transactions challenge](https://fly.io/dist-sys/6b). If you’re having trouble, jump over to the [Fly.io Community forum](https://community.fly.io/) for help.
