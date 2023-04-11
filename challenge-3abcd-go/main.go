package main

import (
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var messages []float64

func appendMessages(msg_chan chan float64) {
	for {
		elem := <-msg_chan
		// if elem == 0 {
		// 	continue
		// }
		messages = append(messages, elem)
	}
}

var topology []interface{}
var recv []float64

func infinite_retry(body map[string]any, dest string, n *maelstrom.Node) {
	chan_done := false
	for {
		n.RPC(dest, body, func(msg maelstrom.Message) error {
			log.Println("Sent and recieved ack for internal_broadcast", msg)
			chan_done = true
			return nil
		})
		time.Sleep(2 * time.Second)
		if chan_done {
			break
		}
	}
}

func main() {
	n := maelstrom.NewNode()
	log.Println("Starting node...")
	msg_chan := make(chan float64, 10000)
	go appendMessages(msg_chan)
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		value := body["message"].(float64)
		msg_chan <- value
		recv = append(recv, value)

		// Update the message type to return back.
		body["type"] = "internal_broadcast"
		resp["type"] = "broadcast_ok"
		resp["msg_id"] = body["msg_id"]
		resply := n.Reply(msg, resp)
		for _, v := range n.NodeIDs() { //Not efficient, just blasts to all nodes
			if v == n.ID() {
				continue
			}
			go infinite_retry(body, v, n)
		}
		// Echo the original message back with the updated message type.
		return resply
	})

	n.Handle("internal_broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		value := body["message"].(float64)
		msg_chan <- value

		// Update the message type to return back.
		resp["type"] = "internal_broadcast_ok"
		resp["msg_id"] = body["msg_id"]
		resply := n.Reply(msg, resp)
		return resply
	})
	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)

		// Update the message type to return back.
		body["type"] = "read_ok"
		body["messages"] = messages

		log.Println(body)
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)

		// Update the message type to return back.
		resp["type"] = "topology_ok"
		topology = body["topology"].(map[string]interface{})[n.ID()].([]interface{})

		log.Println(body)
		return n.Reply(msg, resp)
	})
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
