
package housecall 

import (
    "github.com/pkg/errors"
    
    "net/http"
    "context"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// Gets the info about our current company
func (this *HouseCall) Company (ctx context.Context, token string) (*Company, error) {
    company := &Company{}
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + token 

    errObj, err := this.send (ctx, http.MethodGet, "company", header, nil, &company)
    if err != nil { return nil, errors.WithStack(err) } // bail
    if errObj != nil { return nil, errObj.Err() } // something else bad

    // we're here, we're good
    return company, nil 
}

// Gets the companies schedule.  The return from HCP is a little... tough to interpret
// this just returns time.Times for the start and end of the longest time during the day
func (this *HouseCall) Schedule (ctx context.Context, token string) (*Schedule, error) {
  sch := &Schedule{}
  header := make(map[string]string)
  header["Authorization"] = "Bearer " + token 

  errObj, err := this.send (ctx, http.MethodGet, "company/schedule_availability", header, nil, sch)
  if err != nil { return nil, errors.WithStack(err) } // bail
  if errObj != nil { return nil, errObj.Err() } // something else bad

  // we're here, we're good
  return sch, nil 
}

