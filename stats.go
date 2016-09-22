package main

import ()

type Stats struct {
	times []float64
}

func (s *Stats) AddTime(t float64) {
	s.times = append(s.times, t)
}

func (s *Stats) Count() int {
	return len(s.times)
}

func (s *Stats) Sum() float64 {
	sum := 0.0
	for _, v := range s.times {
		sum += v
	}
	return sum
}
