## Overview : 
This repo holds my solutions to the gossip glomers Challenges that has a series of challenges built using `maelstrom` sdk (go and rust). These are a bunch of introductory distributed system challenges that I think give a good first look into dist sys. I personally found them not so challenging.\

Link to these challenges (I hope they never take these down, if they do, the challenges sit in the readme of each solution) : [Gossip Glomer Challenges](https://fly.io/dist-sys/)

## Notes on Challenges : 

### Challenge-1 : Echo

This is a rather simple challenge of just echoing back all the messaged recived by a given node. This challenge only exists to get you used to the maelstrom SDK and runner. I wrote this solution in Rust, back when there was not great support for rust maelstrom SKD.
I figured it was making these challenges much harder not in the way these challenges intended, hence I moved to Golang from the next challenge.

### Challenge-2 : 

### Challenge-3a : 
### Challenge-3b : 
### Challenge-3c : 
### Challenge-3d : 
### Challenge-3e :
### Challenge-4 :
### Challenge-5a :
### Challenge-5b :
### Challenge-5c :
### Challenge-6a :
### Challenge-6b :

extension of challenge 6a, but now we move a multi node replicated kV store with a rather weak gaurentee. `Read Uncommitted` is what is asked of us. In this gaurentee almost anything is allowed. They also want you to check against partition for total availability.

tests 1: 

```
./maelstrom test -w txn-rw-register --bin ~/go/bin/maelstrom-txn --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted
```
test 2:
```
./maelstrom test -w txn-rw-register --bin ~/go/bin/maelstrom-txn --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition
```

### Challenge-6c :

This is an extension of challenge 6b. Now they still want a multi node KV store, but a `read committed` transactional gaurentee is required. Despite only being a small addition, this makes the challenge hard and worthy of a final `boss` challenge (xD)

The solution is rather simple, the key is [BEWARE BEFORE OPENING, its worth working out on your own]: 

<details>
  <summary>Spoiler warning</summary>

  Jespen asks for total availability, this is only possible if stale reads are allowed. In this challenge Jespen accepts all stale read. To handle intermediated reads make sure to get one lock per transaction and not break up the transaction.
  
  
  ```javascript
  console.log("I'm a code block!");
  ```
  ![image](./6c.jpeg)
  
</details>

test :

```
./maelstrom test -w txn-rw-register --bin ~/go/bin/maelstrom-txn --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-committed --availability total –-nemesis partition
```




### Notes on topics : 

#### Sequential consistency 
Sequential consistency is a strong safety property for concurrent systems. Informally, sequential consistency implies that operations appear to take place in some total order, and that that order is consistent with the order of operations on each individual process.

> Sequential consistency cannot provide gaurentees of ordering across processes

#### Linearizability


#### Read Uncommitted 

#### Read committed
