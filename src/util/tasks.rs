use std::error::Error;
use std::sync::Arc;

use chrono::{DateTime, Utc};
use mongodb::Database;
use tokio::sync::mpsc::{channel, Receiver, Sender};
use tokio::sync::Mutex;
use tokio::time::{sleep, Duration};

use super::types::ClockState;

pub struct NewSender {
    /// True if connecting, false if disconnecting
    pub connect: bool,
    pub tx: Sender<String>,
}

/// Initialization for threads
pub async fn init_threads(db: Arc<Database>) -> (Sender<NewSender>, Sender<String>, Arc<Mutex<ClockState>>) {
    // Create a shared vector for holding senders for each client thread
    let tx_vector = Arc::new(Mutex::<Vec<Sender<String>>>::new(vec![]));
    let tx_vector_sender_listener = tx_vector.clone();
    let tx_vector_mongodb_listener = tx_vector.clone();
    let tx_vector_clock = tx_vector.clone();

    // Generate channel for inter-task communication for mongodb stats listener
    let (vector_tx, vector_rx) = channel::<NewSender>(10);

    // Create the clock state and sender
    let clock_state = Arc::new(Mutex::new(ClockState::new()));
    let clock_state_task = clock_state.clone();
    let (clock_tx, clock_rx) = channel::<String>(10);

    // Thread for listening for new recievers
    tokio::spawn(async move {
        new_sender_listen(tx_vector_sender_listener, vector_rx).await;
    });

    // Thread for listening for mongodb updates
    tokio::spawn(async move {
        match mongo_listen(tx_vector_mongodb_listener, db).await {
            Ok(()) => None,
            Err(err) => Some({
                eprintln!("ERROR in mongo_listen: {}", err);
            }),
        }
    });

    // Spawn the clock thread
    tokio::spawn(async move {
        clock_task(clock_state_task, clock_rx, tx_vector_clock).await;
    });

    return (vector_tx, clock_tx, clock_state);
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

/// Thread that creates the MongoDB listener and listens for changes to the database
async fn mongo_listen(
    sender_list: Arc<Mutex<Vec<Sender<String>>>>,
    db: Arc<Database>,
) -> Result<(), Box<dyn Error>> {
    // Instantiate variables for tracking debounce time
    let mut last_update = DateTime::<Utc>::MIN_UTC.timestamp_millis();
    let mut debounced = false;

    // Create and listen to change stream
    let mut change_stream = db.watch(None, None).await?;
    while change_stream.is_alive() {
        if let Some(_) = change_stream.next_if_any().await? {
            let now = Utc::now().timestamp_millis();
            if now - last_update > 10000 {
                // If last update was more than 10 seconds ago, send update
                // and reset debounce flag
                last_update = now;
                debounced = false;
                update_senders(sender_list.clone(), "stats".to_string()).await;
            } else if !debounced {
                // If not yet debounced an update, run the debounce task
                let debounce_sender_list = sender_list.clone();
                let diff = now - last_update;
                debounced = true;
                tokio::spawn(async move {
                    debounce_task(debounce_sender_list, diff).await;
                });
            };
        }
    }

    Ok(())
}

/// Thread that debounces the change_stream.
/// This will basically wait for 11 seconds before sending an update
/// so too many updates don't get sent at once
async fn debounce_task(sender_list: Arc<Mutex<Vec<Sender<String>>>>, diff: i64) {
    // Sleep for 11 seconds to catch all events that happen
    // Extra second is allocated for overlap with next debounce cycle
    sleep(Duration::from_millis(
        (11000 - diff).try_into().unwrap_or_else(|_| 0),
    ))
    .await;
    update_senders(sender_list, "stats".to_string()).await;
}

/// Sends the update message to all message channels
async fn update_senders(sender_list: Arc<Mutex<Vec<Sender<String>>>>, msg: String) {
    let mut invalid = Vec::new();
    let mut editable_sender_list = sender_list.lock().await;
    for (i, s) in editable_sender_list.iter().enumerate() {
        match s.send(msg.clone()).await {
            Ok(()) => (),
            Err(_) => {
                // Insert invalid senders in reverse order
                invalid.insert(0, i);
            }
        };
    }
    // Remove invalid senders from the sender vector
    if !invalid.is_empty() {
        for i in invalid.iter() {
            editable_sender_list.remove(*i);
        }
    }
}

async fn clock_task(clock_state: Arc<Mutex<ClockState>>, mut rx: Receiver<String>, sender_list: Arc<Mutex<Vec<Sender<String>>>>) {
    loop {
        // Wait for message
        let message = match rx.recv().await {
            Some(x) => x,
            None => continue,
        };

        // Get mutable clock state
        let mut mut_clock_state = clock_state.lock().await;

        // Change flag
        let mut changed = false;

        // If reset, reset the clock
        if message == "reset" {
            mut_clock_state.reset();
            changed = true;
        }

        // If paused, stop the clock
        if message == "pause" {
            mut_clock_state.pause();
            changed = true;
        }

        // If unpaused, start the clock
        if message == "unpause" {
            mut_clock_state.resume();
            changed = true;
        }

        // If changed, send update to all clients
        if changed {
            update_senders(sender_list.clone(), "clock".to_string()).await;
        }
    }
}
