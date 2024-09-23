use std::sync::mpsc::{channel, Receiver, Sender};
use std::sync::{Arc, Mutex};
use std::thread;
use std::thread::JoinHandle;

pub struct ThreadJob {
    pub start: usize,
    pub vec: Vec<i32>,
    /*
    	* job_fn - функція, яка буде виконана в окремому потоці
    	* призначена для ізоляції "бізнес-логіки"
    	*/
    pub job_fn: fn(usize, Vec<i32>) -> i32,
}

pub struct ThreadPool {
    pool: Vec<JoinHandle<()>>,
    job_sender: Option<Sender<ThreadJob>>,
    res_recv: Receiver<i32>,
}

impl ThreadPool {
    fn spawn_thread(
        // id: usize,
        job_recv: Arc<Mutex<Receiver<ThreadJob>>>,
        res_sender: Sender<i32>,
    ) -> JoinHandle<()> {
        thread::spawn(move || {
            // println!("потік №{} запустився", id);
            loop {
                let job = job_recv.lock().unwrap().recv();
                match job {
                    Ok(job) => {
                        // println!("потік №{} отримав роботу", id);
                        let result = (job.job_fn)(job.start, job.vec);
                        // println!("потік №{} завершив роботу з результатом: {:?}", id, result);
                        res_sender.send(result).unwrap();
                    }
                    Err(_) => {
                        // println!("потік №{} завершився", id);
                        break;
                    }
                }
            }
        })
    }

    pub fn new(size: usize) -> ThreadPool {
        let (job_sender, job_recv) = channel();
        let (res_sender, res_recv) = channel();
        let job_recv: Arc<Mutex<Receiver<ThreadJob>>> = Arc::new(Mutex::new(job_recv));
        let mut pool = Vec::with_capacity(size);

        for _ in 0..size {
            // запускаємо потоки та додаємо їх в пул потоків
            pool.push(ThreadPool::spawn_thread(
                //id,
                job_recv.clone(),
                res_sender.clone(),
            ));
        }

        ThreadPool {
            pool,
            job_sender: Some(job_sender),
            res_recv,
        }
    }

    pub fn process(&self, jobs: Vec<ThreadJob>) -> Vec<i32> {
        let jobs_len = jobs.len();
        let mut results = Vec::with_capacity(jobs.len());

        if let Some(sender) = &self.job_sender {
            for job in jobs {
                sender.send(job).unwrap();
            }
        }

        for _ in 0..jobs_len {
            let result = self.res_recv.recv().unwrap();
            results.push(result);
        }
        results
    }

    pub fn join(&mut self) {
        if let Some(job_sender) = self.job_sender.take() {
            drop(job_sender);
        }
        for handle in self.pool.drain(..) {
            handle.join().unwrap();
        }
    }
}
