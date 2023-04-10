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

struct Server {
    args: Vec<String>,
    net: process::ProcNet<Generate, ()>,
    id: Id,
    ids: Vec<Id>,
    msg_id: AtomicU64,
}

impl Default for Server {
    fn default() -> Self {
        Server {
            args: Default::default(),
            net: Default::default(),
            id: Default::default(),
            ids: Default::default(),
            msg_id: Default::default(),
        }
    }
}

#[derive(Deserialize, Serialize, Debug, Eq, PartialEq)]
#[serde(tag = "type")]
pub enum Generate {
    #[serde(rename = "generate")]
    Generate { msg_id: msg::MsgId },
    #[serde(rename = "generate_ok")]
    GenerateOk {
        in_reply_to: msg::MsgId,
        msg_id: Option<msg::MsgId>,
        id: Value,
    },
}
#[async_trait]
impl process::Process<Generate, ()> for Server {
    fn init(
        &mut self,
        args: Vec<String>,
        net: process::ProcNet<Generate, ()>,
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
                    match full_message.body {
                        msg::Body::Workload(Generate::Generate { msg_id }) => {
                            let uid = Uuid::new_v4().to_string();
                            self.net
                                .txq
                                .send(msg::Msg {
                                    src: self.id.clone(),
                                    dest: msg_src,
                                    body: msg::Body::Workload(Generate::GenerateOk {
                                        in_reply_to: msg_id,
                                        msg_id: Some(
                                            self.msg_id
                                                .fetch_add(1, std::sync::atomic::Ordering::SeqCst),
                                        ),
                                        id: json!(uid),
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
    let server_instance: Server = Default::default();
    let r = Arc::new(runtime::Runtime::new(env::args().collect(), server_instance).await?);
    let (r1, r2, r3) = (r.clone(), r.clone(), r.clone());
    let t1 = spawn(async move { r1.run_io_egress().await });
    let t2 = spawn(async move { r2.run_io_ingress().await });
    let t3 = spawn(async move { r3.run_process().await });

    // ... wait until the Maelstrom system closes stdin and stdout
    info!("running");
    let _ignored = tokio::join!(t1, t2, t3);

    info!("stopped");

    Ok(())
}
