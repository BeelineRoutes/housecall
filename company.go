
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

