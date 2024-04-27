(ns maelstrom.workload.txn-list-append
  "A transactional workload over a map of keys to lists of elements. Clients
  submit a single transaction via a `txn` request, and expect a
  completed version of that transaction in a `txn_ok` response.

  A transaction is an array of micro-operations, which should be executed in
  order:

  ```edn
  [op1, op2, ...]
  ```

  Each micro-op is a 3-element array comprising a function, key, and value:

  ```edn
  [f, k, v]
  ```

  There are two functions. A *read* observes the current value of a specific
  key. `[\"r\", 5, [1, 2]]` denotes that a read of key 5 observed the list `[1,
  2]`. When clients submit reads, they leave their values `null`:  `[\"r\", 5,
  null]`. The server processing the transaction should replace that value with
  whatever the observed value is for that key: `[\"r\", 5, [1, 2]]`.

  An *append* adds an element to the end of the key's current value. For
  instance, `[\"append\", 5, 3]` means \"add 3 to the end of the list for key
  5.\" If key 5 were currently `[1, 2]`, the resulting value would become `[1,
  2, 3]`. Appends have values provided by the client, and are returned
  unchanged.

  For example, assume the current state of the database is `{1 [8]}`, and you
  receive a request body like:

  ```json
  {\"type\": \"txn\",
   \"msg_id\": 1,
   \"txn\": [[\"r\", 1, null], [\"append\", 1, 6], [\"append\", 2, 9]]}
  ```

  You might return a response like:

  ```json
  {\"type\": \"txn_ok\",
   \"in_reply_to\": 1,
   \"txn\": [[\"r\", 1, [8]], [\"append\", 1, 6], [\"append\", 2, 9]]}
  ```

  First you read the current value of key 1, returning the list [8]. Then you
  append 6 to key 1. Then you append 9 to key 2, implicitly creating it. The
  resulting state of the database would be `{1 [8, 6], 2 [9]}`.

  Appends in this workload are always integers, and are unique per key. Key
  `x` will only ever see at most one append of `0`, at most one append of `1`,
  and so on.

  Unlike lin-kv, nonexistent keys should be returned as `null`. Lists are
  implicitly created on first append.

  This workload can check many kinds of consistency models. See the
  `--consistency-models` CLI option for details."
  (:refer-clojure :exclude [read])
  (:require [elle.core :as elle]
            [maelstrom [client :as c]
                       [net :as net]]
            [jepsen [client :as client]
                    [generator :as gen]
                    [independent :as independent]]
            [jepsen.tests.cycle.append :as append]
            [schema.core :as s]))

(def Key     s/Any)
(def Element s/Any)

(def ReadReq [(s/one (s/eq "r") "f") (s/one Key "k") (s/one (s/eq nil) "v")])
(def ReadRes [(s/one (s/eq "r") "f") (s/one Key "k") (s/one [Element] "v")])
(def Append  [(s/one (s/eq "append") "f") (s/one Key "k") (s/one Element "v")])

(c/defrpc txn!
  "Requests that the node execute a single transaction. Servers respond with a
  `txn_ok` message, and a completed version of the requested transaction; e.g.
  with read values filled in. Keys and list elements may be of any type."
  {:type  (s/eq "txn")
   :txn   [(s/either ReadReq Append)]}
  {:type  (s/eq "txn_ok")
   :txn   [(s/either ReadRes Append)]})

(defn kw->str
  "We use keywords for our :f's. Converts keywords to strings in a txn."
  [txn]
  (mapv (fn [[f k v]]
          [(name f) k v])
        txn))

(defn str->kw
  "We use keywords for our :f's. Converts strings to keywords in a txn."
  [txn]
  (mapv (fn [[f k v]]
          [(keyword f) k v])
        txn))

(defn client
  "Construct a linearizable key-value client for the given network"
  ([net]
   (client net nil nil))
  ([net conn node]
   (reify client/Client
     (open! [this test node]
       (client net (c/open! net) node))

     (setup! [this test])

     (invoke! [_ test op]
       (c/with-errors op #{}
         (->> (:value op)
              kw->str
              (array-map :txn)
              (txn! conn node)
              :txn
              str->kw
              (assoc op :type :ok, :value))))

     (teardown! [_ test])

     (close! [_ test]
       (c/close! conn))

     client/Reusable
     (reusable? [this test]
       true))))

(defn workload
  "Constructs a workload for linearizable registers, given option from the CLI
  test constructor. Options are:

    :net                  A Maelstrom network
    :key-count            Number of keys to work on at any single time
    :min-txn-length       Minimum number of ops per transaction
    :max-txn-length       Maximum number of ops per transaction
    :max-writes-per-key   How many elements to (try) appending to a single key.
    :consistency-models   What kinds of consistency models to check for."
  [opts]
  (-> (append/test opts)
      (assoc :client (client (:net opts)))))
