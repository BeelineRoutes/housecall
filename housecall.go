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
func (this *HouseCall) TokensFromCode (ctx context.Context, code string) (string, string, error) {
    req := this.seedOAuth()
    req.Code = code // copy this over
    req.GrantType = "authorization_code"

    resp := oauthResponse{}

    header := make(map[string]string)
    header["Content-Type"] = "application/json"

    // make our call
    errObj, err := this.send (ctx, http.MethodPost, "oauth/token", header, req, &resp)
    if err != nil { return "", "", err } // bail

    return resp.AccessToken, resp.RefreshToken, errObj.Err() // if errObj is nil, this will return a nil error
}

// Gets new tokens using a previously retreived refresh token
func (this *HouseCall) TokensFromRefresh (ctx context.Context, refresh string) (string, string, error) {
    req := this.seedOAuth()
    req.RefreshToken = refresh // copy this over
    req.GrantType = "refresh_token"

    resp := oauthResponse{}

    header := make(map[string]string)
    header["Content-Type"] = "application/json"

    // make our call
    errObj, err := this.send (ctx, http.MethodPost, "oauth/token", header, req, &resp)
    if err != nil { return "", "", err } // bail

    return resp.AccessToken, resp.RefreshToken, errObj.Err() // if errObj is nil, this will return a nil error
}
