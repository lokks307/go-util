package periodic

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type taskStatus int

const (
	running taskStatus = 1 + iota
	stop
)

//TaskInfo struct keep information about job
type TaskInfo struct {
	taskFunction interface{}
	taskParams   []interface{}
	interval     time.Duration
	ticker       *time.Ticker

	immediately bool
	status      taskStatus
}

//Scheduler struct keep TaskInfos
type Scheduler struct {
	taskList map[string]*TaskInfo
	rwMutex  *sync.RWMutex
}

//NewScheduler make new scheduler struct
func NewScheduler() *Scheduler {
	scheduler := &Scheduler{
		taskList: make(map[string]*TaskInfo),
		rwMutex:  new(sync.RWMutex),
	}
	return scheduler
}

//RegisterTask regiseter task
func (s *Scheduler) RegisterTask(interval time.Duration, immediately bool, taskNameKey string, taskFunc interface{}, params ...interface{}) error {
	typ := reflect.TypeOf(taskFunc)
	if typ.Kind() != reflect.Func {
		return errors.New("only function can be registered")
	}

	f := reflect.ValueOf(taskFunc)
	if len(params) != f.Type().NumIn() {
		return errors.New("the number of params is not matched")
	}

	if _, ok := s.taskList[taskNameKey]; ok {
		return errors.New("this task function is already registred")
	}

	s.taskList[taskNameKey] = &TaskInfo{
		taskFunction: taskFunc,
		taskParams:   params,
		interval:     interval,
		immediately:  immediately,
		status:       stop,
	}
	return nil
}

func (t *TaskInfo) run() {
	f := reflect.ValueOf(t.taskFunction)
	in := make([]reflect.Value, len(t.taskParams))

	for k, param := range t.taskParams {
		in[k] = reflect.ValueOf(param)
	}

	t.ticker = time.NewTicker(t.interval)
	go func() {
		if t.immediately {
			for ; true; <-t.ticker.C {
				f.Call(in)
			}
		} else {
			for range t.ticker.C {
				f.Call(in)
			}
		}
	}()
}

//Run registered tasks ( if params do not exist, run all tasks. on the other hand, run specific tasks)
func (s *Scheduler) Run(taskNames ...string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			if task.status == running {
				continue
			}
			task.run()
			task.status = running
		}
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			if task.status == running {
				continue
			}
			task.run()
			task.status = running
		}
	}
}

//Stop registered tasks ( if params do not exist, stop all tasks. on the other hand, stop specific tasks)
func (s *Scheduler) Stop(taskNames ...string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			if task.status == stop {
				continue
			}
			task.ticker.Stop()
			task.status = stop
		}
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			if task.status == stop {
				continue
			}
			task.ticker.Stop()
			task.status = stop
		}
	}
}

//Cancel registered tasks ( if params do not exist, remove all tasks. on the other hand, remove specific tasks)
func (s *Scheduler) Cancel(taskNames ...string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			if task.status == running {
				task.ticker.Stop()
			}
		}
		s.taskList = make(map[string]*TaskInfo)
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			if task.status == running {
				task.ticker.Stop()
			}
			delete(s.taskList, taskName)
		}
	}
}
