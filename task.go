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

type success func(buf []byte)
type failure func(err string)

type Task struct {
	Url     string
	Success success
	Failure failure
}

func NewTaks(url string, suc success, fail failure) *Task {
	return &Task{url, suc, fail}
}

type response struct {
	data chan []byte
	err  chan string
}

func (t *Task) create() *response {
	res := &response{make(chan []byte), make(chan string)}
	go func() {
		for {
			select {
			case err := <-res.err:
				t.Failure(err)
				break
			case data := <-res.data:
				t.Success(data)
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
			r.err <- err.Error()
		} else {
			buf, err := ioutil.ReadAll(res.Body)
			if err != nil {
				r.err <- err.Error()
			} else {
				r.data <- buf
			}
		}
	}()
}
