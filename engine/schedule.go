package engine

import (
	"sync"

	"github.com/zhfxm/simple-crawler/collect"
	"go.uber.org/zap"
)

type Scheduler interface {
	Schedule()
	Push(...*collect.Request)
	Pull() *collect.Request
}

type Crawler struct {
	out         chan collect.ParseResult
	failures    map[string]*collect.Request
	failureLock sync.Mutex
	options
}

type Schedule struct {
	requestCh   chan *collect.Request
	workerCh    chan *collect.Request
	priReqQueue []*collect.Request
	reqQueue    []*collect.Request
	Logger      *zap.Logger
}

func NewSchedule(logger *zap.Logger) *Schedule {
	s := &Schedule{}
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	s.requestCh = requestCh
	s.workerCh = workerCh
	return s
}

func (s *Schedule) Schedule() {
	var req *collect.Request
	var ch chan *collect.Request
	for {
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerCh
		}
		if req == nil && len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			req = nil
			ch = nil
		}
	}
}

func (s *Schedule) Pull() *collect.Request {
	return <-s.workerCh
}

func (s *Schedule) Push(requests ...*collect.Request) {
	for _, req := range requests {
		s.requestCh <- req
	}
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOption
	for _, opt := range opts {
		opt(&options)
	}
	e := &Crawler{}
	e.options = options
	e.out = make(chan collect.ParseResult)
	e.failures = make(map[string]*collect.Request)
	return e
}

func (e *Crawler) Run() {
	// 启动调度器
	go e.Schedule()
	// 创建 worker
	for i := 0; i < e.WorkerCount; i++ {
		go e.CreateWorker()
	}
	// 处理结果
	e.HandlerResult()
}

func (e *Crawler) Schedule() {
	// 1、启动调度器
	go e.scheduler.Schedule()
	// 2、添加任务
	go e.scheduler.Push(e.Seeds...)
}

// 创建 worker
func (s *Crawler) CreateWorker() {
	for {
		r := s.scheduler.Pull()
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.SetFailuer(r)
			s.Logger.Error("can't fetch ", zap.Error(err), zap.String("url", r.Url))
			continue
		}
		res := r.ParseFunc(body, r)
		if len(res.Requests) > 0 {
			go s.scheduler.Push(res.Requests...)
		}
		if len(res.Items) > 0 {
			s.out <- res
		}
	}
}

func (e *Crawler) HandlerResult() {
	for {
		select {
		case result := <-e.out:
			for _, item := range result.Items {
				e.Logger.Sugar().Info("request result:", item.(string))
			}
		}
	}
}

func (e *Crawler) SetFailuer(req *collect.Request) {
	e.failureLock.Lock()
	defer e.failureLock.Unlock()
	if _, ok := e.failures[req.Unique()]; !ok {
		e.failures[req.Unique()] = req
	}
}