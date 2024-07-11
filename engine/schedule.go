package engine

import (
	"time"

	"github.com/zhfxm/simple-crawler/collect"
	"go.uber.org/zap"
)

type ScheduleEngine struct {
	requestCh   chan *collect.Request
	workerCh    chan *collect.Request
	out         chan collect.ParseResult
	WorkerCount int
	Logger      *zap.Logger
	Seeds       []*collect.Request
	Fetcher     collect.Fetcher
}

func (s *ScheduleEngine) Run() {
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	out := make(chan collect.ParseResult)
	s.requestCh = requestCh
	s.workerCh = workerCh
	s.out = out
	// 执行调度
	go s.schedule()
	// 创建 worker
	for i := 0; i < s.WorkerCount; i++ {
		go s.CreateWorker()
	}
	s.HandlerResult()
}

func (s *ScheduleEngine) schedule() {
	reqQueue := s.Seeds
	for {
		var req *collect.Request
		var ch chan *collect.Request
		if len(reqQueue) > 0 {
			req = reqQueue[0]
			ch = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			reqQueue = append(reqQueue, r)
		case ch <- req:
			reqQueue = reqQueue[1:]
		}
	}
}

func (s *ScheduleEngine) CreateWorker() {
	for {
		req := <-s.workerCh
		time.Sleep(1 * time.Second)
		body, err := s.Fetcher.Get(req)
		if err != nil {
			s.Logger.Error("request failed", zap.Error(err))
			continue
		}
		result := req.ParseFunc(body, req)
		s.out <- result
	}
}

func (s *ScheduleEngine) HandlerResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				s.requestCh <- req
			}
			for _, item := range result.Items {
				s.Logger.Sugar().Info("request result:", item.(string))
			}
		}
	}

}
