
package housecall 

import (
	"github.com/stretchr/testify/assert"

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
	token, refresh, err := hc.TokensFromCode (ctx, cfg.OAuthCode)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 64, len(token), "expecting a 64 char hash for the token")
	assert.Equal (t, 64, len(refresh), "expecting a 64 char hash for the refresh token")

	t.Logf ("Using Token: %s :: Refresh Token: %s", token, refresh)

	cfg.Token = token // copy this back out to our config file
	out, _ := json.MarshalIndent (cfg, "", "    ")
	err = ioutil.WriteFile ("test.cfg", out, 0666)
	if err != nil { t.Fatal (err) }

	// now try to refresh the token
	token, refresh, err = hc.TokensFromRefresh (ctx, refresh)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 64, len(token), "expecting a 64 char hash for the token")
	assert.Equal (t, 64, len(refresh), "expecting a 64 char hash for the refresh token")
}
