package glog

import (
	"log"
	"sync"
)

var once sync.Once
var glog *log.Logger

func Glog() *log.Logger {
	once.Do(func() {
		glog = InitGlog()
	})

	return glog
}

func InitGlog() *log.Logger {
	return log.New(log.Writer(), "", log.LstdFlags|log.Llongfile|log.LUTC)
}
