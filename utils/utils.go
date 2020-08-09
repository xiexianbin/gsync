package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func DateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Pid() int {
	return os.Getpid()
}

func GoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func Println(a ...interface{}) {
	str := fmt.Sprintf("[%s %d/%d]", DateTime(), Pid(), GoroutineID())
	a = append(a, str)
	_, _ = fmt.Fprintln(os.Stdout, a...)
}

func Md5sum(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	return fmt.Sprintf("%X", hash)
}
