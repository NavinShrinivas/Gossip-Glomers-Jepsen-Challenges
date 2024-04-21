// Challenge 6A Single node KV store
package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

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

      for _, v := range txns{
         inner_tx := v.([]interface{})
         op := inner_tx[0].(string)
         key := inner_tx[1].(float64)
         var value float64
         if op == "w"{
            value = inner_tx[2].(float64)
         }
      }

		resp["type"] = "send_ok"
		resp["offset"] = 0
		resply := n.Reply(msg, resp)
		return resply
	})


	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
