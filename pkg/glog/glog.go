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

func Log(format string, a ...any) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("runtime.Caller failed")
	}

	var builder strings.Builder
	builder.WriteString("%v %v:%v ")
	builder.WriteString(format)

	fmt.Println(fmt.Sprintf(builder.String(), time.Now(), file, line, a))
}

func Error(format string, a ...any) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("runtime.Caller failed")
	}

	var builder strings.Builder
	builder.WriteString("%v:%v ")
	builder.WriteString(format)

	return fmt.Errorf(fmt.Sprintf(builder.String(), file, line, a))
}

func InitGlog() *log.Logger {
	return log.Default()
}
