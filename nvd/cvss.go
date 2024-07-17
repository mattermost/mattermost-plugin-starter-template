package nvd

import (
	"bytes"
	"fmt"
)

// CVSSMetricV2 describes metadata associated with a CVE's CVSS V2 scoring.
type CVSSMetricV2 struct {
	Source              string         `json:"source"`
	Type                string         `json:"type"`
	CVSSData            CVSSData       `json:"cvssData"`
	BaseSeverity        CVSSSeverityV2 `json:"baseSeverity"`
	ExploitabilityScore float32        `json:"exploitabilityScore"`
	ImpactScore         float32        `json:"impactScore"`
}

// CVSSSeverityV2 describes one of a set of pre-defined values indicating
// thresholds of vulnerability signficance.
type CVSSSeverityV2 string

const (
	SeverityHigh CVSSSeverityV2 = "HIGH"
	SeverityMed  CVSSSeverityV2 = "MEDIUM"
	SeverityLow  CVSSSeverityV2 = "LOW"
)

var (
	bArrSpace = []byte(" ")
	bArrQuote = []byte("\"")
	bArrEmpty = []byte("")
)

// UnmarshalJSON fulfills the json.Unmarshaler interface to ensure the received
// value is parsed as a true CVSSSeverityV2 type and is valid.
func (s *CVSSSeverityV2) UnmarshalJSON(b []byte) (err error) {
	// throw away unnecessary characters
	b = bytes.ReplaceAll(b, bArrQuote, bArrEmpty)
	b = bytes.ReplaceAll(b, bArrSpace, bArrEmpty)

	// severity cast the byte slice
	*s = CVSSSeverityV2(b)

	// confirm the received value is inline with expected values
	switch *s {
	case SeverityHigh, SeverityMed, SeverityLow:
		return nil
	default:
		return fmt.Errorf("unexpected baseSeverity value '%s'", *s)
	}
}

// CVSSData is primarily implemented currently to reach the BaseScore field
// which contains a 0.0-10.0 value indicating a general vulnerability
// significance.
type CVSSData struct {
	Version   string  `json:"version"`
	BaseScore float32 `json:"baseScore"`
}
