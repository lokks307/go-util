package periodic

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/subchen/go-trylock/v2"
)

const (
	running int32 = 1 + iota
	stop
	deleted
)

//TaskInfo struct keep information about job
type TaskInfo struct {
	taskFunction interface{}
	taskParams   []interface{}
	interval     time.Duration
	ticker       *time.Ticker
	immediately  bool
	concurrency  bool
	locker       trylock.TryLocker
	status       int32
	name         string
}

type Option struct {
	Interval      time.Duration
	Immediately   bool
	ConcurrentRun bool
	Name          string
	Func          interface{}
	Parameters    []interface{}
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

func (s *Scheduler) RegisterTaskOption(option Option) error {
	err := s.RegisterTask(option.Interval, option.Immediately, option.Name, option.Func, option.Parameters...)
	if err != nil {
		return err
	}

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if task, ok := s.taskList[option.Name]; ok {
		task.concurrency = option.ConcurrentRun
	}

	return nil
}

//RegisterTask regiseter task
func (s *Scheduler) RegisterTask(interval time.Duration, immediately bool, taskNameKey string, taskFunc interface{}, params ...interface{}) error {

	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

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
		concurrency:  true,
		locker:       trylock.New(),
		name:         taskNameKey,
	}
	return nil
}

func (t *TaskInfo) call() {
	f := reflect.ValueOf(t.taskFunction)
	in := make([]reflect.Value, len(t.taskParams))

	for k, param := range t.taskParams {
		in[k] = reflect.ValueOf(param)
	}

	go func() {
		if !t.concurrency {
			if lok := t.locker.TryLock(nil); lok {
				defer t.locker.Unlock()
				f.Call(in)
			}
		} else {
			f.Call(in)
		}
	}()
}

func (t *TaskInfo) resume() {
	t.ticker.Reset(t.interval)
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
				if atomic.LoadInt32(&t.status) == deleted {
					t.ticker.Stop()
					break
				}
				go func() {
					if !t.concurrency {
						if lok := t.locker.TryLock(nil); lok {
							defer t.locker.Unlock()
							f.Call(in)
						}
					} else {
						f.Call(in)
					}
				}()
			}
		} else {
			for range t.ticker.C {
				if atomic.LoadInt32(&t.status) == deleted {
					t.ticker.Stop()
					break
				}

				go func() {
					if !t.concurrency {
						if lok := t.locker.TryLock(nil); lok {
							defer t.locker.Unlock()
							f.Call(in)
						}
					} else {
						f.Call(in)
					}
				}()
			}
		}
	}()
}

//Run registered tasks ( if params do not exist, run all tasks. on the other hand, run specific tasks)
func (s *Scheduler) Run(taskNames ...string) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			if atomic.LoadInt32(&task.status) == running {
				continue
			}
			task.run()
			atomic.StoreInt32(&task.status, running)
		}
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			if atomic.LoadInt32(&task.status) == running {
				continue
			}
			task.run()
			atomic.StoreInt32(&task.status, running)
		}
	}
}

//Stop registered tasks ( if params do not exist, stop all tasks. on the other hand, stop specific tasks)
func (s *Scheduler) Stop(taskNames ...string) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			if atomic.LoadInt32(&task.status) == stop {
				continue
			}
			task.ticker.Stop()
			atomic.StoreInt32(&task.status, stop)
		}
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			if atomic.LoadInt32(&task.status) == stop {
				continue
			}
			task.ticker.Stop()
			atomic.StoreInt32(&task.status, stop)
		}
	}
}

//Cancel registered tasks ( if params do not exist, remove all tasks. on the other hand, remove specific tasks)
func (s *Scheduler) Cancel(taskNames ...string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			atomic.StoreInt32(&task.status, deleted) // will automatically break loop
		}
		s.taskList = make(map[string]*TaskInfo)
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			atomic.StoreInt32(&task.status, deleted)
			delete(s.taskList, taskName)
		}
	}
}

func (s *Scheduler) Call(taskNames ...string) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if len(taskNames) == 0 {
		for _, task := range s.taskList {
			if atomic.LoadInt32(&task.status) == running {
				task.call()
				task.resume()
			}
		}
		return
	}

	for _, taskName := range taskNames {
		if task, ok := s.taskList[taskName]; ok {
			if atomic.LoadInt32(&task.status) == running {
				task.call()
				task.resume()
			}
		}
	}

}
