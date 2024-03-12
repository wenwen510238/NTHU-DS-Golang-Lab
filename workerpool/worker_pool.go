package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type Task struct {
	Func func(args ...interface{}) *Result
	Args []interface{}
}

type Result struct {
	Value interface{}
	Err   error
}

type WorkerPool interface {
	Start(ctx context.Context)
	Tasks() chan *Task
	Results() chan *Result
}

type workerPool struct {
	numWorkers int
	tasks      chan *Task
	results    chan *Result
	wg         *sync.WaitGroup
}

var _ WorkerPool = (*workerPool)(nil)

func NewWorkerPool(numWorkers int, bufferSize int) *workerPool {
	return &workerPool{
		numWorkers: numWorkers,
		tasks:      make(chan *Task, bufferSize),
		results:    make(chan *Result, bufferSize),
		wg:         &sync.WaitGroup{},
	}
}

func (wp *workerPool) Start(ctx context.Context) {
	// TODO: implementation
	//

	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.run(ctx)
	}
	wp.wg.Wait()
	close(wp.results)

	// Starts numWorkers of goroutines, wait until all jobs are done.
	// Remember to closed the result channel before exit.
}

func (wp *workerPool) Tasks() chan *Task {
	return wp.tasks
}

func (wp *workerPool) Results() chan *Result {
	return wp.results
}

func (wp *workerPool) run(ctx context.Context) {
	// TODO: implementation
	//
	defer wp.wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("工作結束")
			return
		case task, ok := <-wp.tasks:
			if !ok {
				fmt.Println("通道已關閉")
				return
			}
			select {
			case <-ctx.Done():
				fmt.Println("任務被取消")
				return
			default:
			}
			result := task.Func(task.Args...)
			wp.results <- result
			// select {
			// case wp.results <- result:
			// 	// wp.results <- result
			// case <-ctx.Done():
			// 	fmt.Println("收到取消信號")
			// 	return
			// }
		}
	}
	// Keeps fetching task from the task channel, do the task,
	// then makes sure to exit if context is done.
}
