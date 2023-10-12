
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestThirdEmployees (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of employees
	employees, err := hc.ListEmployees (ctx, cfg.AccessToken)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(employees) > 0, "expecting at least 1 employee")
	assert.NotEqual (t, "", employees[0].Id, "not filled in")
	assert.NotEqual (t, "", employees[0].FirstName, "not filled in")
	assert.NotEqual (t, "", employees[0].LastName, "not filled in")
	assert.NotEqual (t, "", employees[0].Email, "not filled in")
	assert.NotEqual (t, "", employees[0].Mobile, "not filled in")
	assert.NotEqual (t, "", employees[0].Color, "not filled in")
	
	/*
	for _, e := range employees {
		t.Logf ("%+v\n", e)
	}
	*/
}

