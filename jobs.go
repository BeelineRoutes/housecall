/** ****************************************************************************************************************** **
	Calls related to jobs

    There's a couple of filters used for requesting jobs, these have been broken out into their own functions.

    Updating jobs allows for setting a target time as well as an employee.
    Multiple employees may be assigned as well using UpdateJobDispatch
** ****************************************************************************************************************** **/

package housecall 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "net/url"
    "context"
    "time"
    "encoding/json"
    "log"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// Returns a list of the jobs that are marked as "unscheduled".
func (this *HouseCall) ListUnscheduledJobs (ctx context.Context, token string, pageLimit int) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("work_status[]", "unscheduled")
    params.Set("sort_direction", "desc")

    if pageLimit == 0 { pageLimit = 1 } // just to make it work
    if pageLimit > 200 { pageLimit = 200 } // let's not go crazy here
    
    for i := 1; i <= pageLimit; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Jobs...)

        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, nil // we're done
}

// gets a specific page of unscheduled jobs, ordered by most recently added
func (this *HouseCall) PageUnscheduledJobs (ctx context.Context, token string, page int) ([]Job, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("work_status[]", "unscheduled")
    params.Set("sort_direction", "desc")
    params.Set("page", fmt.Sprintf("%d", page)) // set our next page

    resp := jobListResponse{}

    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    return resp.Jobs, nil // we're good
}

// returns all jobs that are within our start and finish ranges
func (this *HouseCall) ListJobs (ctx context.Context, token string, start, finish time.Time) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    params.Set("scheduled_start_min", start.Format(time.RFC3339))
    params.Set("scheduled_start_max", finish.Format(time.RFC3339))
    
    for i := 1; i <= 100; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        for _, job := range resp.Jobs {
            jstr, _ := json.Marshal(job) // for error handling
            err = this.fillJobAppointments (ctx, token, &job, start, finish)

            if err == nil {
                if len(job.AssignedEmployees) > 0 {
                    ret = append (ret, job)
                } // no error for no crew members, it's expected with appointments
            } else {
                log.Printf("%v :: %s\n", err, string(jstr))
            }
        }
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Wrapf (ErrTooManyRecords, "received over %d jobs in your history", len(ret))
}

// returns all jobs that are within our start and finish ranges
func (this *HouseCall) ListMissedJobs (ctx context.Context, token string, start, finish time.Time) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    params.Set("scheduled_start_min", start.Format(time.RFC3339))
    params.Set("scheduled_start_max", finish.Format(time.RFC3339))
    // params.Set("work_status[]", "scheduled") // 2023-10-12 NT This is a big one, we can't use this status anymore because of appointments
    
    for i := 1; i <= 10; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        for _, job := range resp.Jobs {
            jstr, _ := json.Marshal(job) // for error handling

            err = this.fillJobAppointments (ctx, token, &job, start, finish)
            if err == nil {
                if len(job.AssignedEmployees) > 0 {
                    ret = append (ret, job)
                } // no error for no crew members, it's expected with appointments
            } else {
                log.Printf("%v :: %s\n", err, string(jstr))
            }
        }
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, nil // don't error about this
}

// returns a list of jobs for a specific employee over the target date range
func (this *HouseCall) ListJobsFromEmployee (ctx context.Context, token string, employeeId string, start, finish time.Time) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    params.Set("scheduled_start_min", start.Format(time.RFC3339))
    params.Set("scheduled_start_max", finish.Format(time.RFC3339))
    params.Set("employee_ids[]", employeeId)
    
    for i := 1; i <= 100; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        for _, job := range resp.Jobs {
            jstr, _ := json.Marshal(job) // for error handling

            err = this.fillJobAppointments (ctx, token, &job, start, finish)
            if err == nil {
                if len(job.AssignedEmployees) > 0 {
                    ret = append (ret, job)
                } // no error for no crew members, it's expected with appointments
            } else {
                log.Printf("%v :: %s\n", err, string(jstr))
            }
        }
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Wrapf (ErrTooManyRecords, "received over %d jobs in your history", len(ret))
}

// returns a list of jobs that are associated with the customer
func (this *HouseCall) ListJobsFromCustomer (ctx context.Context, token string, customerId string) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("customer_id", customerId)
    
    for i := 1; i <= 100; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        for _, job := range resp.Jobs {
            jstr, _ := json.Marshal(job) // for error handling

            err = this.fillJobAppointments (ctx, token, &job, job.Schedule.Start, job.Schedule.End)
            if err == nil {
                if len(job.AssignedEmployees) > 0 {
                    ret = append (ret, job)
                } // no error for no crew members, it's expected with appointments
            } else {
                log.Printf("%v :: %s\n", err, string(jstr))
            }
        }
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Wrapf (ErrTooManyRecords, "received over %d jobs in your history", len(ret))
}

// gets the info about a specific job
func (this *HouseCall) GetJob (ctx context.Context, token, jobId string) (*Job, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    job := &Job{}
    
    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs/%s", jobId), header, nil, job)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // see if there were appointments associated with this job
    job.Appointments, err = this.GetJobAppointments (ctx, token, jobId)

    // we're here, we're good
    return job, nil
}

// updates the target scheduled time for a job
// at least 1 employee is required for this
// if startTime is zero, then this will remove the scheduled time from the job
func (this *HouseCall) UpdateJobSchedule (ctx context.Context, token, jobId string, employeeIds []string, startTime time.Time, 
                                            duration, arrivalWindow time.Duration, notifyCustomer bool) error {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    
    if startTime.IsZero() {
        errObj, err := this.send (ctx, http.MethodDelete, fmt.Sprintf("jobs/%s/schedule", jobId), header, nil, nil)
        if err != nil { return errors.WithStack(err) } // bail
        if errObj != nil { 
            if errObj.StatusCode != http.StatusGone {
                return errObj.Err() // something else bad
            } // otherwise we're good with this error here
        }

    } else { // updating
        schedule := &JobSchedule {
            Start: startTime,
            End: startTime.Add (duration),
            Window: int(arrivalWindow.Minutes()),
            Notify: notifyCustomer,
        }

        // add in our assigned employee
        for _, id := range employeeIds {
            schedule.DispatchedEmployees = append (schedule.DispatchedEmployees, DispatchedEmployee{id}) 
        }

        errObj, err := this.send (ctx, http.MethodPut, fmt.Sprintf("jobs/%s/schedule", jobId), header, schedule, nil)
        if err != nil { return errors.WithStack(err) } // bail
        if errObj != nil { 
            if errObj.StatusCode != http.StatusGone {
                return errObj.Err() // something else bad
            } // otherwise we're good with this error here
        }
    }

    // we're here, we're good
    return nil
}

// sets the list of all assigned employees for a job
// only updates the list of employees assigned to a job
func (this *HouseCall) UpdateJobDispatch (ctx context.Context, token, jobId string, employeeIds ...string) error {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    
    dispatch := &JobDispatch {}

    // add in our employees
    for _, id := range employeeIds {
        dispatch.DispatchedEmployees = append (dispatch.DispatchedEmployees, DispatchedEmployee{id}) 
    }

    errObj, err := this.send (ctx, http.MethodPut, fmt.Sprintf("jobs/%s/dispatch", jobId), header, dispatch, nil)
    if err != nil { return errors.WithStack(err) } // bail
    if errObj != nil { 
        if errObj.StatusCode != http.StatusGone {
            return errObj.Err() // something else bad
        } // otherwise we're good with this error here
    }
    
    // we're here, we're good
    return nil
}

//----- APPOINTMENTS
// jobs can have appointments now...
func (this *HouseCall) GetJobAppointments (ctx context.Context, token, jobId string) ([]Appointment, error) {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    var resp struct {
        Appointments []Appointment
    }

    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs/%s/appointments", jobId), header, nil, &resp)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // we're here, we're good
    return resp.Appointments, nil
}

// for the most part we only want a specific appointment assigned to the job within the start/end times
// this picks the first one
func (this *HouseCall) fillJobAppointments (ctx context.Context, token string, job *Job, start, finish time.Time) error {
    if job.Schedule.End.Sub(job.Schedule.Start) < time.Hour * 4 { return nil } // going to assume if the duration of the job is short, then there's no appointments
    // just trying to save time by not checking every job for appointments

    // get a list of appointments
    apps, err := this.GetJobAppointments (ctx, token, job.Id)
    if err != nil { return err }

    if len(apps) < 2 { return nil } // just a long job i guess

    // find the first appointment that starts before the finish time and ends after the start time
    for _, app := range apps {
        if app.Start.Before(finish) && app.End.After(start) {
            // this is in our window
            job.Appointments = make([]Appointment, 1) // reset this list
            job.Appointments[0] = app // this one wins

            // i'm manually updating the status of this job, as the job has a single status, with multiple appointments
            // which makes no sense
            if app.End.Before(time.Now()) && job.IsActive() {
                job.WorkStatus = WorkStatus_scheduled // go back to a scheduled state
            }

            // the arrival window for this appointment also updates the jobs one
            job.Schedule.Window = app.Window
            // as do the start and end times
            job.Schedule.Start = app.Start 
            job.Schedule.End = app.End 
            
            // make sure the assigned crew members match this single appointment
            finalCrew := make([]Employee, 0)
            for _, emp := range job.AssignedEmployees {
                // make sure they're in this app list
                for _, id := range app.AssignedEmployees {
                    if emp.Id == id {
                        // this emp is still assigned
                        finalCrew = append(finalCrew, emp)
                        break // just for speed
                    }
                }
            }

            // copy over this final crew list for the job
            // IMPORTANT, this might be empty now! calling function needs to skip jobs with no assigned crew members
            job.AssignedEmployees = finalCrew
            return nil // we're good
        }
    }

    // this is actually bad, we couldn't find an appointment within this date range
    // i don't think this should happen... 
    return errors.WithStack (ErrAppNotFound)
}

func (this *HouseCall) UpdateJobAppointmentSchedule (ctx context.Context, token, jobId, optionId string, employeeIds []string, startTime time.Time, 
                                                        duration, arrivalWindow time.Duration) error {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    // updating
    var req struct {
        Start time.Time `json:"start_time"`
        End time.Time `json:"end_time"`
        Window int `json:"arrival_window_in_minutes"`
        DispatchedEmployees []string `json:"dispatched_employees_ids"`
    }

    req.Start = startTime
    req.End = startTime.Add (duration)
    req.Window = int(arrivalWindow.Minutes())
    req.DispatchedEmployees = employeeIds
    
    errObj, err := this.send (ctx, http.MethodPut, fmt.Sprintf("jobs/%s/appointments/%s", jobId, optionId), header, req, nil)
    if err != nil { return errors.WithStack(err) } // bail
    if errObj != nil { return errObj.Err() } // something else bad

    // we're here, we're good
    return nil
}

// creates a new job in the system
func (this *HouseCall) CreateJob (ctx context.Context, token, customerId, addressId string, 
                                    startTime time.Time, duration, arrivalWindow time.Duration, 
                                    employeeIds, tags []string, lineItems []LineItem, leadSource string) (*Job, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    header["Content-Type"] = "application/json; charset=utf-8"
    
    job := &createJob {
        CustomerId: customerId,
        AddressId: addressId,
        LineItems: lineItems,
        Tags: tags,
        LeadSource: leadSource,
    }

    // add in our employee
    for _, id := range employeeIds {
        job.Employees = append (job.Employees, id) 
    }
    
    job.Schedule.Start = startTime
    job.Schedule.End = startTime.Add (duration)
    job.Schedule.Window = fmt.Sprintf ("%d", int(arrivalWindow.Minutes()))

    resp := &Job{}
    
    errObj, err := this.send (ctx, http.MethodPost, "jobs", header, job, resp)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad
    
    // we're here, we're good
    return resp, nil
}

// returns the line items associated with the job
// this will include all the different "kinds"
func (this *HouseCall) GetLineItems (ctx context.Context, token, jobId string) ([]*LineItem, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    var ret struct {
        Data []*LineItem
    }
    
    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs/%s/line_items", jobId), header, nil, &ret)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // we're here, we're good
    return ret.Data, nil
}
