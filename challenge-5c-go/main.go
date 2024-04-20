// Challenge 5B Multi Node Kafka
package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

//There are no recency requirements so acknowledged send messages do not need to return in poll messages immediately.

//We can avoid CAS if we simply handle all keys in a single node.
//Using a good hash function, we can send all same keys to the same node

var kv *maelstrom.KV
var kv2 *maelstrom.KV

func main() {
	n := maelstrom.NewNode()
	kv = maelstrom.NewLinKV(n)
	kv2 = maelstrom.NewSeqKV(n)
	n.Handle("send", func(msg maelstrom.Message) error {
		number_of_nodes := len(n.NodeIDs())
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		key := body["key"].(string)
		value := body["msg"].(float64)

		//To find onwership :
		vali, _ := strconv.Atoi(key)
		key_ownership_node := int(vali % number_of_nodes)
		if n.NodeIDs()[key_ownership_node] != n.ID() {
			ret_msg, err := n.SyncRPC(context.Background(),n.NodeIDs()[key_ownership_node], msg.Body)
			if err != nil {
				log.Println(err)
				return err
			}
			if err := json.Unmarshal(ret_msg.Body, &body); err != nil {
				log.Println(err)
				return err
			}
			resp["type"] = "send_ok"
			resp["offset"] = body["offset"].(float64)
			resply := n.Reply(msg, resp)
			return resply

		}

		//Loop from :
		for {
			current_log, err := kv.Read(context.Background(), key)
			if err != nil && maelstrom.ErrorCode(err) == maelstrom.KeyDoesNotExist {
				possible_future_log_arr := []float64{value}
				//We need to use CAS here, unlike in challenge 4 as :
				//Unlike grow only counter, we cannot just overwrite (add) previous value, we need to save it
				err = kv.Write(context.Background(), key, possible_future_log_arr)
				if err != nil {
					continue
				}
				resp["type"] = "send_ok"
				resp["offset"] = 0
				resply := n.Reply(msg, resp)
				return resply
			} else if err != nil {
				log.Panic(err)
				return err
			} else {
				current_log_arr := current_log.([]interface{})
				possible_future_log_arr := append(current_log_arr, value)
				//We need to use CAS here, unlike in challenge 4 as :
				//Unlike grow only counter, we cannot just overwrite (add) previous value, we need to save it
				err = kv.Write(context.Background(), key, possible_future_log_arr)
				if err != nil {
					continue //Re try
				}
				resp["type"] = "send_ok"
				resp["offset"] = len(current_log_arr)
				resply := n.Reply(msg, resp)
				return resply
			}
		}
	})

	n.Handle("poll", func(msg maelstrom.Message) error {
		//You dont need any Sequential guarentees here :)
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		var offseted_log_output map[string][][]float64 = make(map[string][][]float64)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		key_offsets := body["offsets"].(map[string]interface{})
		for k, v := range key_offsets {
			current_log, err := kv.Read(context.Background(), k)
			if err != nil && maelstrom.ErrorCode(err) == maelstrom.KeyDoesNotExist {
				continue //If key doesn't exists, we can avoid
			} else if err != nil {
				continue
			} else {
				current_log_arr := current_log.([]interface{})
				temp := [][]float64{}
				for i, v_i := range current_log_arr {
					if float64(i) < v.(float64) {
						continue
					} else {
						temp = append(temp, []float64{float64(i), v_i.(float64)})
					}
				}
				offseted_log_output[k] = temp
			}
		}
		resp["type"] = "poll_ok"
		resp["msgs"] = offseted_log_output
		resply := n.Reply(msg, resp)
		return resply
	})

	n.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		commit_key_offsets := body["offsets"].(map[string]interface{})
		for k, v := range commit_key_offsets {
			kv2.Write(context.Background(), k, v.(float64))
		}
		resp["type"] = "commit_offsets_ok"
		resply := n.Reply(msg, resp)
		return resply
	})

	n.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		var resp map[string]any = make(map[string]any)
		var commited_offset_output map[string]float64 = make(map[string]float64)
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Println(err)
			return err
		}
		log.Println(body)
		keys := body["keys"].([]interface{})
		for _, v := range keys {
			val, err := kv2.Read(context.Background(), v.(string))
			if err != nil {
				continue
			} else {
				commited_offset_output[v.(string)] = float64(val.(int))
			}
		}
		resp["type"] = "list_committed_offsets_ok"
		resp["offsets"] = commited_offset_output
		resply := n.Reply(msg, resp)
		return resply
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
