// Challenge 6B Multi node read uncommited KV Store
// To achieve 1PL on one side broadcast is enough!
// As reads dont care about getting on "unrollbackable" commited writes
package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var local_map map[float64]float64
var n *maelstrom.Node
var writes_buffer_channel chan SyncWrites
var mu *sync.Mutex

type SyncWrites struct {
	nodeID string // Node to send this write to to be synced
	txn    [][]interface{}
}

func try_write(pair SyncWrites, writes_buffer_channel chan SyncWrites) {
	var body map[string]any = make(map[string]any)
	body["type"] = "txn"
	body["txn"] = pair.txn
	_, err := n.SyncRPC(context.Background(), pair.nodeID, body)
	if err != nil {
		//If we get any error, we put the work back into the queue
		writes_buffer_channel <- pair
	}
}

func one_way_write_syncer(writes_buffer_channel chan SyncWrites) {
	// Doesnt care about write being committed like one would in 2PL
	// Using a go routime spawning pattern for handling indefinite try outs
	for {
		write_pair := <-writes_buffer_channel
		go try_write(write_pair, writes_buffer_channel)
	}

}

func main() {
	n = maelstrom.NewNode()
	mu = new(sync.Mutex)
	local_map = make(map[float64]float64)
	writes_buffer_channel = make(chan SyncWrites)
	go one_way_write_syncer(writes_buffer_channel)
	n.Handle("txn", func(msg maelstrom.Message) error {
		// number_of_nodes := len(n.NodeIDs())
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		txns := body["txn"].([]interface{})
		output_txns := [][]interface{}{}


      //Not splitting writes and maintaining transatcion garnularity by grouping all writes to a single transaction
		sync_writes := [][]interface{}{}

		//Full transactions locks to handle G1c
		mu.Lock()
		for _, v := range txns {
			inner_tx := v.([]interface{})
			op := inner_tx[0].(string)
			key := inner_tx[1].(float64)
			var value float64
			if op == "w" {
				value = inner_tx[2].(float64)
				local_map[key] = value
				output_txns = append(output_txns, []interface{}{"w", key, value})
				sync_writes = append(sync_writes, inner_tx)
			} else {
				//read
				value, ok := local_map[key]
				if !ok {
					output_txns = append(output_txns, []interface{}{"r", key, nil})
				} else {
					output_txns = append(output_txns, []interface{}{"r", key, value})
				}
			}
		}
		mu.Unlock()
		for _, i := range n.NodeIDs() {
			if i != n.ID() {
				//Do not send to self
				syncwrite := SyncWrites{
					nodeID: i,
					txn:    sync_writes,
				}
				writes_buffer_channel <- syncwrite
			}
		}

		// When we move to a read committed guarenteed, sending back "txn_ok" pretty much means  the writes are commited
		resp["type"] = "txn_ok"
		resp["txn"] = output_txns
		resply := n.Reply(msg, resp)
		return resply
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
