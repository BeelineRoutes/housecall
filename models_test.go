
package housecall 

import (
	
	"github.com/stretchr/testify/assert"
	//"github.com/pkg/errors"

	"testing"
	"time"
)

// used for testing so I don't have my actual credentials in the repo
type testConfig struct {
	ClientId, ClientSecret, RedirectUrl, OAuthCode, Token string 
}

func TestModelsError (t *testing.T) {
	var err *Error 

	assert.Equal (t, nil, err.Err(), "nil for a nil object")

	// give it some memory
	err = &Error{}
	assert.NotEqual (t, nil, err.Err(), "Should return an error")
	
}

func TestModelsOAuthResponse (t *testing.T) {
	resp := oauthResponse {
		Expires: 1000,
		Created: 1619557886,
	}

	tm, _ := time.Parse ("2006-01-02 15:04:05", "2021-04-27 21:28:06")
	assert.Equal (t, tm.Unix(), resp.ExpiresAt().Unix(), "expiring time")
}


func TestModelsNewHouseCall (t *testing.T) {
	_, err := NewHouseCall ("", "", "")
	
	assert.NotEqual (t, nil, err, "should have errored")

	// make it work
	hc, err := NewHouseCall ("ca96e4cd990507c2995b9633bd9caa679bee26e99f98572ba54751ab4ff24886", 
		"1fd00f12ab1d3d13c6bf746aa1868bd591af098100d800195b56b6fa97795d73", "https://google.com")
	
	assert.Equal (t, nil, err, "should have worked")
	assert.NotEqual (t, nil, hc, "should have a valid object")

	req := hc.seedOAuth()
	assert.Equal (t, "ca96e4cd990507c2995b9633bd9caa679bee26e99f98572ba54751ab4ff24886", req.ClientId, "client id match")
	assert.Equal (t, "1fd00f12ab1d3d13c6bf746aa1868bd591af098100d800195b56b6fa97795d73", req.ClientSecret, "client secret match")
	assert.Equal (t, "https://google.com", req.Redirect, "redirect url match")

}
