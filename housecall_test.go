
package housecall 

import (
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"

	"testing"
	"context"
	"os"
	"encoding/json"
	"time"
	"io/ioutil"
)

func newHouseCall (t *testing.T) (*HouseCall, *testConfig) {
	// read our local config
	config, err := os.Open("test.cfg")
	if err != nil { t.Fatal (err) }

	cfg := &testConfig{}

	jsonParser := json.NewDecoder (config)
	err = jsonParser.Decode (cfg)
	if err != nil { t.Fatal (err) }
	
	hc, err := NewHouseCall (cfg.ClientId, cfg.ClientSecret, cfg.RedirectUrl)
	if err != nil { t.Fatal (err) }

	return hc, cfg
}


func TestSecondHouseCall (t *testing.T) {
	// read our local config
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// generate a token
	oauth, err := hc.TokensFromCode (ctx, cfg.OAuthCode)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 64, len(oauth.AccessToken), "expecting a 64 char hash for the token")
	assert.Equal (t, 64, len(oauth.RefreshToken), "expecting a 64 char hash for the refresh token")

	t.Logf ("Using Token: %s :: Refresh Token: %s", oauth.AccessToken, oauth.RefreshToken)

	cfg.AccessToken = oauth.AccessToken // copy this back out to our config file
	out, _ := json.MarshalIndent (cfg, "", "    ")
	err = ioutil.WriteFile ("test.cfg", out, 0666)
	if err != nil { t.Fatal (err) }

	// now try to refresh the token
	oauth, err = hc.TokensFromRefresh (ctx, oauth.RefreshToken)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 64, len(oauth.AccessToken), "expecting a 64 char hash for the token")
	assert.Equal (t, 64, len(oauth.RefreshToken), "expecting a 64 char hash for the refresh token")
}

// tests an error from a bad token on refresh
func TestSecondHouseCallRefreshBad (t *testing.T) {
	// read our local config
	hc, _ := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// now try to refresh the token
	_, err := hc.TokensFromRefresh (ctx, "8a1270c05adeceb877d45965ed8d048bd071ad28084ea5ae952094a8ebfeaa49")
	// t.Logf("%s\n", errors.Cause(err))
	assert.Equal (t, ErrInvalidCode, errors.Cause(err))
}

