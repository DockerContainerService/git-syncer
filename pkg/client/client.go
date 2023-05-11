package client

import (
	"fmt"
	"github.com/DockerContainerService/git-syncer/pkg/task"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type taskList struct {
	tasks []task.Task
	lock  sync.Mutex
}

func (t *taskList) Push(tk task.Task) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.tasks = append(t.tasks, tk)
}

func (t *taskList) All() []task.Task {
	return t.tasks
}

func (t *taskList) Len() int {
	return len(t.tasks)
}

type Client struct {
	taskList *taskList
	config   *Config

	retries        int
	routineNum     int
	privateKeyFile string

	failedTaskList *taskList
}

func Default(configFile string, retries, routineNum int, privateKeyFile string) (client *Client, err error) {
	config, err := ParseConfig(configFile)
	if err != nil {
		return
	}
	client = &Client{
		config: config,

		retries:        retries,
		routineNum:     routineNum,
		privateKeyFile: privateKeyFile,

		taskList: &taskList{
			tasks: make([]task.Task, 0),
			lock:  sync.Mutex{},
		},
		failedTaskList: &taskList{
			tasks: make([]task.Task, 0),
			lock:  sync.Mutex{},
		},
	}
	return
}

func (c *Client) Run() {
	c.generateTask()
	ch := make(chan struct{}, c.routineNum)
	var wg sync.WaitGroup
	for _, t := range c.taskList.All() {
		ch <- struct{}{}
		wg.Add(1)
		logrus.Infof("Run task %s: %s => %s", t.GetTitle(), t.GetSrcRepo(), t.GetDstRepo())
		go func(t task.Task) {
			defer wg.Done()
			num := 0
			for num < c.retries+1 {
				if num != 0 {
					logrus.Warnf("Retry task %s, times %d", t.GetTitle(), num)
				}
				err := t.Run()
				if err != nil {
					logrus.Debugf("Run task %s %s => %s err: %v", t.GetTitle(), t.GetSrcRepo(), t.GetDstRepo(), err)
					if num == c.retries-1 {
						t.SetError(err)
						c.failedTaskList.Push(t)
						logrus.Errorf("Run task %s %s => %s err: %v", t.GetTitle(), t.GetSrcRepo(), t.GetDstRepo(), err)
					}
					num++
				} else {
					break
				}
			}
			<-ch
		}(t)
	}
	wg.Wait()
	msg := fmt.Sprintf("Finished, %d task succeeded, %d task failed", c.taskList.Len()-c.failedTaskList.Len(), c.failedTaskList.Len())
	if c.failedTaskList.Len() > 0 {
		msg += ". failed task list: \n\n"
		for _, t := range c.failedTaskList.All() {
			msg += fmt.Sprintf("%s: %s => %s, err: %+v\n", t.GetTitle(), t.GetSrcRepo(), t.GetDstRepo(), t.GetError())
		}
	}
	msg = strings.TrimSuffix(msg, "\n")
	logrus.Infof(msg)
}

func (c *Client) generateTask() {
	for k, v := range *c.config {
		t := task.NewMigrationTask(k, v, c.privateKeyFile)
		c.taskList.Push(t)
		logrus.Infof("Generate task %s: %s => %s", t.GetTitle(), k, v)
	}
}
