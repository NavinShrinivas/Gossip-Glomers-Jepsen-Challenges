package main

//Sequentially consistenet means, it appears to execute events in the exact order they appear in the service

import (
	"context"
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var local_gcounter float64
var kv *maelstrom.KV

func update_global_kv_record(n *maelstrom.Node) {
	for {
		if n.ID() == "" {
			continue
		}
		value, err := kv.Read(context.Background(), n.ID())
		if err != nil && maelstrom.ErrorCode(err) != 20 { //20 means the key does not exist
			log.Panic(err)
		} else {
			if value != nil && int(local_gcounter) > value.(int) {
				kv.Write(context.Background(), n.ID(), local_gcounter)
			} else if value == nil {
				kv.Write(context.Background(), n.ID(), local_gcounter)
			}
		}
	}
}

func get_sum_of_values(n *maelstrom.Node) float64 {
	value_sum := 0.0
	for _, v := range n.NodeIDs() {
		value, err := kv.Read(context.Background(), v)
		if err != nil {
			log.Println(err)
		} else {
			value_sum += float64(value.(int))
		}
	}
	return float64(value_sum)
}

func main() {
	n := maelstrom.NewNode()
	kv = maelstrom.NewSeqKV(n)
	local_gcounter = 0
	log.Println("Starting node...")


   //Shouldnt this be causing racy behaviour? 
   //We are updating the gloabl counters based on local counters, even before we are
   //sure if local counters are upto date. 
   //I think its depending on local counter being updated before 1st read comes in and 
   //global counter being updated before second read.
	go update_global_kv_record(n)
	n.Handle("add", func(msg maelstrom.Message) error {

		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		value := body["delta"].(float64)
		local_gcounter += value
		resp["type"] = "add_ok"
		resply := n.Reply(msg, resp)
		return resply
	})

	n.Handle("read", func(msg maelstrom.Message) error {

		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		resp["type"] = "read_ok"
		sum := get_sum_of_values(n)
		resp["value"] = sum
		resply := n.Reply(msg, resp)
		return resply
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
