// Copied from: https://www.youtube.com/watch?v=Utf-A4rODH8 and https://gist.github.com/julz/c0017fa7a40de0543001
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
    parent()
  case "child":
    child()
  default:
    panic("Wrong input parameters provided")
  }
}

func parent() {
  cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
  }
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  
  must(cmd.Run())
}

func child() {
  fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())
  
  must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
  must(os.MkdirAll("rootfs/oldrootfs", 0700))
	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	must(os.Chdir("/"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
  
  must(cmd.Run())
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}
