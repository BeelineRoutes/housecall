/** ****************************************************************************************************************** **
    Customer related calls
    
** ****************************************************************************************************************** **/

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

// Returns a list of the customers using a 'simple' searching keyword or short address part
// returns in order of most recently created.
// converted to page through the results, allows us to do an empty search for customers
func (this *HouseCall) SearchCustomers (ctx context.Context, token, search string) ([]Customer, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "200")
    params.Set("sort_direction", "desc")
    params.Set("sort_by", "created_at")
    params.Set("q", search)

    var ret []Customer

    for i := 1; i <= 1000; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("page", fmt.Sprintf("%d", i)) // set our next page

        resp := customerListResponse{}
    
        errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("customers?%s", params.Encode()), header, nil, &resp)
        if err != nil { return nil, errors.WithStack(err) } // bail
        if errObj != nil { return nil, errObj.Err() } // something else bad

        // we're here, we're good
        ret = append (ret, resp.Customers...)
        
        if i >= resp.TotalPages { return ret, nil } // we finished
    }
    // this only happens if we have too many pages... shouldn't happen anyway, but return what we got
    return ret, nil 
}

// creates the customer and returns their id
func (this *HouseCall) CreateCustomer (ctx context.Context, token string, customer *Customer) error {

    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 
    header["Content-Type"] = "application/json; charset=utf-8"

    resp := &Customer{}

    errObj, err := this.send (ctx, http.MethodPost, "customers", header, customer, resp)
    if err != nil { return errors.WithStack(err) } // bail
    if errObj != nil { return errObj.Err() } // something else bad

    // return a shallow copy of this new user, we really just want the id
    *customer = *resp 
    return nil 
}
