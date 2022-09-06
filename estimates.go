/** ****************************************************************************************************************** **
	Calls related to estimates

    There's a couple of filters used for requesting estimates, very similar to jobs

    Updating 
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

// Returns a list of the estimates that are marked as "unscheduled".
func (this *HouseCall) ListUnscheduledEstimates (ctx context.Context, token string, pageLimit int) ([]Estimate, error) {
    ret := make([]Estimate, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("work_status[]", "unscheduled")

    if pageLimit == 0 { pageLimit = 1 } // just to make it work
    if pageLimit > 200 { pageLimit = 200 } // let's not go crazy here
    
    for i := 1; i <= pageLimit; i++ { // stay in a loop as long as we're pulling estimates
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := estimateListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("estimates?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Estimates...)

        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, nil // we're done
}

// returns a list of estimates for a specific employee over the target date range
func (this *HouseCall) ListEstimates (ctx context.Context, token string, employeeId string, start, finish time.Time) ([]Estimate, error) {
    ret := make([]Estimate, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    params.Set("scheduled_start_min", start.Format(time.RFC3339))
    params.Set("scheduled_start_max", finish.Format(time.RFC3339))
    if len(employeeId) > 0 {
        params.Set("employee_ids[]", employeeId)
    }
    
    for i := 1; i <= 100; i++ { // stay in a loop as long as we're pulling estimates
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := estimateListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("estimates?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Estimates...)
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Wrapf (ErrTooManyRecords, "received over %d estimates in your history", len(ret))
}

// gets the info about a specific estimate
func (this *HouseCall) GetEstimate (ctx context.Context, token, estId string) (*Estimate, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    est := &Estimate{}
    
    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("estimates/%s", estId), header, nil, est)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // we're here, we're good
    return est, nil
}

// updates the target estimate time for a estimate
// at least 1 employee is required for this
// if startTime is zero, then this will remove the scheduled time from the estimate
/*
func (this *HouseCall) UpdateEstimateSchedule (ctx context.Context, token, estId string, employeeIds []string, startTime time.Time, 
                                            duration, arrivalWindow time.Duration, notifyCustomer bool) error {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    
    if startTime.IsZero() {
        errObj, err := this.send (ctx, http.MethodDelete, fmt.Sprintf("estimates/%s/schedule", estId), header, nil, nil)
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
        for _, id := range employeeIds {
            schedule.DispatchedEmployees = append (schedule.DispatchedEmployees, DispatchedEmployee{id}) 
        }

        errObj, err := this.send (ctx, http.MethodPut, fmt.Sprintf("estimates/%s/schedule", estId), header, schedule, nil)
        if err != nil { return errors.WithStack(err) } // bail
        if errObj != nil { return errObj.Err() } // something else bad
    }

    // we're here, we're good
    return nil
}
*/

// creates a new estimate in the system
func (this *HouseCall) CreateEstimate (ctx context.Context, token, customerId, addressId string, 
                                    startTime time.Time, duration, arrivalWindow time.Duration, notifyCustomer bool,
                                    employeeIds, tags []string, leadSource, note, message string, 
                                    options []CreateEstimateOption) (*Estimate, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    header["Content-Type"] = "application/json; charset=utf-8"
    
    est := &createEstimate {
        Note: note,
        Message: message,
        CustomerId: customerId,
        AddressId: addressId,
        Tags: tags,
        LeadSource: leadSource,
        Employees: employeeIds,
        Options: options,
    }

    est.Schedule.Start = startTime
    est.Schedule.End = startTime.Add (duration)
    est.Schedule.Window = fmt.Sprintf ("%d", int(arrivalWindow.Minutes()))
    est.Schedule.NotifyCustomer = notifyCustomer

    resp := &Estimate{}
    
    errObj, err := this.send (ctx, http.MethodPost, "estimates", header, est, resp)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad
    
    // we're here, we're good
    return resp, nil
}

