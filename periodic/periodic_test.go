package periodic

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var globalVal1 int
var globalVal2 int

func task() {
	globalVal1++
}

func taskWithParams(a, b int) {
	globalVal2 = globalVal2 + a + b
}

func Test_Scheduler_RunAndStop(t *testing.T) {
	globalVal1 = 0
	globalVal2 = 0

	s := NewScheduler()
	s.RegisterTask(time.Millisecond*500, true, "task1", task)
	s.RegisterTask(time.Millisecond*500, false, "task2", taskWithParams, 1, 1)

	s.Run()
	time.Sleep(time.Millisecond * 750)
	s.Stop()

	// Run And Stop
	assert.Equal(t, 2, globalVal1, "globalVal1 should be 2")
	assert.Equal(t, 2, globalVal2, "globalVal2 should be 2")

	// Run task1 Again
	globalVal1 = 0
	globalVal2 = 0
	s.Run("task1")

	time.Sleep(time.Millisecond * 750)
	s.Stop()

	assert.NotEqual(t, 0, globalVal1, "globalVal1 should not be 0")
	assert.Equal(t, 0, globalVal2, "globalVal2 should be 0")
}

func Test_Scheduler_StopSpecificTask(t *testing.T) {
	globalVal1 = 0
	globalVal2 = 0

	s := NewScheduler()
	s.RegisterTask(time.Millisecond*500, true, "task1", task)
	s.RegisterTask(time.Millisecond*500, false, "task2", taskWithParams, 1, 1)

	s.Run()
	time.Sleep(time.Millisecond * 750)
	s.Stop("task2")

	globalVal2 = 0
	time.Sleep(time.Second * 1)
	s.Stop()

	assert.NotEqual(t, 0, globalVal1, "globalVal1 should not be 0")
	assert.Equal(t, 0, globalVal2, "globalVal2 should be 0")
}

func Test_Scheduler_CancelTask(t *testing.T) {
	globalVal1 = 0
	globalVal2 = 0

	s := NewScheduler()
	s.RegisterTask(time.Millisecond*500, true, "task1", task)
	s.RegisterTask(time.Millisecond*500, false, "task2", taskWithParams, 1, 1)

	s.Run()
	time.Sleep(time.Millisecond * 750)
	s.Cancel("task1")

	assert.NotEqual(t, 0, globalVal1, "globalVal1 should not be 0")

	taskListLen := len(s.taskList)
	assert.Equal(t, taskListLen, 1, "Task list size should be 1")

	s.Cancel()

	taskListLen = len(s.taskList)
	assert.Equal(t, taskListLen, 0, "Task list size shoud be 0")
}

func Test_Scheduler_AnonymousFunc(t *testing.T) {
	globalVal1 = 0
	globalVal2 = 0

	anonymousFunc := func() {
		globalVal1++
	}

	anonymousFuncWithParam := func(a, b int) {
		globalVal2 = globalVal2 + a + b
	}

	s := NewScheduler()
	s.RegisterTask(time.Millisecond*500, true, "anonymous1", anonymousFunc)
	s.RegisterTask(time.Millisecond*500, false, "anonymous2", anonymousFuncWithParam, 1, 1)

	s.Run()
	time.Sleep(time.Millisecond * 750)
	s.Stop()

	assert.Equal(t, 2, globalVal1, "globalVal1 should be 2")
	assert.Equal(t, 2, globalVal2, "globalVal2 should be 2")
}

func Test_Scheduler_Resume(t *testing.T) {

	nowTs := time.Now().UnixNano()

	s := NewScheduler()
	s.RegisterTask(time.Millisecond*500, true, "aa", func() {
		nowTs2 := time.Now().UnixNano()
		fmt.Println((nowTs2 - nowTs), " - tick aa")
	})

	s.RegisterTask(time.Millisecond*500, true, "bb", func() {
		nowTs2 := time.Now().UnixNano()
		fmt.Println((nowTs2 - nowTs), " - tick bb")
	})

	s.Run()

	time.Sleep(time.Millisecond * 250)

	s.Call("aa")

	time.Sleep(time.Millisecond * 3000)

	s.Cancel("bb")

	time.Sleep(time.Millisecond * 3000)
}

func TestSchedulerOption(t *testing.T) {
	nowTs := time.Now().UnixNano() / 1e6

	s := NewScheduler()
	s.RegisterTaskOption(Option{
		Immediately: true,
		Interval:    time.Millisecond * 500,
		Name:        "aa",
		Func: func() {
			nowTs2 := time.Now().UnixNano() / 1e6
			fmt.Println((nowTs2 - nowTs), " - tick aa")
			time.Sleep(1 * time.Second)
		},
	})

	s.Run()

	time.Sleep(time.Second * 10)

}
