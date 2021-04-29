
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestHouseCallSecondJobs (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs
	jobs, err := hc.ListJobs (ctx, cfg.Token)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 job")
	assert.Equal (t, "job_87a68ea3cd514c869cbe0d0c8b30c7f1", jobs[0].Id, "target job id")

	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
}

