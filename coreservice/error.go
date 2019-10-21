package coreservice

type errType int

// Define the err type to distinguish err
const (
	ErrTypeDBError errType = iota
	ErrTypeNotFound
	ErrTypeValidation
)

type ModelError struct {
	Kind errType
	Err  error
}
