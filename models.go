
package housecall 

import (
	"github.com/pkg/errors"

	//"fmt"
	"net/url"
	"time"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

const apiURL = "https://api.housecallpro.com"


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

type Pro struct {
	Id string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	FullName string `json:"full_name"`
	Initials string `json:"initials"`
	Email string `json:"email"`
	Mobile string `json:"mobile_number"`
	MessagingId string `json:"messaging_uuid"`
	Color string `json:"color_hex"`
	Avatar string `json:"avatar_url"`
	HasAvatar bool `json:"has_avatar"`
	OrgName string `json:"organization_name"`
	Admin bool `json:"is_admin"`
	Archived bool `json:"is_archived"`
}

//----- ADDRESSES ---------------------------------------------------------------------------------------------------------//

type Address struct {
	Id string `json:"id"`
	Street string `json:"street"`
	Street2 string `json:"street_line_2"`
	City string `json:"city"`
	State string `json:"state"`
	Zip string `json:"zip"`
	Country string `json:"country"`
	Lat string `json:"latitude"`
	Lon string `json:"longitude"`
	TimeZone string `json:"time_zone"`
	CustomerId string `json:"customer"`
	Customer Customer
}

//----- CUSTOMERS ---------------------------------------------------------------------------------------------------------//

type Customer struct {
	Id string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Email string `json:"email"`
	Mobile string `json:"mobile_number"`
	Home string `json:"home_number"`
	Work string `json:"work_number"`
	Company string `json:"company"`
	JobTitle string `json:"job_title"`
	Notes string `json:"notes"`
	Notifications bool `json:"notifications_enabled"`
	// addresses

}


//----- JOBS ---------------------------------------------------------------------------------------------------------//

type Job struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CustomerId string `json:"customer_id"`
	Customer Customer
	AddressId string `json:"address_id"`
	Address Address
	Note string `json:"note"`
	Invoice string `json:"invoice_number"`
	Balance int64 `json:"outstanding_balance"`
	Date time.Time `json:"scheduled_date"`
	Description string `json:"description"`
	Tags struct {
		Data []string `json:"data"`
	} `json:"tags"`
	Pros struct {
		Data[]string `json:"data"`
	} `json:"pros"`
	Schedule struct {
		Data struct {
			Start time.Time `json:"start_time"`
			End time.Time `json:"end_time"`
			Window int `json:"arrival_window_minutes"`
		}
	}
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

