package glog

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var glog *log.Logger

func Glog() *log.Logger {
	once.Do(func() {
		glog = InitGlog()
	})

	return glog
}

func Log(format string, a ...any) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		Glog().Println("runtime.Caller failed")
	}

	var builder strings.Builder
	builder.WriteString("%v %v:%v ")
	builder.WriteString(format)

	Glog().Println(fmt.Sprintf(builder.String(), time.Now(), file, line, a))
}

func InitGlog() *log.Logger {
	return log.New(log.Writer(), "", log.LstdFlags|log.Llongfile|log.LUTC)
}
