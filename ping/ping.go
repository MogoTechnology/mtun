package ping

import "C"
import (
	"crypto/tls"
	"encoding/json"
	"net"
	"sync"
	"time"
)

const factor = 2

func ping(url string, timeout int64) (int64, error) {
	start := time.Now()
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout: time.Duration(timeout) * time.Millisecond * factor,
	}, "tcp", url, &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return 0, err
	}
	defer conn.Close()

	end := time.Now()
	dur := end.Sub(start).Milliseconds() / factor
	return dur, nil
}

type Result struct {
	Ping  int64
	Score float64
}

func Ping(url string, proto int64, timeout int64, count int64) *Result {
	goUrl := url + ":1443"

	failCount := int64(0)
	sumPing := int64(0)
	for i := count; i > 0; i-- {
		p, err := ping(goUrl, timeout)
		if err != nil {
			failCount++
			continue
		}
		sumPing += p
	}

	avePing := int64(0)
	if failCount != count {
		avePing = sumPing / (count - failCount)
	}
	basicScore := float64(1)
	if avePing < 1000 {
		if avePing == 0 {
			basicScore = 0
		} else {
			basicScore = float64(1000-avePing) / float64(200)
		}
	}
	score := basicScore
	if failCount != 0 {
		score = basicScore - (float64(failCount)*2 - 0.5)
	}
	return &Result{
		Ping:  avePing,
		Score: score,
	}
}

type PingRequest struct {
	Url    string
	NodeID int64
	Proto  int64
}

type PingResponse struct {
	NodeID int64
	Url    string
	Ping   int64
	Score  float64
}

func PingMany(request string, timeout int64, count int64) (string, error) {
	requests := make([]PingRequest, 0)
	err := json.Unmarshal([]byte(request), &requests)
	if err != nil {
		return "", err
	}

	ret := make([]PingResponse, 0, len(requests))
	wg := sync.WaitGroup{}
	for i, req := range requests {
		wg.Add(1)
		go func(index int, req PingRequest) {
			defer wg.Done()
			r := Ping(req.Url, req.Proto, timeout, count)
			ret[index] = PingResponse{
				NodeID: req.NodeID,
				Url:    req.Url,
				Ping:   r.Ping,
				Score:  r.Score,
			}
		}(i, req)
	}
	wg.Wait()

	bytes, err := json.Marshal(ret)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type Manager struct {
	req     chan *PingRequest
	resp    chan *PingResponse
	wait    chan bool
	timeout int64
	count   int64
}

func NewManager(worker int64, timeout int64, count int64) *Manager {
	m := new(Manager)
	m.timeout = timeout
	m.count = count
	m.req = make(chan *PingRequest, worker)
	m.resp = make(chan *PingResponse, worker)
	m.wait = make(chan bool, worker)
	go m.Loop()
	return m
}

func (m *Manager) SetTimeout(timeout int64) {
	m.timeout = timeout
}

func (m *Manager) SetCount(count int64) {
	m.count = count
}

func (m *Manager) Loop() {
	for {
		select {
		case req := <-m.req:
			m.wait <- true

			go func() {
				defer func() {
					<-m.wait
				}()
				r := Ping(req.Url, req.Proto, m.timeout, m.count)
				m.resp <- &PingResponse{
					NodeID: req.NodeID,
					Url:    req.Url,
					Ping:   r.Ping,
					Score:  r.Score,
				}
			}()
		}
	}
}

func (m *Manager) AddRequest(ip string, nodeID int64, proto int64) {
	go func() {
		m.req <- &PingRequest{
			Url:    ip,
			NodeID: nodeID,
			Proto:  proto,
		}
	}()
}

func (m *Manager) WaitResult() (*PingResponse, error) {
	resp := <-m.resp
	return resp, nil
}
