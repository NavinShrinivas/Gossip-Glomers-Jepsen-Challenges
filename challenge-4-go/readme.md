# Challenge #4: Grow-Only Counter

In this challenge, you’ll need to implement a stateless, grow-only counter which will run against Maelstrom’s [`g-counter` workload](https://github.com/jepsen-io/maelstrom/blob/main/doc/workloads.md#workload-g-counter). This challenge is different than before in that your nodes will rely on a [sequentially-consistent](https://jepsen.io/consistency/models/sequential) key/value store service provided by Maelstrom.

## [](https://fly.io/dist-sys/4//#specification)Specification

Your node will need to accept two RPC-style message types: `add` & `read`. Your service need only be eventually consistent: given a few seconds without writes, it should converge on the correct counter value.

Please note that the final read from each node should return the final & correct count.

### [](https://fly.io/dist-sys/4//#rpc-add)RPC: `add`

Your node should accept `add` requests and increment the value of a single global counter. Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"add"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"delta"</span><span class="p">:</span><span class="w"> </span><span class="mi">123</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


and it will need to return an `"add_ok"` acknowledgement message:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"add_ok"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


### [](https://fly.io/dist-sys/4//#rpc-read)RPC: `read`

Your node should accept `read` requests and return the current value of the global counter. Remember that the counter service is only [sequentially consistent](https://jepsen.io/consistency/models/sequential). Your node will receive a request message body that looks like this:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"read"</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


and it will need to return a `"read_ok"` message with the current value:

Unwrap text Copy to clipboard


<span class="p">{</span><span class="w">
  </span><span class="nl">"type"</span><span class="p">:</span><span class="w"> </span><span class="s2">"read_ok"</span><span class="p">,</span><span class="w">
  </span><span class="nl">"value"</span><span class="p">:</span><span class="w"> </span><span class="mi">1234</span><span class="w">
</span><span class="p">}</span><span class="w">
</span>


### [](https://fly.io/dist-sys/4//#service-seq-kv)Service: `seq-kv`

Maelstrom provides a [sequentially-consistent](https://jepsen.io/consistency/models/sequential) key/value store called [`seq-kv`](https://github.com/jepsen-io/maelstrom/blob/main/doc/services.md#seq-kv) which has `read`, `write`, & `cas` operations. The Go library provides a [`KV`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#KV) wrapper for this service that you can instantiate with [`NewSeqKV()`](https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#NewSeqKV):

Unwrap text Copy to clipboard


<span class="n">node</span> <span class="o">:=</span> <span class="n">maelstrom</span><span class="o">.</span><span class="n">NewNode</span><span class="p">()</span>
<span class="n">kv</span> <span class="o">:=</span> <span class="n">maelstrom</span><span class="o">.</span><span class="n">NewSeqKV</span><span class="p">(</span><span class="n">node</span><span class="p">)</span>


The API is as follows:

Unwrap text Copy to clipboard


func (kv *KV) Read(ctx context.Context, key string) (any, error)
    Read returns the value for a given key in the key/value store. Returns an
    *RPCError error with a KeyDoesNotExist code if the key does not exist.

func (kv *KV) ReadInt(ctx context.Context, key string) (int, error)
    ReadInt reads the value of a key in the key/value store as an int.

func (kv *KV) Write(ctx context.Context, key string, value any) error
    Write overwrites the value for a given key in the key/value store.

func (kv *KV) CompareAndSwap(ctx context.Context, key string, from, to any, createIfNotExists bool) error
    CompareAndSwap updates the value for a key if its current value matches the
    previous value. Creates the key if createIfNotExists is true.

    Returns an *RPCError with a code of PreconditionFailed if the previous value
    does not match. Return a code of KeyDoesNotExist if the key did not exist.


## [](https://fly.io/dist-sys/4//#evaluation)Evaluation

Build your Go binary as `maelstrom-counter` and run it against Maelstrom with the following command:

Unwrap text Copy to clipboard


./maelstrom test -w g-counter --bin ~/go/bin/maelstrom-counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition


This will run a 3-node cluster for 20 seconds and increment the counter at the rate of 100 requests per second. It will induce network partitions during the test.

If you’re successful, right on! Continue on to the [Kafka-Style Log challenge](https://fly.io/dist-sys/5a). If you’re having trouble, ask for help on the [Fly.io Community forum](https://community.fly.io/).
