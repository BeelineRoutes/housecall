
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestThirdEvents (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	start, err := time.Parse("2006-01-02", "2022-08-28")
	if err != nil { t.Fatal (err) }

	end, err := time.Parse("2006-01-02", "2022-08-30")
	if err != nil { t.Fatal (err) }

	// get our list of jobs, only unscheduled ones
	events, err := hc.ListEvents (ctx, cfg.Token, start, end)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 1, len(events))
	assert.Equal (t, "evt_608c05349bda4bdf9e92529c07342c55", events[0].Id)
	
	/*
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	*/
}
