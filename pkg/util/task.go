package util

import (
	"sync"
)

type TaskPool struct {
	tasks chan func()
	group sync.WaitGroup
}

func NewTask(num int) *TaskPool {
	tasks := make(chan func(), num)
	for i := 0; i < int(num); i++ {
		go func() {
			for task := range tasks {
				task()
			}
		}()
	}
	return &TaskPool{tasks: tasks}
}

func (c *TaskPool) Run(task func()) {
	c.group.Add(1)
	c.tasks <- func() {
		defer c.group.Done()
		task()
	}

}
func (c *TaskPool) Wait() {
	c.group.Wait()
}
func (c *TaskPool) Close() {
	c.Wait()
	close(c.tasks)
}
