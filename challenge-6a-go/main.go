// Challenge 6A Single node KV store
package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var local_map map[float64]float64
var mu *sync.Mutex

func main() {
	n := maelstrom.NewNode()
	local_map = make(map[float64]float64)
	mu = new(sync.Mutex)

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
			} else {
				//read
            mu.Lock()
				value, ok := local_map[key]
            mu.Unlock()
				if !ok {
					output_txns = append(output_txns, []interface{}{"r", key, nil})
				} else {
					output_txns = append(output_txns, []interface{}{"r", key, value})
				}
			}
		}

		resp["type"] = "txn_ok"
		resp["txn"] = output_txns
		resply := n.Reply(msg, resp)
		return resply
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
