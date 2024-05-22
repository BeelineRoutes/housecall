
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestThirdEstimatesUnscheduled (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, only unscheduled ones
	jobs, err := hc.ListUnscheduledEstimates (ctx, cfg.AccessToken, 1)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 estimate")
	assert.NotEqual (t, "", jobs[0].Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Customer.Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Address.Id, "not filled in")
	
	
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	
}
