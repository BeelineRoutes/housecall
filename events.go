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

// returns a list of events over the target date range
// this includes any event that overlaps the passed time
func (this *HouseCall) ListEvents (ctx context.Context, token string, start, end time.Time) ([]Event, error) {
    ret := make([]Event, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    
    for i := 1; i <= 10; i++ { // stay in a loop as long as we're pulling estimates
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := eventListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("events?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        // make sure this event fits within the range
        for _, event := range resp.Events {
            if event.Schedule.Start.Before(end) && event.Schedule.End.After(start) {
                ret = append (ret, event)
            }
        }

        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Errorf ("received over %d events in your history", len(ret))
}

