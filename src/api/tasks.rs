use std::error::Error;
use std::sync::Arc;
// use std::time::Duration;
// use tokio::time;

use mongodb::Database;
use tokio::sync::mpsc::{channel, Receiver, Sender};
use tokio::sync::Mutex;

pub struct NewSender {
    /// True if connecting, false if disconnecting
    pub connect: bool,
    pub tx: Sender<String>,
}

/// Initialization for threads
pub async fn init_threads(db: Arc<Database>) -> Sender<NewSender> {
    // Create a shared vector for holding senders for each client thread
    let tx_vector = Arc::new(Mutex::<Vec<Sender<String>>>::new(vec![]));
    let tx_vector_sender_listener = tx_vector.clone();
    let tx_vector_mongodb_listener = tx_vector.clone();

    // Generate channel for inter-task communication for mongodb stats listener
    let (vector_tx, vector_rx) = channel::<NewSender>(10);

    // Thread for listening for new recievers
    tokio::spawn(async move {
        new_sender_listen(tx_vector_sender_listener, vector_rx).await;
    });

    // Thread for listening for mongodb updates
    tokio::spawn(async move {
        match mongo_listen(tx_vector_mongodb_listener, db).await {
            Ok(()) => None,
            Err(err) => Some({
                println!("ERROR in mongo_listen: {}", err);
            }),
        }
    });

    return vector_tx;
}

/// Thread to listen for new senders (new thread connections)
async fn new_sender_listen(
    sender_list: Arc<Mutex<Vec<Sender<String>>>>,
    mut rx: Receiver<NewSender>,
) {
    loop {
        // Wait for message (happens when a new admin client connects)
        let message = match rx.recv().await {
            Some(x) => x,
            None => continue,
        };

        // Retrieve the vector from the mutex
        let mut editable_sender_list = sender_list.lock().await;

        // Based on the message, either append or delete the Sender from the vector
        if message.connect {
            editable_sender_list.push(message.tx);
        } else {
            if let Some(pos) = editable_sender_list
                .iter()
                .position(|x| x.same_channel(&message.tx))
            {
                editable_sender_list.remove(pos);
            }
        }
    }
}

/// Creates the MongoDB listener and listens for changes to the database
async fn mongo_listen(
    sender_list: Arc<Mutex<Vec<Sender<String>>>>,
    db: Arc<Database>,
) -> Result<(), Box<dyn Error>> {
    let mut change_stream = db.watch(None, None).await?;
    while change_stream.is_alive() {
        // TODO: simply continue if change stream is invalid and show error
        if let Some(_) = change_stream.next_if_any().await? {
            let editable_sender_list = sender_list.lock().await;
            for s in editable_sender_list.iter() {
                s.send("update".to_string()).await?;
            }
        }
    }

    Ok(())
}
