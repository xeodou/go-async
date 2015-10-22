/*
* @Author: xeodou
* @Date:   2015
 */

package task

import (
	"io/ioutil"
	"net/http"
)

// type TaskManager struct {
// 	Tasks []Task
// }

// var instance *TaskManager
// var once sync.Once

// func Instance() *TaskManager {
// 	once.Do(func() {
// 		instance = &TaskManager{}
// 	})
// 	return instance
// }

// func (tm *TaskManager) Add(t ...Task) {
// 	tm.Tasks = append(tm.Tasks, t...)
// }

type TaskInterface interface {
	Success(buf []byte)
	Failure(err error)
}

type Task struct {
	Url      string
	Listener TaskInterface
}

func NewTask(url string, listener TaskInterface) *Task {
	return &Task{url, listener}
}

type response struct {
	data chan []byte
	err  chan error
}

func (t *Task) create() *response {
	res := &response{make(chan []byte), make(chan error)}
	go func() {
		for {
			select {
			case err := <-res.err:
				t.Listener.Failure(err)
				break
			case data := <-res.data:
				t.Listener.Success(data)
			}
		}
	}()
	return res
}

func (t *Task) Runtask() {
	r := t.create()
	go func() {
		res, err := http.Get(t.Url)
		if err != nil {
			r.err <- err
		} else {
			buf, err := ioutil.ReadAll(res.Body)
			if err != nil {
				r.err <- err
			} else {
				r.data <- buf
			}
		}
	}()
}
