package nvd

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	// standard URI scheme for all requests.
	apiRequestScheme = "https"
	// the standard NVD API hostname for all requests.
	apiHostname = "services.nvd.nist.gov"
)

// standard http client used by all nvd package requests.
var client = http.Client{Timeout: 5 * time.Second}

// doRequestWithUnmarshal performs an HTTP request and attempts an unmarshal of
// the response body to provided type T, returning the resulting instance of T
// along with any error encountered.
func doRequestWithUnmarshal[T any](
	r *http.Request,
	responseType T,
) (T, error) {
	// execute the request
	response, err := client.Do(r)
	if err != nil {
		return responseType, err
	}
	defer response.Body.Close()

	// attempt response decode
	err = json.NewDecoder(response.Body).Decode(&responseType)

	return responseType, err
}
