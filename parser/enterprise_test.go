package main

import "testing"

func TestRun(t *testing.T) {
	run()
}

func BenchmarkRun(b *testing.B) {
	for b.Loop() {
		run()
	}
}
