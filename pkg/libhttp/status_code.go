package libhttp

type StatusCode int

//go:generate stringer -type=StatusCode  -linecomment
const (
	OK       StatusCode = 200 // OK
	NotFound StatusCode = 404 // Not Found
)
