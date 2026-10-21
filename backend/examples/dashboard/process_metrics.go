// Author: HyperparameterSI

package main

// ProcessSample is one Go-process observation. Byte counters exclude WebView
// subprocesses and are not RSS. Timestamp is Unix milliseconds.
type ProcessSample struct {
	// Sequence is a monotonic cursor, independent of wall-clock adjustments.
	Sequence  uint64 `json:"sequence"`
	Timestamp int64  `json:"timestamp"`
	// CPU is percent of one logical core; multi-core usage can exceed 100.
	CPU          float64 `json:"cpu"`
	CPUAvailable bool    `json:"cpuAvailable"`
	Heap         uint64  `json:"heap"`
	HeapInuse    uint64  `json:"heapInuse"`
	HeapIdle     uint64  `json:"heapIdle"`
	HeapReleased uint64  `json:"heapReleased"`
	// RuntimeMemory is Sys minus HeapReleased, in bytes.
	RuntimeMemory uint64 `json:"runtimeMemory"`
	Stack         uint64 `json:"stack"`
	Objects       uint64 `json:"objects"`
	// AllocRate is allocated bytes per second between consecutive samples.
	AllocRate  float64 `json:"allocRate"`
	TotalAlloc uint64  `json:"totalAlloc"`
	Goroutines int     `json:"goroutines"`
	// ThreadRecords counts runtime thread-profile records, not live OS threads.
	ThreadRecords int     `json:"threadRecords"`
	GCCycles      uint32  `json:"gcCycles"`
	GCDelta       uint32  `json:"gcDelta"`
	GCPauseMS     float64 `json:"gcPauseMS"`
	// GCPauseP95MS covers the latest at most 256 GC pauses, in milliseconds.
	GCPauseP95MS   float64 `json:"gcPauseP95MS"`
	GCTotalPauseMS float64 `json:"gcTotalPauseMS"`
	NextGC         uint64  `json:"nextGC"`
	// GCCPUFraction is a lifetime ratio, not the CPU percentage of this interval.
	GCCPUFraction    float64 `json:"gcCPUFraction"`
	SampleDurationMS float64 `json:"sampleDurationMS"`
}

// ProcessSnapshot describes a process and a chronological set of observations.
type ProcessSnapshot struct {
	Samples    []ProcessSample `json:"samples"`
	PID        int             `json:"pid"`
	OS         string          `json:"os"`
	Arch       string          `json:"arch"`
	GoVersion  string          `json:"goVersion"`
	CPUs       int             `json:"cpus"`
	GOMAXPROCS int             `json:"gomaxprocs"`
	Uptime     float64         `json:"uptime"`
	StartedAt  int64           `json:"startedAt"`
	IntervalMS int             `json:"intervalMS"`
}
