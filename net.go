// just the actual sending/receiving stuff

package housecall 

import (
    "github.com/pkg/errors"

    "fmt"
    "net/http"
    "context"
    "encoding/json"
    "io/ioutil"
    "bytes"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// handles making the request and reading the results from it 
// if there's an error the Error object will be set, otherwise it will be nil
func (this *HouseCall) finish (req *http.Request, out interface{}) (*Error, error) {
	resp, err := http.DefaultClient.Do (req)
	
	if err != nil { return nil, errors.WithStack (err) }
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll (resp.Body)

    if resp.StatusCode > 399 { 
		errObj := &Error{}
        err = errors.WithStack (json.Unmarshal (body, errObj)) // capture this error
		if err != nil { err = errors.Wrap (err, string(body)) }

        errObj.StatusCode = resp.StatusCode
        return errObj, err
    }
	
	if out != nil { err = errors.WithStack (json.Unmarshal (body, out)) }
	
	return nil, err // we're good
}


  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

func (this *HouseCall) send (ctx context.Context, requestType, link string, header map[string]string, 
						in, out interface{}) (*Error, error) {
	var jstr []byte 
	var err error 

	if in != nil {
		jstr, err = json.Marshal (in)
		if err != nil { return nil, errors.WithStack (err) }

		header["Content-Type"] = "application/json; charset=utf-8"

		fmt.Println (string(jstr))
	}
	
	req, err := http.NewRequestWithContext (ctx, requestType, fmt.Sprintf ("%s/%s", apiURL, link), bytes.NewBuffer(jstr))
	if err != nil { return nil, errors.Wrap (err, link) }

	for key, val := range header { req.Header.Set (key, val) }
	errObj, err := this.finish (req, out)
	
	return errObj, errors.Wrapf (err, " %s : %s", link, string(jstr))
}
