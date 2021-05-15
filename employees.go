
package housecall 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "context"
    "net/url"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// Returns a list of the employees (pros) in asc by last name
func (this *HouseCall) ListEmployees (ctx context.Context, token string) ([]Employee, error) {
    ret := make([]Employee, 0) // main list to return
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "100")
    params.Set("sort_direction", "asc")
    params.Set("sort_by", "last_name")

    for i := 1; i <= 10000; i++ { // stay in a loop as long as we're pulling employees
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page
        resp := employeeListResponse{}
        
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("employees?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Employees...) // add this to our list

        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    return ret, errors.Errorf ("received over %d employees", len(ret))
}

