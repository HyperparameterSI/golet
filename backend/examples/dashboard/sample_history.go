// Author: HyperparameterSI

package main

const historyLimit = 300

// sampleHistory is a fixed-capacity ring. Its owner must synchronize access.
// Samples are stored by value so replacing an entry releases no external data.
type sampleHistory struct {
	samples  [historyLimit]ProcessSample
	next     int
	count    int
	sequence uint64
}

func (h *sampleHistory) append(sample ProcessSample) {
	h.sequence++
	sample.Sequence = h.sequence
	h.samples[h.next] = sample
	h.next = (h.next + 1) % historyLimit
	h.count = min(h.count+1, historyLimit)
}

// since copies only samples newer than sequence, in collection order. A cursor
// from a prior process falls back to a full snapshot after the process restarts.
func (h *sampleHistory) since(sequence uint64) []ProcessSample {
	if sequence > h.sequence {
		sequence = 0
	}
	count := h.count
	if remaining := h.sequence - sequence; remaining < uint64(count) {
		count = int(remaining)
	}
	result := make([]ProcessSample, count)
	start := (h.next - count + historyLimit) % historyLimit
	first := min(count, historyLimit-start)
	copy(result, h.samples[start:start+first])
	copy(result[first:], h.samples[:count-first])
	return result
}
