package myqueue

import (
	"encoding/json"
	"time"

	"github.com/xblockchainlabs/myqueue/models"
)

func AddJob(name string, data interface{}, delay time.Duration) (id uint, err error) {
	params, err := json.Marshal(data)
	if err != nil {
		return
	}
	j := models.Job{}
	j.Name = name
	j.Params = params

	id, err = models.SchedJob(j, time.Now().UTC().Add(delay))
	return
}
