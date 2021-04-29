# HouseCallPro
Package provides a wrapper for interacting with the HouseCallPro API.  Written in pure GoLang

`go get github.com/BeelineRoutes/housecall`

## House Call OAuth flow
[Offical documentation for the API can be found here](https://pro.housecallpro.com/docs/alpha/api.html).

First have your user go to 
`https://api.housecallpro.com/oauth/authorize?response_type=code&client_id=clientId&redirect_uri=https://your-url.com`

This will return to your redirect url with a "code" as a url param
`https://dev-api.geniemobile.net/housecall?code=urlParamCode`

## Usage
```go

import (
    "github.com/BeelineRoutes/housecall"
    "log"
)

clientId := "ca96e4cd990507c2995b9633bd9caa679bee26e99f98572ba54751ab4ff24886" // your special client id
clientSecret := "1fd00f12ab1d3d13c6bf746aa1868bd591af098100d800195b56b6fa97795d73" // your secret
redirectUrl := "https://your-domain.com" // this can be whatever, you tell House Call what you want when you create your account

hc, err := NewHouseCall (clientId, clientSecret, redirectUrl) // create the hc object
if err != nil { log.Fatal (err) }

// now you can do something like convert the code you got from the url to a long-lived token and refresh token
params := r.URL.Query()
code := ""
if len(params["code"]) > 0 && len(params["code"][0]) > 0 {
    code = params[term][0]
} else {
    t.Fatal ("was expecting a url param 'code' to be set")
}

token, refresh, err := hc.TokensFromCode (context.TODO(), code)
if err != nil { log.Fatal (err) }

// token is used as the bearer for future calls
// refresh can be used to generate a new token when it expires
```

## History
This is actively being developed and while the goal is to prevent breaking changes, this is still in an early alpha.

- 0.1 Initial version allows for validating a client OAuth token and retrieving the Token and Refresh tokens


