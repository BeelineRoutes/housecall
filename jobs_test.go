
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
	jobs, err := hc.ListUnscheduledJobs (ctx, cfg.Token)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 job")
	assert.NotEqual (t, "", jobs[0].Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Customer.Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Address.Id, "not filled in")
	
	/*
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	*/
}

func TestThirdFutureJobs (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute * 10)
	defer cancel()

	// get our list of jobs, only unscheduled ones
	jobs, err := hc.ListJobs (ctx, cfg.Token, time.Now(), time.Now().AddDate (0, 2, 0))
	if err != nil { t.Fatal (err) }

	t.Logf("got %d jobs\n", len(jobs))

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 job")
	assert.NotEqual (t, "", jobs[0].Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Customer.Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Address.Id, "not filled in")
	
	/*
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	*/
}


func TestThirdJobScheduleUpdate (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, we need one of these to update
	jobs, err := hc.ListUnscheduledJobs (ctx, cfg.Token)
	if err != nil { t.Fatal (err) }

	if len(jobs) == 0 { t.Fatal ("need at least 1 job to do this test") }

	t.Logf("targing job %s\n", jobs[0].Id)

	// now update the schedule to be something
	targetDate := time.Now().AddDate (0, 0, 7) // 1 week in the future

	// get our list of employees
	employees, err := hc.ListEmployees (ctx, cfg.Token)
	if err != nil { t.Fatal (err) }
	if len(employees) == 0 { t.Fatal ("you need an employee to assign a scheduled job to") }

	err = hc.UpdateJobSchedule (ctx, cfg.Token, jobs[0].Id, employees[0].Id, targetDate, time.Minute * 33, time.Minute * 34, false) // weird things so we know we updated
	if err != nil { t.Fatal (err) }

	job, err := hc.GetJob (ctx, cfg.Token, jobs[0].Id) // get this job to verify we updated it
	if err != nil { t.Fatal (err) }

	assert.Equal (t, targetDate.Format("2006-01-02 15:04:05"), job.Schedule.Start.Format("2006-01-02 15:04:05"), "start time")
	assert.Equal (t, targetDate.Add(time.Minute * 33).Format("2006-01-02 15:04:05"), job.Schedule.End.Format("2006-01-02 15:04:05"), "end time")
	assert.Equal (t, 34, job.Schedule.Window, "job window")

	// all good, now clear it

	err = hc.UpdateJobSchedule (ctx, cfg.Token, jobs[0].Id, "", time.Time{}, 0, 0, false)
	if err != nil { t.Fatal (err) }

	job, err = hc.GetJob (ctx, cfg.Token, jobs[0].Id) // get this job to verify we updated it
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, job.Schedule.Start.IsZero(), "start is zero")
	assert.Equal (t, true, job.Schedule.End.IsZero(), "end is zero")
}

/*
func TestSimple (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	targetDate := time.Now().AddDate (0, 0, 2) // 1 week in the future

	err := hc.UpdateJobSchedule (ctx, cfg.Token, "job_bed0d8b73e164e0a8be68b71603a9a5c", "pro_2a51082b07424ba9976da29c7d4fcbac", targetDate, time.Minute * 30, time.Minute * 30, false) // weird things so we know we updated
	if err != nil { t.Fatal (err) }

}


*/

