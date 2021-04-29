
package housecall 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "context"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

//----- OUATH -------------------------------------------------------------------------------------------------------//

// Takes the passed code we got from the params of the redirect url and converts it to long-live token and refresh token
func (this *HouseCall) ListJobs (ctx context.Context, token string) ([]Job, error) {
    ret := make([]Job, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    for i := 1; i <= 10000; i++ { // stay in a loop as long as we're pulling jobs
        var resp struct {
            Data struct {
                Data []Job `json:"data"`
            } `json:"data"`
            TotalCount int `json:"total_count"`
        }

        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("alpha/jobs?page=%d&page_size=100", i), header, nil, &resp)
        if err != nil { return nil, err } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Data.Data...) // add this to our list

        if len(ret) >= resp.TotalCount { return ret, nil } // we finished
    }
    return ret, errors.Errorf ("received over %d jobs in your history", len(ret) * 100)
}
