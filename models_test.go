
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

//----- JOBS -------------------------------------------------------------------------------------------------------//

func TestModelsJobs (t *testing.T) {
	resp := jobListResponse{}

	err := json.Unmarshal ([]byte(jobListJson), &resp)
	if err != nil { t.Fatal (err) }
}




//----- JSON TESTS -------------------------------------------------------------------------------------------------------//

const jobListJson = `{"object":"paginated_jobs","data":{"object":"list","data":[{"object":"job","id":"job_3f38bd290d0f4c3aa1a95a3253cfd69b","name":null,"created_at":"2021-05-01T14:45:48Z","updated_at":"2021-05-01T14:45:49Z","address_id":"adr_7462b044c17542089105b303ae5ae3bd","customer_id":"cus_5b4847bd860040f3943bc6ea6b020cec","invoice_number":"5","note":null,"outstanding_balance":0,"total_amount":0,"scheduled_date":"","time_zone":"America/New_York","description":"Repair - Windshield Repair","tags":{"object":"list","data":["My Website"],"url":"/jobs/job_3f38bd290d0f4c3aa1a95a3253cfd69b/tags"},"pros":{"object":"list","data":["pro_d9edbbaec96f4657a500dc3b4ef00ab3"],"url":"/jobs/job_3f38bd290d0f4c3aa1a95a3253cfd69b/pros"},"work_status_timestamps":{"start":null,"on_the_way":null,"finish":null,"travel_duration":null,"on_job_duration":null},"latest_invoice":null,"url":"/jobs/job_3f38bd290d0f4c3aa1a95a3253cfd69b","shared":false,"work_status":"needs scheduling","attachments_count":0,"primary_pro_name":"Nathan Thomas","discount":0,"job_source":"Online Booking","paid":false,"schedule":{"object":"schedule","data":{"object":"schedule","start_time":null,"end_time":null,"arrival_window_minutes":0,"recurrence_uuid":null,"recurrence_rule":null},"url":"/jobs/job_3f38bd290d0f4c3aa1a95a3253cfd69b/schedule"},"segment_count":1},{"object":"job","id":"job_87a68ea3cd514c869cbe0d0c8b30c7f1","name":null,"created_at":"2021-02-02T05:10:43Z","updated_at":"2021-02-02T05:10:44Z","address_id":"adr_17078b2751f14c74b6acb5694b09618e","customer_id":"cus_f3b190bbb90f45c3b1e65e8b7a1d31d4","invoice_number":"4","note":"","outstanding_balance":25000,"total_amount":25000,"scheduled_date":"2021-02-06T00:30:00-05:00","time_zone":"America/New_York","description":"Quickie service","tags":{"object":"list","data":[],"url":"/jobs/job_87a68ea3cd514c869cbe0d0c8b30c7f1/tags"},"pros":{"object":"list","data":["pro_d9edbbaec96f4657a500dc3b4ef00ab3"],"url":"/jobs/job_87a68ea3cd514c869cbe0d0c8b30c7f1/pros"},"work_status_timestamps":{"start":null,"on_the_way":null,"finish":null,"travel_duration":null,"on_job_duration":null},"latest_invoice":null,"url":"/jobs/job_87a68ea3cd514c869cbe0d0c8b30c7f1","shared":false,"work_status":"scheduled","attachments_count":0,"primary_pro_name":"Nathan Thomas","discount":0,"job_source":"pro","paid":false,"schedule":{"object":"schedule","data":{"object":"schedule","start_time":"2021-02-06T05:30:00Z","end_time":"2021-02-06T06:00:00Z","arrival_window_minutes":0,"recurrence_uuid":null,"recurrence_rule":null},"url":"/jobs/job_87a68ea3cd514c869cbe0d0c8b30c7f1/schedule"},"segment_count":1},{"object":"job","id":"job_f150f8d83c8746e1a913a97864fb0f9a","name":null,"created_at":"2021-02-02T05:07:40Z","updated_at":"2021-02-02T05:07:41Z","address_id":"adr_5a375eb481a24b97bee17dd8bd2fa346","customer_id":"cus_1eda014820e34e369675ed91b83aa98d","invoice_number":"3","note":"","outstanding_balance":2500,"total_amount":2500,"scheduled_date":"2021-02-06T08:00:00-05:00","time_zone":"America/New_York","description":"car detailing","tags":{"object":"list","data":[],"url":"/jobs/job_f150f8d83c8746e1a913a97864fb0f9a/tags"},"pros":{"object":"list","data":["pro_d9edbbaec96f4657a500dc3b4ef00ab3"],"url":"/jobs/job_f150f8d83c8746e1a913a97864fb0f9a/pros"},"work_status_timestamps":{"start":null,"on_the_way":null,"finish":null,"travel_duration":null,"on_job_duration":null},"latest_invoice":null,"url":"/jobs/job_f150f8d83c8746e1a913a97864fb0f9a","shared":false,"work_status":"scheduled","attachments_count":0,"primary_pro_name":"Nathan Thomas","discount":0,"job_source":"pro","paid":false,"schedule":{"object":"schedule","data":{"object":"schedule","start_time":"2021-02-06T13:00:00Z","end_time":"2021-02-06T14:00:00Z","arrival_window_minutes":0,"recurrence_uuid":null,"recurrence_rule":null},"url":"/jobs/job_f150f8d83c8746e1a913a97864fb0f9a/schedule"},"segment_count":1},{"object":"job","id":"job_eaee39ebe44947afb68875a62549fac8","name":null,"created_at":"2021-02-02T05:05:34Z","updated_at":"2021-02-02T05:05:35Z","address_id":"adr_ab5776332d0a4716af0710be3a0efd6b","customer_id":"cus_5e4e641853ac43869316e721219c30ab","invoice_number":"2","note":"","outstanding_balance":5000,"total_amount":5000,"scheduled_date":"2021-02-02T11:00:00-05:00","time_zone":"America/New_York","description":"flat tire","tags":{"object":"list","data":[],"url":"/jobs/job_eaee39ebe44947afb68875a62549fac8/tags"},"pros":{"object":"list","data":["pro_d9edbbaec96f4657a500dc3b4ef00ab3"],"url":"/jobs/job_eaee39ebe44947afb68875a62549fac8/pros"},"work_status_timestamps":{"start":null,"on_the_way":null,"finish":null,"travel_duration":null,"on_job_duration":null},"latest_invoice":null,"url":"/jobs/job_eaee39ebe44947afb68875a62549fac8","shared":false,"work_status":"scheduled","attachments_count":0,"primary_pro_name":"Nathan Thomas","discount":0,"job_source":"pro","paid":false,"schedule":{"object":"schedule","data":{"object":"schedule","start_time":"2021-02-02T16:00:00Z","end_time":"2021-02-02T16:30:00Z","arrival_window_minutes":0,"recurrence_uuid":null,"recurrence_rule":null},"url":"/jobs/job_eaee39ebe44947afb68875a62549fac8/schedule"},"segment_count":1},{"object":"job","id":"job_82fcc2c6a62745e1a70142e60720f925","name":null,"created_at":"2021-02-02T04:47:25Z","updated_at":"2021-02-02T04:47:53Z","address_id":"adr_613ab9f8200a4148a154f968f7318c39","customer_id":"cus_bd39e89d17f44117b0f8d5144b9f136e","invoice_number":"1","note":"","outstanding_balance":10000,"total_amount":10000,"scheduled_date":"2021-02-06T10:30:00-05:00","time_zone":"America/New_York","description":"Repair - Windshield Repair","tags":{"object":"list","data":[],"url":"/jobs/job_82fcc2c6a62745e1a70142e60720f925/tags"},"pros":{"object":"list","data":["pro_d9edbbaec96f4657a500dc3b4ef00ab3"],"url":"/jobs/job_82fcc2c6a62745e1a70142e60720f925/pros"},"work_status_timestamps":{"start":null,"on_the_way":null,"finish":null,"travel_duration":null,"on_job_duration":null},"latest_invoice":null,"url":"/jobs/job_82fcc2c6a62745e1a70142e60720f925","shared":false,"work_status":"scheduled","attachments_count":0,"primary_pro_name":"Nathan Thomas","discount":0,"job_source":"pro","paid":false,"schedule":{"object":"schedule","data":{"object":"schedule","start_time":"2021-02-06T15:30:00Z","end_time":"2021-02-06T16:30:00Z","arrival_window_minutes":0,"recurrence_uuid":null,"recurrence_rule":null},"url":"/jobs/job_82fcc2c6a62745e1a70142e60720f925/schedule"},"segment_count":1}],"url":"/jobs"},"page":1,"page_size":100,"total_page_count":1,"total_count":5}`
