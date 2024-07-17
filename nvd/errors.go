package nvd

import (
	"fmt"
)

// ErrNoResults is returned when a request is made against the NVD API and the
// returned cvePage contains zero results.
type ErrNoResults struct {
	queryStr string
}

func (e ErrNoResults) Error() string {
	return fmt.Sprintf(
		"no results returned for query: %s\n",
		e.queryStr)
}
