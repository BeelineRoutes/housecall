
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestThirdJobs (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, only unscheduled ones
	jobs, err := hc.ListJobs (ctx, cfg.Token, time.Now())
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 job")
	assert.Equal (t, "job_94dec270539c4566be8b11173323ef5f", jobs[0].Id, "target job id")

	/*
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	*/
}

