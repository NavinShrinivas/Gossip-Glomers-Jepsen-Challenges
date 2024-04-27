# Challenge #3e: Efficient Broadcast, Part II

In this challenge, we’ll make our [Efficient, Multi-Node Broadcast](https://fly.io/dist-sys/3d) implementation even more efficient. Why settle for a fast distributed system when you could always make faster?

## [](https://fly.io/dist-sys/3e//#specification)Specification

With the same node count of `25` and a message delay of `100ms`, your challenge is to achieve the following performance metrics:

-   Messages-per-operation is below `20`
-   Median latency is below `1 second`
-   Maximum latency is below `2 seconds`

## [](https://fly.io/dist-sys/3e//#evaluation)Evaluation

Build your Go binary as `maelstrom-broadcast` and run it against Maelstrom with the same command as before:

Unwrap text Copy to clipboard


./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100


On success, congratulations, you’ve completed the Broadcast challenge. Move on to the [Grow-Only Counter challenge](https://fly.io/dist-sys/4). If you’re having trouble, poke your head in at [Fly.io Community forum](https://community.fly.io/) and ask for some help.
