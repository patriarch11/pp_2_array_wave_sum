use crate::pool::{ThreadJob, ThreadPool};

mod pool;

const VEC_LENGTH: usize = 7;

fn add(start: usize, vec: Vec<i32>) -> i32 {
    let length = vec.len();
    let is_length_even = length % 2 == 0;

    if !is_length_even {
        /*
        	* якщо довжина вектора - це непарне число,
        	* а стартовий індекс знаходиться посередині,
        	* то повертаємо число, яке знаходиться на стартовому індексі
         */
        if start as f64 >= (length as f64 / 2.0).floor() {
            return vec[start];
        }
    }
    let end = length - start - 1;
    vec[start] + vec[end]
}

fn wave_sum(v: Vec<i32>, pool: &ThreadPool, mut wave_num: i32) -> i32 {
    let stop_idx = (v.len() as f64 / 2.0).ceil() as usize;
    let mut jobs = Vec::with_capacity(stop_idx);

    for i in 0..stop_idx {
        jobs.push(ThreadJob {
            start: i,
            vec: v.clone(),
            job_fn: add,
        })
    }

    let res = pool.process(jobs);

    println!("Хвиля №{}, результат {:?}", wave_num, res);

    if res.len() > 1 {
        // рекурсивний кейс
        wave_num += 1;
        return wave_sum(res, pool, wave_num);
    }
    res[0] // термінальний кейс
}

fn main() {
    let num_cores = std::thread::available_parallelism().unwrap().get();
    let vec: Vec<i32> = (1..VEC_LENGTH as i32).collect();
    let sync_sum: i32 = vec.iter().sum();
    println!("Максимальна к-ть потоків, для оптимальної роботи: {}", num_cores);
    println!("Довжина вектора: {}", VEC_LENGTH);
    println!("Вектор: {:?}", vec);
    println!("Синхронно порахований результат: {}", sync_sum);

    let mut pool = ThreadPool::new(num_cores);
    let res = wave_sum(vec, &pool, 0);
    pool.join();
    println!("Результат хвилевого алгоритму: {:?}", res);
}
