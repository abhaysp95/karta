package main

import (
	"fmt"
	"sync"
	"time"
)

type Job func()

type Pool struct {
	queue chan Job
	queueSize int
	wg sync.WaitGroup
}

func NewPool(workCount int) *Pool {
	pool := &Pool {
		queue: make(chan Job, workCount),
		queueSize: workCount,
	}

	// fmt.Println("Job with workCount: ", workCount)
	return pool
}

func (p *Pool) AddJob(job Job) {
	p.wg.Add(1)
	p.queue <- job
	// fmt.Println("adding job")
}

func (p *Pool) ProcessQueue() {
	/* defer p.wg.Done()
	for {
		select {
		case job := <- p.queue:
			go func() {
				defer p.wg.Done()
				job()
			}()
		default:
			fmt.Println("Waiting for job")
			time.Sleep(time.Second * 1)
		}
	} */
	for i := 0; i < p.queueSize; i++ {
		go func() {
			for job := range p.queue {
				job()
				p.wg.Done()
			}
		} ()
	}
}

func DownloadJob(count int) Job {
	return func() {
		fmt.Println("Downloading Video: ", count)
		time.Sleep(time.Second * 10)
		fmt.Println("Downloaded Video: ", count)
	}
}

func (p *Pool) WaitJob() {
	close(p.queue)
	p.wg.Wait()
}

func main() {
	pool := NewPool(5)
	var count int = 30

	// pool.wg.Add(1)
	go pool.ProcessQueue()
	for i := 0; i < count; i++ {
		fmt.Println("Calling for count: ", i)
		pool.AddJob(DownloadJob(i))
	}

	pool.WaitJob()
}
