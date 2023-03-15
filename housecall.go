/** ****************************************************************************************************************** **
    HouseCall Pro API wrapper
    written for GoLang
    Created 2021-04-27 by Nathan Thomas 
    Courtesy of BeelineRoutes.com

    current docs in v1
    https://docs.housecallpro.com/

** ****************************************************************************************************************** **/

package housecall 

import (
    "github.com/pkg/errors"
    
    //"fmt"
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
func (this *HouseCall) TokensFromCode (ctx context.Context, code string) (*OauthResponse, error) {
    req := this.seedOAuth()
    if req == nil { return nil, errors.WithStack(ErrInvalidCode) }

    req.Code = code // copy this over
    req.GrantType = "authorization_code"

    resp := &OauthResponse{}

    header := make(map[string]string)
    header["Content-Type"] = "application/json"

    // make our call
    errObj, err := this.send (ctx, http.MethodPost, "oauth/token", header, req, resp)
    if err != nil { return nil, err } // bail

    return resp, errObj.Err() // if errObj is nil, this will return a nil error
}

// Gets new tokens using a previously retreived refresh token
func (this *HouseCall) TokensFromRefresh (ctx context.Context, refresh string) (*OauthResponse, error) {
    req := this.seedOAuth()
    if req == nil { return nil, errors.WithStack(ErrInvalidCode) }
    req.RefreshToken = refresh // copy this over
    req.GrantType = "refresh_token"

    resp := &OauthResponse{}

    header := make(map[string]string)
    header["Content-Type"] = "application/json"

    // make our call
    errObj, err := this.send (ctx, http.MethodPost, "oauth/token", header, req, resp)
    if err != nil { return nil, err } // bail

    return resp, errObj.Err() // if errObj is nil, this will return a nil error
}
