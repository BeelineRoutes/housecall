
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
	jobs, err := hc.ListUnscheduledJobs (ctx, cfg.AccessToken, 1)
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

	// get our list of jobs
	jobs, err := hc.ListJobs (ctx, cfg.AccessToken, time.Now().AddDate(0, 0, 0), time.Now().AddDate (0, 0, 1))
	if err != nil { t.Fatal (err) }

	t.Logf("got %d jobs\n", len(jobs))

	assert.Equal (t, true, len(jobs) > 0, "expecting at least 1 job")
	assert.NotEqual (t, "", jobs[0].Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Customer.Id, "not filled in")
	assert.NotEqual (t, "", jobs[0].Address.Id, "not filled in")
	
	
	for _, j := range jobs {
		t.Logf ("%+v\n", j)
	}
	
}

func TestThirdJobAppointments (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute * 10)
	defer cancel()

	// get our list of apps
	apps, err := hc.GetJobAppointments (ctx, cfg.AccessToken, "job_fa1846167bf54c8aa3615cb709c72129")
	if err != nil { t.Fatal (err) }

	t.Logf("got %d apps\n", len(apps))

	assert.Equal (t, true, len(apps) > 0, "expecting at least 1 app")
		
	for _, j := range apps {
		t.Logf ("%+v\n", j)
	}
	
}


func TestThirdJobScheduleUpdate1 (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, we need one of these to update
	jobs, err := hc.ListUnscheduledJobs (ctx, cfg.AccessToken, 1)
	if err != nil { t.Fatal (err) }

	if len(jobs) == 0 { t.Fatal ("need at least 1 job to do this test") }

	t.Logf("targing job %s\n", jobs[0].Id)

	// now update the schedule to be something
	targetDate := time.Now().AddDate (0, 0, 7) // 1 week in the future

	// get our list of employees
	employees, err := hc.ListEmployees (ctx, cfg.AccessToken)
	if err != nil { t.Fatal (err) }
	if len(employees) == 0 { t.Fatal ("you need an employee to assign a scheduled job to") }

	err = hc.UpdateJobSchedule (ctx, cfg.AccessToken, jobs[0].Id, append(make([]string, 0), employees[0].Id), targetDate, time.Minute * 33, time.Minute * 34, false) // weird things so we know we updated
	if err != nil { t.Fatal (err) }

	job, err := hc.GetJob (ctx, cfg.AccessToken, jobs[0].Id) // get this job to verify we updated it
	if err != nil { t.Fatal (err) }

	assert.Equal (t, targetDate.Format("2006-01-02 15:04:05"), job.Schedule.Start.Format("2006-01-02 15:04:05"), "start time")
	assert.Equal (t, targetDate.Add(time.Minute * 33).Format("2006-01-02 15:04:05"), job.Schedule.End.Format("2006-01-02 15:04:05"), "end time")
	assert.Equal (t, 34, job.Schedule.Window, "job window")

	// all good, now clear it

	err = hc.UpdateJobSchedule (ctx, cfg.AccessToken, jobs[0].Id, make([]string, 0), time.Time{}, 0, 0, false)
	if err != nil { t.Fatal (err) }

	job, err = hc.GetJob (ctx, cfg.AccessToken, jobs[0].Id) // get this job to verify we updated it
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, job.Schedule.Start.IsZero(), "start is zero")
	assert.Equal (t, true, job.Schedule.End.IsZero(), "end is zero")
}

func TestThirdJobLineItems (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of jobs, only unscheduled ones
	lineItems, err := hc.GetLineItems (ctx, cfg.AccessToken, "job_6d1066c319bf4617acfbb9cb038385fb")
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 2, len(lineItems))
	assert.Equal (t, "Tasting Flight", lineItems[0].Name)
	
	for _, li := range lineItems {
		t.Logf ("%+v\n", li)
	}
	
}

// job is deleted/archived so we should get a 410 back
func TestThirdJobArchivedScheduleUpdate (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// now update the schedule to be something
	targetDate := time.Now().AddDate (0, 0, 7) // 1 week in the future

	err := hc.UpdateJobSchedule (ctx, cfg.AccessToken, "job_a823caa00d064af0a0ef7c3f4f3fabc2", make([]string, 0), targetDate, time.Minute * 33, time.Minute * 30, false) // weird things so we know we updated
	if err != nil { t.Fatal (err) }
}

/*
func TestSimple (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	targetDate := time.Now().AddDate (0, 0, 2) // 1 week in the future

	err := hc.UpdateJobSchedule (ctx, cfg.AccessToken, "job_bed0d8b73e164e0a8be68b71603a9a5c", "pro_2a51082b07424ba9976da29c7d4fcbac", targetDate, time.Minute * 30, time.Minute * 30, false) // weird things so we know we updated
	if err != nil { t.Fatal (err) }

}


*/

// really just testing the notifications
func TestThirdJobScheduleUpdate2 (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	
	err := hc.UpdateJobSchedule (ctx, cfg.AccessToken, "job_9436795e1e2645fa988feab850f95b34", append(make([]string, 0), "pro_2a51082b07424ba9976da29c7d4fcbac"), time.Now(), time.Minute * 30, time.Minute * 60, false) // weird things so we know we updated
	if err != nil { t.Fatal (err) }

}
