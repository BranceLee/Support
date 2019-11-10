package coreservice

type errType int

// Define the err type to distinguish err
const (
	ErrTypeDBError errType = iota
	ErrTypeNotFound
	ErrTypeValidation
)

// ModelError wraps the err content and annotates what kind of the error it is
type ModelError struct {
	Kind errType
	Err  error
}
