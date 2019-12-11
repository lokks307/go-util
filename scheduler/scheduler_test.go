package scheduler

import (
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
