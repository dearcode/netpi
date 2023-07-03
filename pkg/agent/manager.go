package agent

import (
	"fmt"
	"time"

	"dearcode.net/crab/http/client"
	"dearcode.net/crab/log"

	"dearcode.net/netpi/pkg/meta"
)

type jobResp struct {
	Status int
	Jobs   []meta.Job
}

func (m *Manager) fetchJobs() ([]meta.Job, error) {
	resp := jobResp{}

	if err := client.New().GetJSON(m.fetchURL, nil, &resp); err != nil {
		log.Errorf("GetJSON error:%v, url:%v", err, m.fetchURL)
		return nil, err
	}
	log.Infof("fetch url:%v, resp:%#v", m.fetchURL, resp)

	if resp.Status != 0 {
		log.Errorf("fetch jobs status:%v", resp.Status)
		return nil, fmt.Errorf("invalid response code:%v", resp.Status)
	}

	return resp.Jobs, nil
}

type Manager struct {
	id        string
	localAddr string
	agentAddr string
	fetchURL  string
}

func NewManager(id, localAddr, proxyHost string) *Manager {
	return &Manager{
		id:        id,
		localAddr: localAddr,
		agentAddr: fmt.Sprintf("%s:9877", proxyHost),
		fetchURL:  fmt.Sprintf("http://%s:8080/proxy/Job/", proxyHost),
	}
}

func (m *Manager) filter(jobs []meta.Job) []meta.Job {
	var js []meta.Job
	for _, j := range jobs {
		if j.ID == m.id {
			js = append(js, j)
		}
	}

	return js
}

func (m *Manager) Run() {

	log.Infof("id:%v local:%v, proxy:%v", m.id, m.localAddr, m.agentAddr)

	for {
		jobs, err := m.fetchJobs()
		if err != nil {
			log.Errorf("fetch jobs error:%v", err)
			time.Sleep(time.Second)
			continue
		}

		for range m.filter(jobs) {
			s := newServer(m.localAddr, m.agentAddr, m.id)
			if err = s.Connect(); err != nil {
				log.Errorf("connect error:%v", err)
				continue
			}
			log.Infof("connect proxy:%v success", m.agentAddr)
			go s.Run()
		}
		time.Sleep(time.Second)
	}

}
