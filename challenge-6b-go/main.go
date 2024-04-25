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
var writes_buffer_channel chan SyncWrite
var mu *sync.Mutex

type SyncWrite struct {
	nodeID string // Node to send this write to to be synced
	a, b   float64
}

func try_write(pair SyncWrite, writes_buffer_channel chan SyncWrite) {
	var body map[string]any = make(map[string]any)
	body["type"] = "txn"
	body["txn"] = [][]interface{}{{"w", pair.a, pair.b}}
	_, err := n.SyncRPC(context.Background(), pair.nodeID, body)
	if err != nil {
		//If we get any error, we put the work back into the queue
		writes_buffer_channel <- pair
	}
}

func one_way_write_syncer(writes_buffer_channel chan SyncWrite) {
	//Doesnt care about commit like one would in 2PL
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
	writes_buffer_channel = make(chan SyncWrite)
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

		for _, v := range txns {
			inner_tx := v.([]interface{})
			op := inner_tx[0].(string)
			key := inner_tx[1].(float64)
			var value float64
			if op == "w" {
				value = inner_tx[2].(float64)
            mu.Lock()
				local_map[key] = value
            mu.Unlock()
				output_txns = append(output_txns, []interface{}{"w", key, value})
				for _, i := range n.NodeIDs() {
               if i!=n.ID(){
                  //Do not send to self
                  syncwrite := SyncWrite{
                     nodeID: i,
                     a:      key,
                     b:      value,
                  }
                  writes_buffer_channel <- syncwrite
               }
				}
			} else {
				//read
            mu.Lock()
				value, ok := local_map[key]
            mu.Unlock()
				if  !ok{
					output_txns = append(output_txns, []interface{}{"r", key, nil})
				} else {
					output_txns = append(output_txns, []interface{}{"r", key, value})
				}
			}
		}

		// When we move to a read committed guarenteed, sending back "txn_ok" pretty much means we are commited the transactions (writes)
		resp["type"] = "txn_ok"
		resp["txn"] = output_txns
		resply := n.Reply(msg, resp)
		return resply
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
