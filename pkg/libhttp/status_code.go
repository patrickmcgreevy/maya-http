package libhttp

type StatusCode int

//go:generate stringer -type=StatusCode
const (
	OK       StatusCode = 201 // OK
	NotFound StatusCode = 404 // Not Found
)
