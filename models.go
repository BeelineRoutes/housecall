
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
	// TODO need to handle the different error objects returned using an overridden unmarshal function
	Error string `json:"error"`
	Description string `json:"error_description"`
	StatusCode int
	
	/*
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
	*/
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

//----- COMPANY ---------------------------------------------------------------------------------------------------------//

type Company struct {
	Id string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email string `json:"support_email"`
	Name string `json:"name"`
	Logo string `json:"logo_url"`
	Address Address `json:"address"`
	Website string `json:"website"`
	DefaultArrivalWindow int `json:"default_arrival_window"`
	TimeZone string `json:"time_zone"`
}

//----- PROS ---------------------------------------------------------------------------------------------------------//

type Employee struct {
	Id string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Mobile string `json:"mobile_number"`
	Color string `json:"color_hex"`
	Avatar string `json:"avatar_url"`
	Role string `json:"role"`
}

type employeeListResponse struct {
	Employees []Employee `json:"employees"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
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
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// this became the most complicated thing, but just trying to return an empty string when appropriate 
func (this Address) ToString() string {
	afterComma := ""
	
	if len(this.State) > 0 && len(this.Zip) > 0 { // we have both
		afterComma = fmt.Sprintf(", %s  %s",this.State, this.Zip)
	} else if len(this.Zip) > 0 {
		afterComma = " " + this.Zip 
	} else if len(this.State) > 0 {
		afterComma = " " + this.State
	}

	beforeComma := this.Street
	
	if len(beforeComma) > 0 && len (this.Street2) > 0 {
		beforeComma += " " + this.Street2
	} else if len(beforeComma) == 0 {
		beforeComma = this.Street2 // just copy this
	}

	if len(beforeComma) > 0 && len (this.City) > 0 {
		beforeComma += " " + this.City
	} else if len(beforeComma) == 0 {
		beforeComma = this.City // just copy this
	}

	if len(beforeComma) > 0 && len(afterComma) > 0 {
		return beforeComma + afterComma
	} else if len(beforeComma) > 0 {
		return beforeComma
	} else if len(afterComma) > 0 {
		return afterComma
	}
	return "" // nothing do'n
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
	WorkTimestamps struct {
		OnMyWay time.Time `json:"on_my_way_at"`
		Started time.Time `json:"started_at"`
		Completed time.Time `json:"completed_at"`
	} `json:"work_timestamps"`
}

/* known work status
needs scheduling
scheduled
in progress
complete unrated
complete rated
user canceled
pro canceled
*/

// returns that the job is in a state where the job is still expected to be completed in the future
func (this *Job) IsPending () bool {
	switch this.WorkStatus {
	case "scheduled", "needs scheduling":
		return true
	}
	return false // this is in a state where the job has been cancelled or already started
}

// returns that the job is in a state where everything is still a "go".  Either it hasn't happened yet, it's happening now, or it will in the future
func (this *Job) IsActive () bool {
	switch this.WorkStatus {
	case "scheduled", "in progress", "complete unrated", "complete rated":
		return true
	}
	return false // not an active job
}

type JobSchedule struct {
	Start time.Time `json:"start_time"`
	End time.Time `json:"end_time"`
	Window int `json:"arrival_window_in_minutes"`
	Notify bool `json:"notify"`
	NotifyPro bool `json:"notify_pro"`
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
