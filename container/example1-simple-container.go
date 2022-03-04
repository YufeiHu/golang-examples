// Copied from: https://www.youtube.com/watch?v=Utf-A4rODH8
// 
// Simple container that shows:
//   - UTS (Unix Timesharing System) namespace (i.e. hostname)
//   - process IDs namespace
//   - filesystem namespace (mount points)
// However, it does not show:
//   - user namespace
//   - IPC
//   - networking namespace
//   - cgroup
// 
// You must have "/home/rootfs" on your system which contains a linux filesystem to make the code run
// 
// Example usage: ./example1-simple-container run /bin/bash

package main

import (
  "fmt"
  "os"
  "os/exec"
  "syscall"
)

func main() {
  switch os.Args[1] {
  case "run":
    run()
  case "child":
    child()
  default:
    panic("Wrong input parameters provided")
  }
}

func run() {
  cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr {
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
  }
  
  must(cmd.Run())
}

func child() {
  fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())
  
  cmd := exec.Command(os.Args[2], os.Args[3:]...)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  
  must(syscall.Chroot("/home/rootfs"))
  must(os.Chdir("/"))
  must(syscall.Mount("proc", "proc", "proc", 0, ""))
  must(cmd.Run())
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}
