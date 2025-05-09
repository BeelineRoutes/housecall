/** ****************************************************************************************************************** **
	Calls related to jobs

    There's a couple of filters used for requesting jobs, these have been broken out into their own functions.

    Updating jobs allows for setting a target time as well as an employee.
** ****************************************************************************************************************** **/

package housecall 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "net/url"
    "context"
    "time"
    "strings"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// Returns a list of the jobs that are marked as "unscheduled".
func (this *HouseCall) ListUnscheduledJobs (ctx context.Context, token string, pageLimit int) ([]*Job, error) {
    ret := make([]*Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("work_status[]", "unscheduled")
    params.Set("sort_direction", "desc")

    if pageLimit == 0 { pageLimit = 1 } // just to make it work
    
    for i := 1; i <= pageLimit; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        if resp.TotalPages > pageLimit {
            return nil, nil // we have too many pages and it would take too long to return them all, ~3 seconds per page request
        }

        // we're here, we're good
        ret = append (ret, resp.Jobs...)

        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, nil // we're done
}

// returns all jobs that are within our start and finish ranges
func (this *HouseCall) ListJobs (ctx context.Context, token string, start, finish time.Time) ([]*Job, error) {
    ret := make([]*Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    params.Set("scheduled_start_min", start.Format(time.RFC3339))
    params.Set("scheduled_start_max", finish.Format(time.RFC3339))
    params.Set("expand[]", "appointments")

    existingApps := make(map[string]struct{})
    
    for i := 1; i <= 10; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        for _, job := range resp.Jobs {
            for _, j := range expandJob (job, start, finish) {
                
                if len(j.Schedule.Appointments) > 0 {
                    if _, exists := existingApps[j.Schedule.Appointments[0].Id]; exists { continue } // skip this one

                    existingApps[j.Schedule.Appointments[0].Id] = struct{}{} // mark it for next time
                }

                ret = append (ret, j) // include this one
            }
        }
        
        if i >= resp.TotalPages { break } // we finished
    }

    // 2025-04-27 NT started doing this as a way to get the second appointment for a job. Still doesn't find middle appointments
    // now find the jobs that end on this date
    params.Del("scheduled_start_min")
    params.Del("scheduled_start_max")

    params.Set("scheduled_end_min", start.Format(time.RFC3339))
    params.Set("scheduled_end_max", finish.Format(time.RFC3339))

    for i := 1; i <= 10; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        for _, job := range resp.Jobs {
            for _, j := range expandJob (job, start, finish) {
                
                if len(j.Schedule.Appointments) > 0 {
                    if _, exists := existingApps[j.Schedule.Appointments[0].Id]; exists { continue } // skip this one

                    existingApps[j.Schedule.Appointments[0].Id] = struct{}{} // mark it for next time
                }

                ret = append (ret, j) // include this one
            }
        }
        
        if i >= resp.TotalPages { break } // we finished
    }
    
    return ret, nil // we're good
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

//----- APPOINTMENTS

// this is how we update the "new" setup for jobs where we have an appointment now
// 2023-10-18 notifications don't work with this endpoint, HCP says they're working on that 
func (this *HouseCall) UpdateJobAppointmentSchedule (ctx context.Context, token, jobId, apptId string, employeeIds []string, startTime time.Time, 
                                                        duration, arrivalWindow time.Duration, notifyCustomer bool) error {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    // updating
    var req struct {
        Start time.Time `json:"start_time"`
        End time.Time `json:"end_time"`
        Notify bool `json:"notify"`
        Window int `json:"arrival_window_minutes"`
        DispatchedEmployees []string `json:"dispatched_employees_ids"`
    }

    req.Start = startTime
    req.End = startTime.Add (duration)
    req.Window = int(arrivalWindow.Minutes())
    req.DispatchedEmployees = employeeIds
    req.Notify = notifyCustomer
    
    errObj, err := this.send (ctx, http.MethodPut, fmt.Sprintf("jobs/%s/appointments/%s", jobId, apptId), header, req, nil)
    if err != nil { return errors.WithStack(err) } // bail
    if errObj != nil { 
        if errObj.StatusCode == http.StatusGone || errObj.StatusCode == http.StatusNotFound {
            return nil // no big deal
        }

        // i'm also getting HouseCall Error : 400 : Archived job :
        // which happens when someone deletes the job and not just the appointment for the job
        // just going to hard code that string
        if errObj.StatusCode == http.StatusBadRequest && strings.Contains (errObj.Err().Error(), "Archived job") {
            return nil // also ignore these errors
        }
        
        // otherwise we're good with this error here
        return errObj.Err() // something else bad
    }

    if err != nil { return errors.WithStack(err) } // bail
    if errObj != nil { return errObj.Err() } // something else bad

    // we're here, we're good
    return nil
}

// creates a new job in the system
func (this *HouseCall) CreateJob (ctx context.Context, token, customerId, addressId string, 
                                    startTime time.Time, duration, arrivalWindow time.Duration, 
                                    employeeIds, tags []string, lineItems []LineItem, leadSource, notes string) (*Job, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    header["Content-Type"] = "application/json; charset=utf-8"
    
    job := &createJob {
        CustomerId: customerId,
        AddressId: addressId,
        LineItems: lineItems,
        Tags: tags,
        LeadSource: leadSource,
        Notes: notes,
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
