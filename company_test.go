
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestThirdCompany (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, only unscheduled ones
	company, err := hc.Company (ctx, cfg.Token)
	if err != nil { t.Fatal (err) }

	// t.Logf("%+v\n", company)

	assert.Equal (t, true, len(company.PhoneNumber) > 0, "phone number: " + company.PhoneNumber)
	assert.Equal (t, true, len(company.Website) > 0, "website: " + company.Website)
	assert.Equal (t, true, len(company.TimeZone) > 0, "time zone: " + company.TimeZone)
	assert.Equal (t, true, len(company.Address.City) > 0, "city: " + company.Address.City)
	assert.Equal (t, true, len(company.Address.Longitude) > 0, "longitude: " + company.Address.Longitude)
}

func TestThirdSchedule (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, only unscheduled ones
	sch, err := hc.Schedule (ctx, cfg.Token)
	if err != nil { t.Fatal (err) }

	// t.Logf("%+v\n", sch)

	loc, err := time.LoadLocation ("America/Denver")
	if err != nil { t.Fatal (err) }

	list, err := sch.DaySchedules (loc)
	if err != nil { t.Fatal (err) }

	
	assert.Equal (t, 7, len(list), "expecting 7 days")
	assert.Equal (t, time.Weekday(1), list[1].Start.Weekday(), "Monday should be our second index")
	assert.Equal (t, time.Minute * 450, list[3].Duration, "wednesday shouild have 7.5 hours of work")
}

