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

	if len(a) > 0 {
		fmt.Println(fmt.Sprintf(builder.String(), time.Now(), file, line, a))
		return
	}

	fmt.Println(fmt.Sprintf(builder.String(), time.Now(), file, line))
}

func Error(format string, a ...any) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("runtime.Caller failed")
	}

	var builder strings.Builder
	builder.WriteString("%v:%v ")
	builder.WriteString(format)

	if len(a) > 0 {
		return fmt.Errorf(fmt.Sprintf(builder.String(), file, line, a))
	}

	return fmt.Errorf(fmt.Sprintf(builder.String(), file, line))
}

func InitGlog() *log.Logger {
	return log.Default()
}
