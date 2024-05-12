package errors

import (
	se "github.com/usk81/ashihara/shared/domain/errors"
)

const (
	// ServiceDomain ...
	ServiceDomain = "holidays"
)

// Error Case
var (
	// CaseBadRequest ...
	CaseBadRequest = se.CaseBadRequest
	// CaseUnauthenticated ...
	CaseUnauthenticated = se.CaseUnauthenticated
	// CasePermissionDenied ...
	CasePermissionDenied = se.CasePermissionDenied
	// CaseNotFound ...
	CaseNotFound = se.CaseNotFound
	// CaseAborted ...
	CaseAborted = se.CaseAborted
	// CaseAlreadyExists ...
	CaseAlreadyExists = se.CaseAlreadyExists
	// CaseResourceExhausted ...
	CaseResourceExhausted = se.CaseResourceExhausted
	// CaseUnavailable ...
	CaseUnavailable = se.CaseUnavailable
	// CaseBackendError ...
	CaseBackendError = se.CaseBackendError
)

func init() {
	se.ServiceDomain = ServiceDomain
}

// NewCause ...
func NewCause(err error, c se.ErrCase) error {
	return se.NewCause(err, ServiceDomain, c)
}
