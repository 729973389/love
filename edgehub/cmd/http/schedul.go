package http

import "time"

type Schedule struct {
	SendData string
	Action   chan string
}

func NewSchedule() *Schedule {
	return &Schedule{GetConfig().SendData, make(chan string, 10)}
}

func (s *Schedule) Run() {
	go func() {
		for {
			s.Action <- s.SendData
			time.Sleep(30 * time.Minute)
		}
	}()
}
