/*
* @Author: xeodou
* @Date:   2015
* @Last Modified by:   xeodou
* @Last Modified time: 2015
 */

package task

import (
	"fmt"
	"sync"
	"testing"
)

const (
	test_url = "https://www.google.com"
	count    = 1
)

var wg sync.WaitGroup
var data string

type Tmp struct{}

func (t *Tmp) Success(buf []byte) {
	data = string(buf)
	wg.Done()
}

func (t *Tmp) Failure(err error) {
	fmt.Printf("%s", err)
	wg.Done()
}

func TestTask(t *testing.T) {
	wg.Add(1)

	task := NewTask(test_url, &Tmp{})
	task.Runtask()

	wg.Wait()

	if data == "" {
		t.Errorf("Run async task failed.")
	}

}
