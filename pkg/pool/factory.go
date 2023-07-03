package pool

import (
	"net"
	"sync"
	"time"

	"dearcode.net/crab/log"

	"dearcode.net/netpi/pkg/meta"
)

type poolType int

const (
	poolTypeJob = iota + 1
	poolTypeWorker
)

type connQueue struct {
	workerID string
	conns    []proxyConn
	last     time.Time
}

type proxyConn struct {
	conn net.Conn
	last time.Time
}

type factory struct {
	workers map[string]*connQueue
	jobs    map[string]*connQueue
	sync.Mutex
}

var Factory *factory

func init() {
	Factory = &factory{
		workers: make(map[string]*connQueue),
		jobs:    make(map[string]*connQueue),
	}
}

func (f *factory) Get(workerID string, t poolType) (proxyConn, bool) {
	f.Lock()
	defer f.Unlock()

	cq := f.workers
	if t == poolTypeJob {
		cq = f.jobs
	}

	q, ok := cq[workerID]
	if !ok {
		return proxyConn{}, false
	}

	if len(q.conns) == 0 {
		return proxyConn{}, false
	}

	c := q.conns[0]
	q.conns = q.conns[1:]
	return c, true

}

func (f *factory) GetWorker(workerID string) (net.Conn, bool) {
	w, ok := f.Get(workerID, poolTypeWorker)
	if !ok {
		return nil, false
	}
	return w.conn, true
}

func (f *factory) PutWorker(workerID string, conn net.Conn) {
	f.Put(workerID, conn, poolTypeWorker)
	log.Infof("worker:%v online worker:%v", conn, workerID)
}

func (f *factory) Put(workerID string, conn net.Conn, t poolType) {
	f.Lock()
	defer f.Unlock()

	cq := f.workers
	if t == poolTypeJob {
		cq = f.jobs
	}

	c := proxyConn{
		conn: conn,
		last: time.Now(),
	}

	q, ok := cq[workerID]
	if !ok {
		aq := &connQueue{
			workerID: workerID,
			conns:    []proxyConn{c},
			last:     time.Now(),
		}
		cq[workerID] = aq
		return
	}

	q.last = time.Now()

	q.conns = append(q.conns, c)
}

func (f *factory) Jobs() []meta.Job {
	f.Lock()
	defer f.Unlock()

	var js []meta.Job

	for id, job := range f.jobs {
		js = append(js, meta.Job{
			Time: job.last,
			ID:   id,
		})
	}

	return js
}

func (f *factory) Wait(workerID string, conn net.Conn) net.Conn {
	f.Put(workerID, conn, poolTypeJob)
	conn, ok := f.GetWorker(workerID)
	for !ok {
		log.Infof("workerID:%v get worker failed", workerID)
		conn, ok = f.GetWorker(workerID)
		time.Sleep(time.Second)
	}

	return conn
}
