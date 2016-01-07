package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	fileName := "temp.txt"

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// 锁定
	// LOCK_EX exclusive  Get an exclusive lock
	// LOCK_NB nonblock   Fail rather than wait
	// LOCK_UN unlock     Remove a lock
	// LOCK_SH shared     Get a shared lock
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("flock", fileName, time.Now())

	time.Sleep(10 * time.Second)

	// 解锁
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("unlock", fileName, time.Now())
}
