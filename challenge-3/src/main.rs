use async_maelstrom::*;
use async_trait::async_trait;

use std::env;
use std::sync::atomic::AtomicU64;
use std::sync::Arc;

use log::{info, warn};
use tokio::spawn;

use serde::{Deserialize, Serialize};
use serde_json::json;
use serde_json::Value;
use uuid::Uuid;

use tokio::sync::mpsc;

struct ReadServer {
    args: Vec<String>,
    net: process::ProcNet<Read, ()>,
    id: Id,
    ids: Vec<Id>,
    msg_id: AtomicU64,
    recv: mpsc::UnboundedReceiver<i32>,
}

impl ReadServer {
    fn new(recv: mpsc::UnboundedReceiver<i32>) -> Self {
        ReadServer {
            args: Default::default(),
            net: Default::default(),
            id: Default::default(),
            ids: Default::default(),
            msg_id: Default::default(),
            recv,
        }
    }
}

#[derive(Deserialize, Serialize, Debug, Eq, PartialEq)]
#[serde(tag = "type")]
pub enum Read {
    #[serde(rename = "broadcast")]
    Read { msg_id: msg::MsgId },
    #[serde(rename = "broadcast_ok")]
    ReadOk {
        in_reply_to: msg::MsgId,
        msg_id: Option<msg::MsgId>,
        messages: Vec<i32>,
    },
}

#[async_trait]
impl process::Process<Read, ()> for ReadServer {
    fn init(
        &mut self,
        args: Vec<String>,
        net: process::ProcNet<Read, ()>,
        id: Id,
        ids: Vec<Id>,
        start_msg_id: msg::MsgId,
    ) {
        self.args = args;
        self.net = net;
        self.id = id;
        self.ids = ids;
        self.msg_id = AtomicU64::new(start_msg_id)
    }
    async fn run(&self) -> Status {
        let mut message: Vec<i32> = Vec::new();
        loop {
            match self.net.rxq.recv().await {
                Ok(full_message) => {
                    let msg_src = full_message.src.clone();
                    println!("{:?}", json!(full_message));
                    match full_message.body {
                        msg::Body::Workload(Read::Read { msg_id }) => {
                            loop {
                                match self.recv.recv().await {
                                    Some(n) => message.push(n),
                                    None => break,
                                }
                            }
                            self.net
                                .txq
                                .send(msg::Msg {
                                    src: self.id.clone(),
                                    dest: msg_src,
                                    body: msg::Body::Workload(Read::ReadOk {
                                        in_reply_to: msg_id,
                                        msg_id: Some(
                                            self.msg_id
                                                .fetch_add(1, std::sync::atomic::Ordering::SeqCst),
                                        ),
                                        messages: message.clone(),
                                    }),
                                })
                                .await?;
                        }
                        _ => {
                            return Ok(());
                        }
                    }
                }
                _ => return Ok(()),
            }
        }
    }
}

struct BroadcastServer {
    args: Vec<String>,
    net: process::ProcNet<Broadcast, ()>,
    id: Id,
    ids: Vec<Id>,
    msg_id: AtomicU64,
    sender: mpsc::UnboundedSender<i32>,
}

// impl Default for BroadcastServer {
//     fn default() -> Self {
//         BroadcastServer {
//             args: Default::default(),
//             net: Default::default(),
//             id: Default::default(),
//             ids: Default::default(),
//             msg_id: Default::default(),
//             sender: Default::default(),
//         }
//     }
// }

impl BroadcastServer {
    fn new(sender: mpsc::UnboundedSender<i32>) -> BroadcastServer {
        BroadcastServer {
            args: Default::default(),
            net: Default::default(),
            id: Default::default(),
            ids: Default::default(),
            msg_id: Default::default(),
            sender,
        }
    }
}

#[derive(Deserialize, Serialize, Debug, Eq, PartialEq)]
#[serde(tag = "type")]
pub enum Broadcast {
    #[serde(rename = "broadcast")]
    Broadcast { msg_id: msg::MsgId, message: i32 },
    #[serde(rename = "broadcast_ok")]
    BroadcastOk {
        in_reply_to: msg::MsgId,
        msg_id: Option<msg::MsgId>,
    },
}

#[async_trait]
impl process::Process<Broadcast, ()> for BroadcastServer {
    fn init(
        &mut self,
        args: Vec<String>,
        net: process::ProcNet<Broadcast, ()>,
        id: Id,
        ids: Vec<Id>,
        start_msg_id: msg::MsgId,
    ) {
        self.args = args;
        self.net = net;
        self.id = id;
        self.ids = ids;
        self.msg_id = AtomicU64::new(start_msg_id)
    }
    async fn run(&self) -> Status {
        loop {
            match self.net.rxq.recv().await {
                Ok(full_message) => {
                    let msg_src = full_message.src.clone();
                    println!("{:?}", json!(full_message));
                    match full_message.body {
                        msg::Body::Workload(Broadcast::Broadcast { msg_id, message }) => {
                            self.sender.send(message).unwrap();
                            self.net
                                .txq
                                .send(msg::Msg {
                                    src: self.id.clone(),
                                    dest: msg_src,
                                    body: msg::Body::Workload(Broadcast::BroadcastOk {
                                        in_reply_to: msg_id,
                                        msg_id: Some(
                                            self.msg_id
                                                .fetch_add(1, std::sync::atomic::Ordering::SeqCst),
                                        ),
                                    }),
                                })
                                .await?;
                        }
                        _ => {
                            return Ok(());
                        }
                    }
                }
                _ => return Ok(()),
            }
        }
    }
}

#[tokio::main]
async fn main() -> Status {
    let (tx, rx): (mpsc::UnboundedSender<i32>, mpsc::UnboundedReceiver<i32>) =
        mpsc::unbounded_channel();
    let broadcast_server_instance: BroadcastServer = BroadcastServer::new(tx);
    let read_server_instance: ReadServer = ReadServer::new(rx);
    let r =
        Arc::new(runtime::Runtime::new(env::args().collect(), broadcast_server_instance).await?);
    let reader =
        Arc::new(runtime::Runtime::new(env::args().collect(), read_server_instance).await?);
    let (r1, r2, r3) = (r.clone(), r.clone(), r.clone());
    let (reader1, reader2, reader3) = (reader.clone(), reader.clone(), reader.clone());
    let t1 = spawn(async move { r1.run_io_egress().await });
    let t2 = spawn(async move { r2.run_io_ingress().await });
    let t3 = spawn(async move { r3.run_process().await });

    // ... wait until the Maelstrom system closes stdin and stdout
    info!("running");
    let _ignored = tokio::join!(t1, t2, t3);

    info!("stopped");

    Ok(())
}
