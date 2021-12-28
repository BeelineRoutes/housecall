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

    if pageLimit == 0 { pageLimit = 1 } // just to make it work
    
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
    return ret, errors.Errorf ("received over %d jobs in your history", len(ret))
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
    
    for i := 1; i <= 1000; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Jobs...)
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Errorf ("received over %d jobs in your history", len(ret))
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
    
    for i := 1; i <= 1000; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Jobs...)
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Errorf ("received over %d jobs in your history", len(ret))
}

// gets the info about a specific job
func (this *HouseCall) GetJob (ctx context.Context, token, jobId string) (*Job, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    job := &Job{}
    
    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs/%s", jobId), header, nil, job)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // we're here, we're good
    return job, nil
}

// updates the target scheduled time for a job
// at least 1 employee is required for this
// if startTime is zero, then this will remove the scheduled time from the job
func (this *HouseCall) UpdateJobSchedule (ctx context.Context, token, jobId, employeeId string, startTime time.Time, 
                                            duration, arrivalWindow time.Duration, notifyCustomer bool) error {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    
    if startTime.IsZero() {
        errObj, err := this.send (ctx, http.MethodDelete, fmt.Sprintf("jobs/%s/schedule", jobId), header, nil, nil)
        if err != nil { return errors.WithStack(err) } // bail
        if errObj != nil { return errObj.Err() } // something else bad

    } else { // updating
        schedule := &JobSchedule {
            Start: startTime,
            End: startTime.Add (duration),
            Window: int(arrivalWindow.Minutes()),
            Notify: notifyCustomer,
        }

        // add in our assigned employee
        schedule.DispatchedEmployees = append (schedule.DispatchedEmployees, DispatchedEmployee{employeeId}) 

        errObj, err := this.send (ctx, http.MethodPut, fmt.Sprintf("jobs/%s/schedule", jobId), header, schedule, nil)
        if err != nil { return errors.WithStack(err) } // bail
        if errObj != nil { return errObj.Err() } // something else bad
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
    if errObj != nil { return errObj.Err() } // something else bad
    
    // we're here, we're good
    return nil
}
