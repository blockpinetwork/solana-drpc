package queue

type Job interface {
	Do()
}

type Worker struct {
	Job 	chan Job
	Stop	chan bool
}

func NewWorker() Worker{
	return Worker{
		Job: make(chan Job),
		Stop: make(chan bool),
	}
}

func (w Worker) Run(wq chan chan Job)  {
	go func() {
		for{
			wq <- w.Job
			select {
			case job := <- w.Job:
				job.Do()
			case <- w.Stop:
				return
			}
		}
	}()
}

type Queue struct {
	queueSize 	int
	JobQueue	chan Job
	WorkerQueue	chan chan Job
}

func NewQueue(size int) *Queue {
	return &Queue{
		queueSize: size,
		JobQueue: make(chan Job),
		WorkerQueue: make(chan chan Job, size),
	}
}

func (q *Queue) Run() {
	for i := 0; i < q.queueSize; i++ {
		NewWorker().Run(q.WorkerQueue)
	}
	go func() {
		for {
			select {
			case job := <- q.JobQueue:
				worker := <- q.WorkerQueue
				worker <- job
			}
		}
	}()
}

func (q *Queue) PushJob(job Job)  {
	q.JobQueue <- job
}