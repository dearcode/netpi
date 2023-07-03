package proxy

import (
	"dearcode.net/crab/log"
	"dearcode.net/doodle/pkg/service"

	"dearcode.net/netpi/pkg/meta"
	"dearcode.net/netpi/pkg/pool"
)

type Job struct {
	meta.Job
}

type JobRequest struct {
	service.RequestHeader
}

type JobResponse struct {
	service.ResponseHeader
	Jobs []Job
}

func (j Job) Get(req JobRequest, resp *JobResponse) {
	for _, job := range pool.Factory.Jobs() {
		resp.Jobs = append(resp.Jobs, Job{Job: job})
	}

	log.Infof("%v req:%#v, resp:%#v", j, req, resp)
}
