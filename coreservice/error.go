package coreservice

type errType int

// Define the err type to distinguish err
const (
	errTypeDBError errType = iota
	errTypeNotFound
)

type modelError struct {
	kind	errType
	err		error
}