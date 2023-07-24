package server

import (
	"GtBase/pkg/constants"
	"log"
	"syscall"
)

type Ioer interface {
	Run(int)
	Wait() ([]*Task, error)
	AddRead(int) error
}

type Task struct {
	eventFd   int
	eventType int32
}

func (t *Task) EventFd() int {
	return t.eventFd
}

func (t *Task) EventType() int32 {
	return t.eventType
}

func CreateTask(eventFd int, eventType int32) *Task {
	return &Task{eventFd: eventFd, eventType: eventType}
}

type EPoller struct {
	epollFd  int
	listenFd int
}

func (p *EPoller) Run(listenFd int) {
	epollFd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatalf("can't use EpollCreate1 becasue %v", err)
		return
	}

	p.epollFd = epollFd
	p.listenFd = listenFd

	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(listenFd),
	}

	erre := syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, listenFd, &event)
	if erre != nil {
		log.Fatalf("failed to add ListenFd because %v", err)
		return
	}
}

func (p *EPoller) Wait() ([]*Task, error) {
	result := make([]*Task, 0, 10)
	events := make([]syscall.EpollEvent, 10)
	n, err := syscall.EpollWait(p.epollFd, events, -1)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		event := events[i]
		if event.Events&syscall.EPOLLIN != 0 {
			if int(event.Fd) == p.listenFd {
				result = append(result, CreateTask(int(event.Fd), constants.IoerAccept))
			} else {
				result = append(result, CreateTask(int(event.Fd), constants.IoerRead))
			}
		}
	}

	return result, nil
}

func (p *EPoller) AddRead(addFd int) error {
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(addFd),
	}

	err := syscall.EpollCtl(p.epollFd, syscall.EPOLL_CTL_ADD, addFd, &event)
	if err != nil {
		return err
	}

	return nil
}
