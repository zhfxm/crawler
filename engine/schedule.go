package engine

import (
	"time"

	"github.com/zhfxm/simple-crawler/collect"
	"go.uber.org/zap"
)

type Schedule struct {
	requestCh   chan *collect.Request
	workerCh    chan *collect.Request
	out         chan collect.ParseResult
	options
}

func NewSchedule(opts ...Option) *Schedule {
	options := defaultOption
	for _, opt := range opts {
		opt(&options)
	}
	s := &Schedule{}
	s.options = options
	return s
}

func (s *Schedule) Run() {
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

func (s *Schedule) schedule() {
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

func (s *Schedule) CreateWorker() {
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

func (s *Schedule) HandlerResult() {
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
