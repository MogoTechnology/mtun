// Package ping 实现网络ping工具包，主要功能是测试目标服务器的网络连接延迟并计算评分。
//
// 评分逻辑
//   - 基础评分：基于平均延迟计算，延迟越低评分越高（1000ms为基准线）
//   - 失败惩罚：失败次数越多，评分越低
package ping

import "C"
import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"
)

const factor = 2

// ping返回单次连接延迟。
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

// Result 存储单次ping的结果，包含平均延迟(Ping)和评分(Score)。
type Result struct {
	Ping  int64
	Score float64
}

// Ping 通过TLS协议连接目标服务器的1443端口，测量网络延迟。
//
// 执行多次ping并计算平均值和评分。
//
// 参数：
//   - url：目标服务器的URL，域名或IP，不带端口号，例如 "www.example.com"。
//   - proto：协议类型，没用到。
//   - timeout：超时时间，单位毫秒。
//   - count：测试次数。
//
// 返回值：
//   - *Result：测试结果，包含延迟和评分。
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
	return &Result{
		Ping:  avePing,
		Score: getScore(count, failCount, avePing),
	}
}

// 结果大概在 0-5 分，5分最优。有失败时会出现负分数结果。
func getScore(count int64, failCount int64, avePing int64) float64 {
	basicScore := float64(1)
	if avePing < 1000 {
		if avePing == 0 {
			basicScore = 0
		} else {
			basicScore = float64(1000-avePing) / float64(200)
		}
	}
	score := basicScore
	if failCount != 0 && count > 0 {
		// 失败率为 62.5% 时，扣除全部 5 分。
		// 0.5 是补偿，表示第一次失败的扣除量减少 0.5 分。
		failRatio := float64(failCount) / float64(count)
		penalty := float64(5)/0.625*failRatio - 0.5
		score = basicScore - penalty
	}
	return score
}

// PingRequest 是ping请求参数，包含目标URL、节点ID和协议。
type PingRequest struct {
	Url    string
	NodeID int64
	Proto  int64 // 没用到
}

// PingResponse 是ping响应结果，包含节点ID、URL、延迟和评分。
type PingResponse struct {
	NodeID int64
	Url    string
	Ping   int64
	Score  float64
}

// PingMany 并发ping多个目标并返回JSON格式结果。
//
// 参数：
//   - request：包含多个目标的JSON字符串，每个目标包含Url、NodeID和Proto字段。
//   - timeout：超时时间，单位毫秒。
//   - count：测试次数。
//
// 返回值：
//   - string：包含多个目标的Ping结果的JSON字符串，每个目标包含NodeID、Url、Ping和Score字段。
//   - error：如果解析JSON字符串或执行Ping测试失败，返回错误信息。
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

// Manager 是ping请求管理器，支持并发控制和结果处理。
type Manager struct {
	req     chan *PingRequest
	resp    chan *PingResponse
	wait    chan bool
	timeout int64
	count   int64
}

// NewManager 创建一个新的Manager实例。
// Manager提供并发ping请求的管理、配置和结果获取。
//
// 参数：
//   - worker：并发工作线程数。
//   - timeout：超时时间，单位毫秒。
//   - count：测试次数。
//
// 返回值：
//   - *Manager：新创建的Manager实例。
//
// 注意：
//   - 该方法会启动工作协程
//   - 调用AddRequest方法添加目标后，需要调用WaitResult方法等待结果。
//   - 可以通过SetTimeout和SetCount方法动态调整超时时间和测试次数。
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
			// wait 个数到达 worker 个数时，阻塞
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

// AddRequest 添加一个ping请求。
//
// 参数：
//   - ip：目标IP地址。
//   - nodeID：节点ID。仅用于原样返回。
//   - proto：协议。没用到。
//
// 注意：
//   - 该方法会异步执行ping请求。
//   - 调用WaitResult方法等待结果。
func (m *Manager) AddRequest(ip string, nodeID int64, proto int64) {
	go func() {
		m.req <- &PingRequest{
			Url:    ip,
			NodeID: nodeID,
			Proto:  proto,
		}
	}()
}

// WaitResult 等待请求完成并返回结果。
//
// 返回值：
//   - *PingResponse：请求的结果。
//   - error：如果超时，返回超时错误。
func (m *Manager) WaitResult() (*PingResponse, error) {
	timeout := time.Duration(m.timeout) * time.Millisecond * factor * time.Duration(m.count) * 2
	select {
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	case resp := <-m.resp:
		return resp, nil
	}
}
