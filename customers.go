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
// returns in order of most recently created. limiting this to 100 matches, if they don't search well, no reason to return all of them
func (this *HouseCall) SearchCustomers (ctx context.Context, token, search string) ([]Customer, error) {
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    params := url.Values{}
    params.Set("page_size", "100")
    params.Set("sort_direction", "desc")
    params.Set("sort_by", "created_at")
    params.Set("q", search)

    resp := customerListResponse{}
    
    errObj, err := this.send (ctx, http.MethodGet, fmt.Sprintf("customers?%s", params.Encode()), header, nil, &resp)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // we're here, we're good
    return resp.Customers, nil 
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
