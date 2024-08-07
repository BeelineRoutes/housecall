/** ****************************************************************************************************************** **
	Calls related to lead sources

** ****************************************************************************************************************** **/

package housecall 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "net/url"
    "context"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type LeadSource struct {
    Id, Name string 
}

type leadResponse struct {
    Total_pages int 
    Lead_sources []*LeadSource
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// returns a list of all the lead sources for an org
func (this *HouseCall) ListLeads (ctx context.Context, token string) ([]*LeadSource, error) {
    ret := make([]*LeadSource, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "100")
    
    for i := 1; i <= 3; i++ { // loop a little but no sense in going crazy here
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        
        resp := &leadResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("lead_sources?%s", params.Encode()), header, nil, resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append(ret, resp.Lead_sources...)
        
        if i >= resp.Total_pages { return ret, nil } // we finished
    }
    return ret, nil 
}
