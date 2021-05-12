
package housecall 

import (
	"github.com/pkg/errors"

	"fmt"
	"net/url"
	"net/http"
	"time"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

const apiURL = "https://api.housecallpro.com"

//----- ERRORS ---------------------------------------------------------------------------------------------------------//

var (
	ErrInvalidCode 		= errors.New("OAuth code not valid")
)


  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

//----- ERRORS ---------------------------------------------------------------------------------------------------------//
type Error struct {
	Error string `json:"error"`
	Description string `json:"error_description"`
	StatusCode int
}

func (this *Error) Err () error {
	if this == nil { return nil } // no error
	switch this.StatusCode {
	case http.StatusUnauthorized:
		if this.Error == "invalid_grant" { // this is for granting access based on the passed code
			return errors.Wrap (ErrInvalidCode, this.Description)
		}

	}
	// just a default
	return errors.Errorf ("HouseCall Error : %d : %s : %s", this.StatusCode, this.Error, this.Description)
}

//----- OAUTH ---------------------------------------------------------------------------------------------------------//

type oauthRequest struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType string `json:"grant_type"`
	Code string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Redirect string `json:"redirect_uri"`
}

type oauthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Expires int64 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	Created int64 `json:"created_at"`
}

// returns a time object of when this oauth will expire
func (this *oauthResponse) ExpiresAt () time.Time {
	return time.Unix (this.Created + this.Expires, 0)
}

//----- PROS ---------------------------------------------------------------------------------------------------------//

type Employee struct {
	Id string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	Mobile string `json:"mobile_number"`
	Color string `json:"color_hex"`
	Avatar string `json:"avatar_url"`
	Role string `json:"role"`
}

//----- ADDRESSES ---------------------------------------------------------------------------------------------------------//

type Address struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Street string `json:"street"`
	Street2 string `json:"street_line_2"`
	City string `json:"city"`
	State string `json:"state"`
	Zip string `json:"zip"`
	Country string `json:"country"`
}

func (this Address) ToString() string {
	return fmt.Sprintf ("%s %s %s, %s  %s", this.Street, this.Street2, this.City, this.State, this.Zip)
}

//----- CUSTOMERS ---------------------------------------------------------------------------------------------------------//

type Customer struct {
	Id string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Mobile string `json:"mobile_number"`
	Home string `json:"home_number"`
	Work string `json:"work_number"`
	Company string `json:"company"`
	Notifications bool `json:"notifications_enabled"`
	Tags []string `json:"tags"`
	Addresses []Address `json:"addresses"`
}


//----- JOBS ---------------------------------------------------------------------------------------------------------//

type Job struct {
	Id string `json:"id"`
	CustomerId string `json:"customer_id"`
	Customer Customer
	Address Address `json:"address"`
	Note string `json:"note"`
	WorkStatus string `json:"work_status"`
	Invoice string `json:"invoice_number"`
	Balance int64 `json:"outstanding_balance"`
	Total int64 `json:"total_amount"`
	Tags []string `json:"tags"`
	Description string `json:"description"`
	AssignedEmployees [] Employee `json:"assigned_employees"`
	Schedule struct {
		Start time.Time `json:"scheduled_start"`
		End time.Time `json:"scheduled_end"`
		Window int `json:"arrival_window"`
	}
}

type jobListResponse struct {
	Jobs []Job `json:"jobs"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

//----- PUBLIC ---------------------------------------------------------------------------------------------------------//

type HouseCall struct {
	clientId, clientSecret, callbackUrl string // for making api calls
}

// populates our oauth request with the data we have from this object
func (this *HouseCall) seedOAuth () *oauthRequest {
	return &oauthRequest {
		ClientId: this.clientId,
		ClientSecret: this.clientSecret,
		Redirect: this.callbackUrl,
	}
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

func NewHouseCall (clientId, clientSecret, callbackUrl string) (*HouseCall, error) {
	// validate some inputs
	if len(clientId) != 64 { return nil, errors.Errorf ("client ID appears invalid.  Expecting a 64 character hash") }
	if len(clientSecret) != 64 { return nil, errors.Errorf ("client secret appears invalid.  Expecting a 64 character hash") }

	u, err := url.Parse(callbackUrl)
	if err != nil { return nil, errors.Wrapf (err, "%s is not a valid url", callbackUrl) }
	if u.Scheme == "" || u.Host == "" { return nil, errors.Errorf ("%s is not a valid url", callbackUrl) }
    
	// looks good
	return &HouseCall { 
		clientId: clientId, 
		clientSecret: clientSecret,
		callbackUrl: callbackUrl,
	}, nil
}

