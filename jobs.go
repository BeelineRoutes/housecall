
package housecall 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "context"
    "time"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// Returns a list of the jobs in desc order up until the previous Date.  Unscheduled jobs will always get included first
func (this *HouseCall) ListJobs (ctx context.Context, token string, previousDate time.Time) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    for i := 1; i <= 10000; i++ { // stay in a loop as long as we're pulling jobs
        resp := jobListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("jobs?page=%d&page_size=100&sort_direction=desc", i), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Jobs...) // add this to our list

        if i >= resp.TotalPages { return ret, nil } // we finished
        
        lastJob := resp.Jobs[len(resp.Jobs)-1]
        if lastJob.Schedule.Start.IsZero() == false && lastJob.Schedule.Start.Before (previousDate) {
            return ret, nil // we hit our previous date limit, so we're done
        }
    }
    return ret, errors.Errorf ("received over %d jobs in your history", len(ret))
}

