use std::sync::{ Arc, Mutex };
use std::thread;

fn main() {
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![]; // Vector to hold the thread handles

    for _ in 0..5 {
        // Spawn 5 threads
        let counter_clone = Arc::clone(&counter);
        let handle = thread::spawn(move || {
            for _ in 0..10000 {
                let mut num = counter_clone.lock().unwrap();
                *num += 1; // Increment the counter
            }
        });

        handles.push(handle); // Store the thread handle
    }

    for handle in handles {
        handle.join().unwrap(); // Wait for all threads to finish
    }

    let numbers: Vec<i32> = (1..1_000_000).collect();
    let sum: i32 = numbers
        .iter()
        .map(|x| x * 2)
        .filter(|x| x % 3 == 0)
        .sum();
}
