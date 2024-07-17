package nvd

// TODO: implement CVETags support in the CVE struct. CVE tags is a polymorphic
//       JSON field and so was skipped for the initial implementation for the
//       sake of time.

// CVE describes the shape of a CVE result returned by the API and contained in
// the top-level 'vulnerabilities' array.
type CVE struct {
	Key struct {
		ID               string           `json:"id"`
		SourceIdentifier string           `json:"sourceIdentifier"`
		Published        Timestamp        `json:"published"`
		LastModified     Timestamp        `json:"lastModified"`
		Descriptions     []CVEDescription `json:"descriptions"`
		Metrics          CVEMetrics       `json:"metrics"`
		Weaknesses       []CVEWeakness    `json:"weaknesses"`
		References       []CVEReference   `json:"references"`
	} `json:"cve"`
}

// CVEDescription contains usually a single-line/sentence briefly describing
// the vulnerability.
type CVEDescription struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// CVEMetrics holds data about the CVE's CVSS v2 qualifications.
//
// NOTE: this will likely be where CVSS v3 and CVSS v4 are supported in the
// future but are not yet.
type CVEMetrics struct {
	CVSSMetricV2 []CVSSMetricV2 `json:"cvssMetricV2"`
}

// CVEWeakness describes the Common Weakness Enumeration Specification (CWE)
// associated with this CVE.
type CVEWeakness struct {
	Source      string                   `json:"source"`
	Type        string                   `json:"type"`
	Description []CVEWeaknessDescription `json:"description"`
}

// CVEWeaknessDescription 's 'Value' field holds the actual CWE specifier
// associated with a CVE.
// https://nvd.nist.gov/vuln/categories
type CVEWeaknessDescription struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// CVEReference holds a URL and source associated with a CVE. These URLs can
// be thinks like the actual announcement, mitigation steps and more. The
// Source field provides a brief description of from where the link was
// submitted (ex: US Government Resource, MITRE).
type CVEReference struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}
