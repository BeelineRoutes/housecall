
package housecall 

import (
	
	"github.com/stretchr/testify/assert"
	//"github.com/pkg/errors"

	"testing"
	"time"
	"encoding/json"
)

// used for testing so I don't have my actual credentials in the repo
type testConfig struct {
	ClientId, ClientSecret, RedirectUrl, OAuthCode, Token string 
}

func TestFirstModelsError (t *testing.T) {
	var err *Error 

	assert.Equal (t, nil, err.Err(), "nil for a nil object")

	// give it some memory
	err = &Error{}
	assert.NotEqual (t, nil, err.Err(), "Should return an error")
	
}

func TestFirstModelsOAuthResponse (t *testing.T) {
	resp := oauthResponse {
		Expires: 1000,
		Created: 1619557886,
	}

	tm, _ := time.Parse ("2006-01-02 15:04:05", "2021-04-27 21:28:06")
	assert.Equal (t, tm.Unix(), resp.ExpiresAt().Unix(), "expiring time")
}


func TestFirstModelsNewHouseCall (t *testing.T) {
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

//----- JOBS -------------------------------------------------------------------------------------------------------//

func TestFirstModelsJobs (t *testing.T) {
	resp := jobListResponse{}

	err := json.Unmarshal ([]byte(jobListJson), &resp)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 6, resp.TotalItems, "total items")
	assert.Equal (t, 1, resp.TotalPages, "total pages")
	assert.Equal (t, 6, len(resp.Jobs), "total jobs")
}




//----- JSON TESTS -------------------------------------------------------------------------------------------------------//

const jobListJson = `{"page":1,"page_size":100,"total_pages":1,"total_items":6,"jobs":[{"id":"job_94dec270539c4566be8b11173323ef5f","invoice_number":"6","description":"Install \u0026 Remove - Oil Change","customer":{"id":"cus_cf1aa167d1b54f88ba00ea95700ef3fc","first_name":"Adam","last_name":"Thomas","email":"aawaite@gmail.com","mobile_number":"5083084921","home_number":null,"work_number":null,"company":null,"notifications_enabled":true,"tags":[]},"address":{"id":"adr_17baddfd19fc4eabac09c670abffc793","type":"billing","street":"40 Riverview Drive","street_line_2":null,"city":"Burlington","state":"VT","zip":"05408","country":null},"note":null,"work_status":"needs scheduling","work_timestamps":{"on_my_way_at":null,"started_at":null,"completed_at":null},"schedule":{"scheduled_start":null,"scheduled_end":null,"arrival_window":0},"total_amount":0,"outstanding_balance":0,"assigned_employees":[{"id":"pro_d9edbbaec96f4657a500dc3b4ef00ab3","first_name":"Nathan","last_name":"Thomas","email":"nathan@necustomsoftware.com","mobile_number":"6175433004","color_hex":"EF9159","avatar_url":"/assets/add_image_thumb.png","role":"field tech"}],"tags":["My Website"]},{"id":"job_3f38bd290d0f4c3aa1a95a3253cfd69b","invoice_number":"5","description":"Repair - Windshield Repair","customer":{"id":"cus_5b4847bd860040f3943bc6ea6b020cec","first_name":"Alissa","last_name":"Thomas","email":"atortfeasor@gmail.com","mobile_number":"6175433004","home_number":null,"work_number":null,"company":null,"notifications_enabled":true,"tags":[]},"address":{"id":"adr_7462b044c17542089105b303ae5ae3bd","type":"billing","street":"23 POTTER PL","street_line_2":null,"city":"Shelburne","state":"VT","zip":"05482","country":null},"note":null,"work_status":"needs scheduling","work_timestamps":{"on_my_way_at":null,"started_at":null,"completed_at":null},"schedule":{"scheduled_start":null,"scheduled_end":null,"arrival_window":0},"total_amount":0,"outstanding_balance":0,"assigned_employees":[{"id":"pro_d9edbbaec96f4657a500dc3b4ef00ab3","first_name":"Nathan","last_name":"Thomas","email":"nathan@necustomsoftware.com","mobile_number":"6175433004","color_hex":"EF9159","avatar_url":"/assets/add_image_thumb.png","role":"field tech"}],"tags":["My Website"]},{"id":"job_87a68ea3cd514c869cbe0d0c8b30c7f1","invoice_number":"4","description":"Quickie service","customer":{"id":"cus_f3b190bbb90f45c3b1e65e8b7a1d31d4","first_name":"Elaina","last_name":"Waite","email":null,"mobile_number":null,"home_number":null,"work_number":null,"company":null,"notifications_enabled":false,"tags":[]},"address":{"id":"adr_17078b2751f14c74b6acb5694b09618e","type":"billing","street":"301 College St","street_line_2":null,"city":"Burlington","state":"VT","zip":"05401","country":"US"},"note":"","work_status":"scheduled","work_timestamps":{"on_my_way_at":null,"started_at":null,"completed_at":null},"schedule":{"scheduled_start":"2021-02-06T05:30:00Z","scheduled_end":"2021-02-06T06:00:00Z","arrival_window":0},"total_amount":25000,"outstanding_balance":25000,"assigned_employees":[{"id":"pro_d9edbbaec96f4657a500dc3b4ef00ab3","first_name":"Nathan","last_name":"Thomas","email":"nathan@necustomsoftware.com","mobile_number":"6175433004","color_hex":"EF9159","avatar_url":"/assets/add_image_thumb.png","role":"field tech"}],"tags":[]},{"id":"job_f150f8d83c8746e1a913a97864fb0f9a","invoice_number":"3","description":"car detailing","customer":{"id":"cus_1eda014820e34e369675ed91b83aa98d","first_name":"Karen","last_name":"Waite","email":null,"mobile_number":null,"home_number":null,"work_number":null,"company":null,"notifications_enabled":false,"tags":[]},"address":{"id":"adr_5a375eb481a24b97bee17dd8bd2fa346","type":"billing","street":"40 Airport Rd","street_line_2":null,"city":"South Burlington","state":"VT","zip":"05403","country":"US"},"note":"","work_status":"scheduled","work_timestamps":{"on_my_way_at":null,"started_at":null,"completed_at":null},"schedule":{"scheduled_start":"2021-02-06T13:00:00Z","scheduled_end":"2021-02-06T14:00:00Z","arrival_window":0},"total_amount":2500,"outstanding_balance":2500,"assigned_employees":[{"id":"pro_d9edbbaec96f4657a500dc3b4ef00ab3","first_name":"Nathan","last_name":"Thomas","email":"nathan@necustomsoftware.com","mobile_number":"6175433004","color_hex":"EF9159","avatar_url":"/assets/add_image_thumb.png","role":"field tech"}],"tags":[]},{"id":"job_eaee39ebe44947afb68875a62549fac8","invoice_number":"2","description":"flat tire","customer":{"id":"cus_5e4e641853ac43869316e721219c30ab","first_name":"Gene","last_name":"Thomas","email":"jo@jothomas.art","mobile_number":"4235963474","home_number":null,"work_number":null,"company":null,"notifications_enabled":false,"tags":[]},"address":{"id":"adr_ab5776332d0a4716af0710be3a0efd6b","type":"billing","street":"4658 Mt Philo Rd","street_line_2":null,"city":"Charlotte","state":"VT","zip":"05445","country":"US"},"note":"","work_status":"scheduled","work_timestamps":{"on_my_way_at":null,"started_at":null,"completed_at":null},"schedule":{"scheduled_start":"2021-02-02T16:00:00Z","scheduled_end":"2021-02-02T16:30:00Z","arrival_window":0},"total_amount":5000,"outstanding_balance":5000,"assigned_employees":[{"id":"pro_d9edbbaec96f4657a500dc3b4ef00ab3","first_name":"Nathan","last_name":"Thomas","email":"nathan@necustomsoftware.com","mobile_number":"6175433004","color_hex":"EF9159","avatar_url":"/assets/add_image_thumb.png","role":"field tech"}],"tags":[]},{"id":"job_82fcc2c6a62745e1a70142e60720f925","invoice_number":"1","description":"Repair - Windshield Repair","customer":{"id":"cus_bd39e89d17f44117b0f8d5144b9f136e","first_name":"Pat","last_name":"Waite","email":null,"mobile_number":"2075776012","home_number":null,"work_number":null,"company":null,"notifications_enabled":false,"tags":[]},"address":{"id":"adr_613ab9f8200a4148a154f968f7318c39","type":"service","street":"1436 Williston Rd","street_line_2":null,"city":"South Burlington","state":"VT","zip":"05403","country":"US"},"note":"","work_status":"scheduled","work_timestamps":{"on_my_way_at":null,"started_at":null,"completed_at":null},"schedule":{"scheduled_start":"2021-02-06T15:30:00Z","scheduled_end":"2021-02-06T16:30:00Z","arrival_window":0},"total_amount":10000,"outstanding_balance":10000,"assigned_employees":[{"id":"pro_d9edbbaec96f4657a500dc3b4ef00ab3","first_name":"Nathan","last_name":"Thomas","email":"nathan@necustomsoftware.com","mobile_number":"6175433004","color_hex":"EF9159","avatar_url":"/assets/add_image_thumb.png","role":"field tech"}],"tags":[]}]}`
