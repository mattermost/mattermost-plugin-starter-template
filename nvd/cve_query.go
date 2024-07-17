package nvd

// TODO: implement CVSS v2 and v3 severities - these are implemented but not yet
//       supported by the NVD API.

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// The NVD CVE API offers a single request path (rest/json/cves/2.0),
	// filtering/refinement of returned data is controlled by the query string.
	// Below are known query keys made available via the NVD API.
	// https://nvd.nist.gov/developers/vulnerabilities

	// QueryKeyCPEName filters results to CVEs associated with a particular CPE
	// name.
	QueryKeyCPEName = "cpeName"
	// QueryKeyCVEID returns a single CVE with a matching ID, if it exists.
	QueryKeyCVEID = "cveId"
	// QueryKeyCVETag returns CVEs holding a matching tag.
	// Tags can be one of several predefined values, see the CVETag type for
	// these.
	QueryKeyCVETag = "cveTag"
	// QueryKeyCVSSV2Severity defines a minimum threshold for CVSS V2 severity -
	// see the CVSSV2Severity type for predefined values.
	QueryKeyCVSSV2Severity = "cvssV2Severity"
	// QueryKeyHasKEV produces only results which CISA has confirmed exploitation of in
	// the wild.
	QueryKeyHasKEV = "hasKev"
	// QueryKeyIsVulnerable produces only results where 1) a CPE is associated and
	// 2) the CPE is considered vulnerable.
	QueryKeyIsVulnerable = "isVulnerable"
	// QueryKeyKeywordSearch filter results by keyword(s)
	// example (single keyword)    : "keywordSearch=Windows"
	// example (multiple keywords) : "keywordSearch=Windows Mac Linux"
	QueryKeyKeywordSearch = "keywordSearch"
	// QueryKeyPubStartDate filters results PUBLISHED only after this date.
	// NOTE: if specified, you MUST ALSO specify pubEndDate.
	QueryKeyPubStartDate = "pubStartDate"
	// QueryKeyPubEndDate filters results PUBLISHED only before this date.
	// NOTE: if specified, you MUST ALSO specify pubStartDate.
	QueryKeyPubEndDate = "pubEndDate"
	// QueryKeyResultsPerPage describes the query string key which indicates the
	// number of results returned per request.
	QueryKeyResultsPerPage = "resultsPerPage"
	// QueryKeyStartIndex describes the index at which paginated results should
	// be returned from.
	QueryKeyStartIndex = "startIndex"
)

const (
	// pathCVEsCVSSV2 is the standard URL used in all NVD CVE requests.
	pathCVEsCVSSV2 = "rest/json/cves/2.0"
	// ISO-8601 time format string (required by NVD for start/end date ranges)
	timeFormatISO8601 = "2006-01-02T15:04:05.000Z"
	// the default results per-page value
	defaultResultsPerPage = 100
)

// NewCVEQuery initializes a boilerplate CVEQuery instance returning a
// pointer to that instance.
//
// Example:
// cves, err := NewCVEQuery().ResultsPerPage(100).Keyword("Mac").Fetch()
func NewCVEQuery() *CVEQuery {
	q := &CVEQuery{
		query: url.Values{},
		u: url.URL{
			Scheme: apiRequestScheme,
			Host:   apiHostname,
			Path:   pathCVEsCVSSV2,
		},
	}

	// set a default 'resultsPerPage' value
	return q.ResultsPerPage(defaultResultsPerPage)
}

// CVEQuery serves as a query builder for retrieving CVE data from the Nist
// Vulnerability Database (NVD).
type CVEQuery struct {
	publishedWithin time.Duration
	query           url.Values
	u               url.URL
}

// String fulfills fmt.Stringer to produce just the query string component of
// the enclosed url.URL.
func (cq *CVEQuery) String() string {
	return cq.u.RawQuery
}

// Fetch executes the currently constructed CVE query, returning all results.
func (cq *CVEQuery) Fetch() ([]CVE, error) {
	// get the url
	u := cq.u

	// apply a published start/end if we have a duration
	if cq.publishedWithin != 0 {
		end := time.Now()
		start := end.Add(-cq.publishedWithin)
		cq.query.Set(QueryKeyPubStartDate, start.Format(timeFormatISO8601))
		cq.query.Set(QueryKeyPubEndDate, end.Format(timeFormatISO8601))
	}

	u.RawQuery = cq.query.Encode()

	// roll 'em up
	return cq.getCVEs(u, nil)
}

// getCVEs is a recursive function to page through all returned NVDQuery
// results.
func (cq *CVEQuery) getCVEs(u url.URL, accumulator []CVE) ([]CVE, error) {
	// init the request
	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return accumulator, err
	}

	// do the request and unmarshal the response body
	result, err := doRequestWithUnmarshal(r, cvePage{})
	if err != nil {
		return accumulator, err
	}

	// if our request has no results -> let's get outta here.
	if result.TotalResults == 0 {
		return nil, ErrNoResults{u.RawQuery}
	}

	// if accumulator is nil, initialize it now with bounds
	if accumulator == nil {
		accumulator = make([]CVE, 0, result.TotalResults)
	}

	// append our current request's returned CVEs
	accumulator = append(accumulator, result.Vulnerabilities...)

	// define the next page's starting index
	nextStartIndex := result.StartIndex + result.ResultsPerPage

	// if we have more results in the tank -> recurse
	if nextStartIndex < result.TotalResults {
		// move the URL query 'startIndex' key-value forward
		q := u.Query()
		q.Set(QueryKeyStartIndex, strconv.Itoa(nextStartIndex))
		u.RawQuery = q.Encode()

		// recurse
		return cq.getCVEs(
			u,
			accumulator)
	}

	// break recursion
	return accumulator, nil
}

// ResultsPerPage allows tuning of the number of CVEs which will be returned
// per paginated request.
func (cq *CVEQuery) ResultsPerPage(n int) *CVEQuery {
	cq.query.Set(QueryKeyResultsPerPage, strconv.Itoa(n))
	return cq
}

// CPEName allows filtering of results to CVEs associated with a particular
// CPE.
func (cq *CVEQuery) CPEName(name string) *CVEQuery {
	cq.query.Set(QueryKeyCPEName, name)
	return cq
}

// CVEID retrieves a single CVE with the provided ID.
func (cq *CVEQuery) CVEID(id string) *CVEQuery {
	cq.query.Set(QueryKeyCVEID, id)
	return cq
}

// CVETag filters results to only those associated with a given tag.
// NOTE: only a SINGLE tag can be provided via this query at a time.
func (cq *CVEQuery) CVETag(tag CVETag) *CVEQuery {
	cq.query.Set(QueryKeyCVETag, string(tag))
	return cq
}

// CVSSV2Severity filters results to only those that meet the provided MINIMUM
// CVSS V2 severity threshold.
func (cq *CVEQuery) CVSSV2Severity(severity CVSSSeverityV2) *CVEQuery {
	cq.query.Set(QueryKeyCVSSV2Severity, string(severity))
	return cq
}

// PublishedWithin filters results to only those published AFTER time.Time start
// and before time.Time end.
func (cq *CVEQuery) PublishedWithin(duration time.Duration) *CVEQuery {
	cq.publishedWithin = duration
	return cq
}

// KeywordSearch filters results to only those containing one or more keywords.
//
// NOTE: multiple keywords are provided by space-delimiting them within a SINGLE
// string.
//
// Examples:
// q.KeywordSearch("Cisco")
// q.KeywordSearch("Mac Windows")
func (cq *CVEQuery) KeywordSearch(keyword string) *CVEQuery {
	cq.query.Set(QueryKeyKeywordSearch, keyword)
	return cq
}
