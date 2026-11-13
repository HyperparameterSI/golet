//go:build darwin

// Author: HyperparameterSI
package main

import (
	"fmt"
	"github.com/ebitengine/purego"
	"os"
	"syscall"
	"unsafe"
)

type darwinSource struct {
	lib        uintptr
	host, task uint32
	err        error
	stats      func(uint32, int32, *uint32, *uint32) int32
	deallocate func(uint32, uint32) int32
	procInfo   func(int32, int32, uint64, unsafe.Pointer, int32) int32
}

func newSource() source {
	s := &darwinSource{}
	lib, e := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_LOCAL)
	if e != nil {
		s.err = e
		return s
	}
	s.lib = lib
	var hostSelf func() uint32
	purego.RegisterLibFunc(&hostSelf, lib, "mach_host_self")
	purego.RegisterLibFunc(&s.stats, lib, "host_statistics")
	purego.RegisterLibFunc(&s.deallocate, lib, "mach_port_deallocate")
	purego.RegisterLibFunc(&s.procInfo, lib, "proc_pidinfo")
	var taskSelf func() uint32
	purego.RegisterLibFunc(&taskSelf, lib, "task_self_trap")
	s.task = taskSelf()
	s.host = hostSelf()
	return s
}
func (s *darwinSource) Close() {
	if s.host != 0 {
		s.deallocate(s.task, s.host)
	}
	if s.lib != 0 {
		_ = purego.Dlclose(s.lib)
	}
}
func (s *darwinSource) Read() raw {
	r := raw{errors: []string{}}
	if s.err != nil {
		r.errors = append(r.errors, s.err.Error())
		return r
	}
	var usage syscall.Rusage
	if e := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); e == nil {
		r.cpuSeconds = float64(usage.Utime.Sec+usage.Stime.Sec) + float64(usage.Utime.Usec+usage.Stime.Usec)/1e6
		r.processOK = true
	} else {
		r.errors = append(r.errors, "进程 CPU: "+e.Error())
	}
	// proc_taskinfo: six uint64 fields followed by twelve int32 fields (Darwin SDK).
	var task struct {
		Values   [6]uint64
		Counters [12]int32
	}
	if s.procInfo(int32(os.Getpid()), 4, 0, unsafe.Pointer(&task), 96) == 96 {
		v := task.Values[1]
		r.rss = &v
	} else {
		r.errors = append(r.errors, "无法读取进程 RSS")
	}
	count := uint32(4)
	if s.stats(s.host, 3, &r.ticks[0], &count) == 0 && count == 4 {
		r.systemOK = true
	} else {
		r.errors = append(r.errors, "无法读取系统 CPU")
	}
	var sysctl func(string, *uint64, *uintptr, uintptr, uintptr) int32
	purego.RegisterLibFunc(&sysctl, s.lib, "sysctlbyname")
	var total uint64
	size := uintptr(8)
	code := sysctl("hw.memsize", &total, &size, 0, 0)
	var vm [15]uint32
	count = 15
	if code == 0 && size == 8 && s.stats(s.host, 2, &vm[0], &count) == 0 && count >= 4 {
		free := uint64(vm[0]) * uint64(os.Getpagesize())
		used := total - min(total, free)
		r.total = &total
		r.used = &used
	} else {
		r.errors = append(r.errors, fmt.Sprint("无法读取系统内存，状态 ", code))
	}
	return r
}
