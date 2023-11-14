/** ****************************************************************************************************************** **
	Data objects
	Converted objects from the HCP api into go-lang equivilants

** ****************************************************************************************************************** **/

package housecall 

import (
	"github.com/pkg/errors"

	"fmt"
	"net/url"
	"net/http"
	"time"
	"strings"
	"strconv"
	"encoding/json"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type WorkStatus string 

const (
	WorkStatus_needsScheduling 		WorkStatus = "needs scheduling"
	WorkStatus_scheduled	 		WorkStatus = "scheduled"
	WorkStatus_inProgress	 		WorkStatus = "in progress"
	WorkStatus_completeUnrated 		WorkStatus = "complete unrated"
	WorkStatus_completeRated 		WorkStatus = "complete rated"
	WorkStatus_userCanceled 		WorkStatus = "user canceled"
	WorkStatus_proCanceled 			WorkStatus = "pro canceled"
	
)

const apiURL = "https://api.housecallpro.com"

//----- ERRORS ---------------------------------------------------------------------------------------------------------//

var (
	ErrInvalidCode 		= errors.New("OAuth code not valid")
	ErrAuthExpired		= errors.New("OAuth expired")
	ErrTooManyRecords	= errors.New("Too many records returned")
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

//----- ERRORS ---------------------------------------------------------------------------------------------------------//
type Error struct {
	ErrMsg, Description string 
	StatusCode int 
}

func (this *Error) UnmarshalJSON (b []byte) error {

	// try this way
	var one struct {
		Error struct {
			Message string 
		}
	}

	err := json.Unmarshal(b, &one)
	if err == nil && len(one.Error.Message) > 0 {
		this.ErrMsg = one.Error.Message

		if strings.Contains(this.ErrMsg, "archived job") {
			this.StatusCode = http.StatusGone // it's gone 410
		}

	} else {
		// that didn't work, try another format
		var two struct {
			Error string 
			Description string `json:"error_description"`
			StatusCode int
		}

		err = json.Unmarshal (b, &two)
		if err == nil {
			this.ErrMsg = two.Error 
			this.Description = two.Description
			this.StatusCode = two.StatusCode
		} else {
			this.ErrMsg = err.Error()
			this.Description = string(b)
		}
	}

	if len(this.ErrMsg) == 0 {
		// this didn't work
		this.ErrMsg = "Unkown struct type"
		this.Description = string(b)
	}
	return nil 
}

func (this *Error) Err () error {
	if this == nil { return nil } // no error
	
	if this.ErrMsg == "invalid_grant" { // this is for granting access based on the passed code
		return errors.Wrap (ErrInvalidCode, this.Description)
	}

	switch this.StatusCode {
	case http.StatusUnauthorized:
		return errors.Wrap (ErrAuthExpired, this.Description) // invalid for another reason, most likely the oauth has been revoked
	
	}
	// just a default
	return errors.Errorf ("HouseCall Error : %d : %s : %s", this.StatusCode, this.ErrMsg, this.Description)
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

type OauthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Expires int64 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	Created int64 `json:"created_at"`
}

// returns a time object of when this oauth will expire
func (this *OauthResponse) ExpiresAt () time.Time {
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

// converts the companies timezone string into a golang location object
func (this Company) ConvertTimezone () (*time.Location, error) {
	if len(this.TimeZone) == 0 { this.TimeZone = "UTC" } // just a default

	return time.LoadLocation (this.TimeZone)
}

//----- SCHEDULE ---------------------------------------------------------------------------------------------------------//

type daySchedule struct {
	Start time.Time 
	Duration time.Duration 
}

type scheduleTime string

// converts this string into a golang time object
func (this scheduleTime) Time (loc *time.Location) (time.Time, error) {
	tm := time.Now() // just use the current calendar year, month, day for seeding these times to return

	parts := strings.Split (string(this), ":") // expecting a string "13:00"
	if len(parts) != 2 { 
		return time.Time{}, errors.Errorf ("bad start time : %s", this) 
	}

	hr, err := strconv.Atoi (parts[0]) // get the hours
	if err != nil { 
		return time.Time{}, errors.Errorf ("bad start time hour : %s : %s", this, err.Error()) 
	}

	min, err := strconv.Atoi (parts[1]) // get the minutes 
	if err != nil { 
		return time.Time{}, errors.Errorf ("bad start time minutes : %s : %s", this, err.Error()) 
	}

	// now create our actual start time, using our timezone
	return time.Date (tm.Year(), tm.Month(), tm.Day(), hr, min, 0, 0, loc), nil 
}

type Schedule struct {
	DailyAvailabilities struct {
		Data []struct {
			DayName string `json:"day_name"`
			ScheduleWindows struct {
				Data []struct {
					StartTime scheduleTime `json:"start_time"`
					EndTime scheduleTime `json:"end_time"`
				} `json:"data"`
			} `json:"schedule_windows"`
		} `json:"data"`
	} `json:"daily_availabilities"`
}

// goes through all the days and returns the earliest start and end for each
// leaves Start and End as IsZero if there's no schedules for that day
// this always returns 7 items, 1 for each day of the week
// loc is the local timezone for this schedule, HCP has the times as a local string
// when returned it converts it to UTC time
// if the schedule seems off I'll add more time to either side depending on how close to noon the times are
func (this *Schedule) DaySchedules (loc *time.Location) ([]daySchedule, error) {
	var ret []daySchedule // this is what we're going to try to fill in

	utc, _ := time.LoadLocation ("UTC")

	for d := 0; d < 7; d++ { // 7 days in a loop
		day := daySchedule{}
		dayName := "saturday" // default

		switch time.Weekday(d) { // figure out what our target date is in housecall
		case time.Sunday:
			dayName = "sunday"
		case time.Monday:
			dayName = "monday"
		case time.Tuesday:
			dayName = "tuesday"
		case time.Wednesday:
			dayName = "wednesday"
		case time.Thursday:
			dayName = "thursday"
		case time.Friday:
			dayName = "friday"
		}

		var early, late time.Time

		// loop through our data looking for the correct date
		for _, data := range this.DailyAvailabilities.Data {
			if strings.EqualFold (data.DayName, dayName) { // this is our target day of the week
				for _, list := range data.ScheduleWindows.Data {
					// parse out our start and end times
					start, err := list.StartTime.Time(loc)
					if err != nil {
						return nil, errors.Wrapf (err, "Weekday : %s", dayName)
					}
					
					end, err := list.EndTime.Time(loc) // and the end time
					if err != nil {
						return nil, errors.Wrapf (err, "Weekday : %s", dayName)
					}

					// what we're looking for is an earlier start time, and a later end time
					if early.IsZero() || start.Before(early) {
						early = start 
					}
					if late.IsZero() || end.After(late) {
						late = end 
					}
				}
			}
		}

		// at this point we know our ealiest start and latest end for the date, so set them in our daySchedule
		day.Start = early.In (utc) // keep everything in utc time
		day.Duration = late.Sub(early) // if these were never set, it still works out

		// now i want the Weekday to return the correct day, so we loop, adding days until it matches
		for day.Start.Weekday() != time.Weekday(d) {
			day.Start = day.Start.AddDate (0, 0, 1) // date only matters so it matches the day of the week for us
		}

		ret = append (ret, day) // add this to our return list
	}

	return ret, nil // we're good
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
	Tags []string `json:"tags"`
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
	LeadSource string `json:"lead_source,omitempty"`
}

type customerListResponse struct {
	Customers []Customer `json:"customers"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

//----- JOBS ---------------------------------------------------------------------------------------------------------//

type Appointment struct {
	Id string `json:"id"`
	Start time.Time `json:"start_time"`
	End time.Time `json:"end_time"`
	Window int `json:"arrival_window_minutes"`
	AssignedEmployees []string `json:"dispatched_employees_ids"`
}

type Job struct {
	Id string `json:"id"`
	CustomerId string `json:"customer_id"`
	Customer Customer
	Address Address `json:"address"`
	Note string `json:"note"`
	WorkStatus WorkStatus `json:"work_status"`
	Invoice string `json:"invoice_number"`
	Balance int64 `json:"outstanding_balance"`
	Total int64 `json:"total_amount"`
	Tags []string `json:"tags"`
	Description string `json:"description"`
	AssignedEmployees []Employee `json:"assigned_employees"`
	Schedule struct {
		Start time.Time `json:"scheduled_start"`
		End time.Time `json:"scheduled_end"`
		Window int `json:"arrival_window"`
	}
	Appointments []Appointment `json:"appointments"`
	WorkTimestamps struct {
		OnMyWay time.Time `json:"on_my_way_at"`
		Started time.Time `json:"started_at"`
		Completed time.Time `json:"completed_at"`
	} `json:"work_timestamps"`
	LeadSource string `json:"lead_source,omitempty"`
}

// returns that the job is in a state where the job is still expected to be completed in the future
func (this *Job) IsPending () bool {
	switch WorkStatus(this.WorkStatus) {
	case WorkStatus_scheduled, WorkStatus_needsScheduling:
		return true
	}
	return false // this is in a state where the job has been cancelled or already started
}

// returns that the job is in a state where everything is still a "go".  Either it hasn't happened yet, it's happening now, or it will in the future
func (this *Job) IsActive () bool {
	switch WorkStatus(this.WorkStatus) {
	case WorkStatus_scheduled, WorkStatus_inProgress, WorkStatus_completeUnrated, WorkStatus_completeRated:
		return true
	}
	return false // not an active job
}

type DispatchedEmployee struct {
	Id string `json:"employee_id"`
}

type JobSchedule struct {
	Start time.Time `json:"start_time"`
	End time.Time `json:"end_time"`
	Window int `json:"arrival_window_in_minutes"`
	Notify bool `json:"notify"`
	NotifyPro bool `json:"notify_pro"`
	DispatchedEmployees []DispatchedEmployee `json:"dispatched_employees"`
}

type JobDispatch struct {
	DispatchedEmployees []DispatchedEmployee `json:"dispatched_employees"`
}

type jobListResponse struct {
	Jobs []Job `json:"jobs"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type LineItem struct {
	Name string `json:"name"`
	Description string `json:"description"`
	UnitPrice int `json:"unit_price"`
	Quantity json.Number `json:"quantity"`
	UnitCost int `json:"unit_cost"`
	Kind string `json:"kind"`
}

type createJob struct {
	CustomerId string `json:"customer_id"`
	AddressId string `json:"address_id"`
	Schedule struct {
		Start time.Time `json:"scheduled_start"`
		End time.Time `json:"scheduled_end"`
		Window string `json:"arrival_window"`
	} `json:"schedule"`
	LineItems []LineItem `json:"line_items"`
	Employees []string `json:"assigned_employee_ids"`
	Tags []string `json:"tags"`
	LeadSource string `json:"lead_source,omitempty"`
}

//----- ESTIMATES -------------------------------------------------------------------------------------------------------//

type Estimate struct {
	Id string `json:"id"`
	EstimateNumber string `json:"estimate_number"`
	WorkStatus WorkStatus `json:"work_status"`
	// LeadSource string `json:"lead_source,omitempty"`
	Customer Customer
	Address Address `json:"address"`
	WorkTimestamps struct {
		OnMyWay time.Time `json:"on_my_way_at"`
		Started time.Time `json:"started_at"`
		Completed time.Time `json:"completed_at"`
	} `json:"work_timestamps"`
	Schedule struct {
		Start time.Time `json:"scheduled_start"`
		End time.Time `json:"scheduled_end"`
		Window int `json:"arrival_window"`
	}
	AssignedEmployees [] Employee `json:"assigned_employees"`
	Options []estimateOption
}

// returns that the job is in a state where the job is still expected to be completed in the future
func (this *Estimate) IsPending () bool {
	switch this.WorkStatus {
	case WorkStatus_scheduled, WorkStatus_needsScheduling:
		return true
	}
	return false // this is in a state where the job has been cancelled or already started
}

type estimateOption struct {
	Id string `json:"id"`
	Name string `json:"name"`
	OptionNumber string `json:"option_number"`
	TotalAmount int64 `json:"total_amount"`
	ApprovalStatus string `json:"approval_status"`
	Status WorkStatus `json:"status"`
	MessageFromPro string `json:"message_from_pro"`
	Tags []string `json:"tags"`
}

func (this *estimateOption) IsPending () bool {
	switch this.Status {
	case WorkStatus_scheduled, WorkStatus_needsScheduling:
		return true
	}
	return false // this is in a state where the job has been cancelled or already started
}

type estimateListResponse struct {
	Estimates []Estimate `json:"estimates"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type CreateEstimateOption struct {
	Name string `json:"name"`
	LineItems []LineItem `json:"line_items"`
}

type createEstimate struct {
	CustomerId string `json:"customer_id"`
	AddressId string `json:"address_id"`
	Note string `json:"note"`
	Message string `json:"message"`
	Schedule struct {
		Start time.Time `json:"start_time"`
		End time.Time `json:"end_time"`
		Window string `json:"arrival_window_in_minutes"`
		NotifyCustomer bool `json:"notify_customer"`
	} `json:"schedule"`
	Employees []string `json:"assigned_employee_ids"`
	Tags []string `json:"tags"`
	LeadSource string `json:"lead_source,omitempty"`
	Options []CreateEstimateOption `json:"options"`
}

//----- EVENTS ---------------------------------------------------------------------------------------------------------//

type recurrence struct {
	freq, count, interval, months, days int 
	dows [7]bool 
	byDayCnt, byMonthDay int
	until time.Time 
}

type Event struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
	Recurrence string `json:"recurrence_rule"`
	AssignedEmployees [] Employee `json:"assigned_employees"`
	Schedule struct {
		Start time.Time `json:"start_time"`
		End time.Time `json:"end_time"`
		TimeZone string `json:"time_zone"`
	} `json:"schedule"`
}


func (this Event) parseRecurrence ()  (*recurrence, error) {
	ret := &recurrence{
		interval: 1, // default this 
	}
	var err error 

	dows := map[string]int{"SU": 0, "MO": 1, "TU": 2, "WE": 3, "TH": 4, "FR": 5, "SA": 6}
	
	for _, tok := range strings.Split(this.Recurrence, ";") { // split it based on semicolons
		// now beak the equal sign out
		parts := strings.Split(tok, "=")
		switch parts[0] {
		case "FREQ":
			switch parts[1] {
			case "DAILY":
				ret.freq = 1
			case "WEEKLY":
				ret.freq = 7
			case "MONTHLY":
				ret.months = 1
			case "YEARLY":
				ret.months = 12
			default:
				return nil, errors.Errorf("Unknown frequency : %s", this.Recurrence)
			}

		case "INTERVAL": // this only gets set if it's not 1
			ret.interval, err = strconv.Atoi(parts[1])
			if err != nil { return nil, errors.Wrapf(err, "interval : %s", this.Recurrence) }

		case "UNTIL":
			ret.until, err = time.Parse("20060102T150405Z", parts[1])
			if err != nil { return nil, errors.Wrapf(err, this.Recurrence) }

		case "COUNT":
			ret.count, err = strconv.Atoi(parts[1])
			if err != nil { return nil, errors.Wrapf(err, "count : %s", this.Recurrence) }

		case "BYDAY":
			// record a true for each day during this week that we want to include an event for
			for _, day := range strings.Split(parts[1], ",") {
				ret.dows[dows[day]] = true 
				ret.byDayCnt++ // so we know how many days are enabled
			}

		case "BYMONTHDAY":
			ret.byMonthDay, err = strconv.Atoi(parts[1])
			if err != nil { return nil, errors.Wrapf(err, "bymonthday : %s", this.Recurrence) }
		}
	}

	ret.days = ret.freq * ret.interval // calculate out days between starts

	if ret.until.IsZero() && ret.count == 0 {
		// not set, so just go out 1 year from now
		ret.until = time.Now().AddDate(1, 0, 0)
	}

	return ret, nil 
}

// figure out the full start and end date/times for this event
func (this Event) scheduleRange () (time.Time, time.Time, error) {
	if len(this.Recurrence) == 0 { return this.Schedule.Start, this.Schedule.End, nil } // just as it is

	recur, err := this.parseRecurrence()
	if err != nil { return time.Time{}, time.Time{}, err } // bail 

	if recur.until.IsZero() {
		// use the count
		recur.count--
		return this.Schedule.Start.AddDate(0, recur.months * recur.count, recur.days * recur.count),
			this.Schedule.End.AddDate(0, recur.months * recur.count, recur.days * recur.count), nil

	} else {
		// we have a target end date, so just use that
		return this.Schedule.Start, recur.until, nil 
	}
}

// creates a list of event objects based on the recurrence schedule
// "FREQ=WEEKLY;INTERVAL=2;UNTIL=20221128T070000Z;BYDAY=SA"
func (this Event) ExtractRecurrence () (events []Event, err error) {
	events = append(events, this) // this one always gets included

	if len(this.Recurrence) == 0 { return } // no recurrence, it's just this event

	// figure it out
	recur, err := this.parseRecurrence()
	if err != nil { return nil, err }

	if recur.months == 0 && recur.days == 0 { return nil, errors.Errorf("Couldn't find a frequency : %s", this.Recurrence) }
	
	// make sure we got some expected things
	if recur.until.IsZero() && recur.count == 0 { return nil, errors.Errorf("Missing until date : %s", this.Recurrence) }

	// now start adding to our return object
	jstr, err := json.Marshal(this) // doing the copying using marshalling
	if err != nil {	return nil, errors.WithStack(err) }

	// we need to know the location and timezone, to this is so if the recurrence travels over daylight savings switches
	loc, err := time.LoadLocation(this.Schedule.TimeZone)
	if err != nil {
		fmt.Printf ("ERROR loading location :: event id: %s :: timezone: %s\n", this.Id, this.Schedule.TimeZone)
		loc, _ = time.LoadLocation("America/Chicago") // default so it works
	}

	// now figure out a time in the local timezone, this is what we keep ramping, so that this always returns
	// an event with the start/end times in UTC
	localStart := this.Schedule.Start.In(loc) // start is in utc, localStart is in whatever timezone this event is in
	localEnd := this.Schedule.End.In(loc)

	// we keep looping, we've already included the "first" event
	for {
		recur.count-- // remove one from our count, in case we're using that as our limit 

		// see if we have multipled days this week each of which need an event
		if recur.byDayCnt > 1 {
			// now loop through each day of the week that we want an event for on this count, this is usualy only 1 but could be more
			// we're going day by day here
			localStart = localStart.AddDate(0, 0, 1)
			localEnd = localEnd.AddDate(0, 0, 1)

			if recur.dows[int(localStart.Weekday())] { // this day is enabled
				nextEvent := Event{}
				err = json.Unmarshal(jstr, &nextEvent) // create our new object based on the last event
				if err != nil {	return nil, errors.WithStack(err) }

				nextEvent.Schedule.Start = localStart.In(time.UTC) // these stay in utc
				nextEvent.Schedule.End = localEnd.In(time.UTC)

				// check to see if we're done
				if recur.count <= 0 && nextEvent.Schedule.Start.After(recur.until) { return  } // we're done

				events = append(events, nextEvent)
				jstr, err = json.Marshal(nextEvent) // update this byte string for next time
				if err != nil {	return nil, errors.WithStack(err) }
			}

		} else {
			// this can just be done in a "traditional" way
		
			nextEvent := Event{}
			err = json.Unmarshal(jstr, &nextEvent) // create our new object based on the last event
			if err != nil {	return nil, errors.WithStack(err) }

			// add it to the start and end
			localStart = localStart.AddDate(0, recur.months, recur.days)
			localEnd = localEnd.AddDate(0, recur.months, recur.days)

			nextEvent.Schedule.Start = localStart.In(time.UTC) // these stay in utc
			nextEvent.Schedule.End = localEnd.In(time.UTC)

			// check to see if we're done
			if recur.count <= 0 && nextEvent.Schedule.Start.After(recur.until) { return  } // we're done

			events = append(events, nextEvent)
			jstr, err = json.Marshal(nextEvent) // update this byte string for next time
			if err != nil {	return nil, errors.WithStack(err) }
		}
	}

	return // we're done!!!
}

type eventListResponse struct {
	Events []Event `json:"events"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}


//----- PUBLIC ---------------------------------------------------------------------------------------------------------//

type HouseCall struct {
	clientId, clientSecret, callbackUrl string // for making api calls
}

// populates our oauth request with the data we have from this object
func (this *HouseCall) seedOAuth () *oauthRequest {
	if this == nil { return nil } // double check
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
