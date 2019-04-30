package connect

// Error is a type implementing the error interface. We use this to define
// constant errors.
type Error string

// Error implements error.
func (e Error) Error() string {
	return string(e)
}
