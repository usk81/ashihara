package repository

import (
	"github.com/usk81/ashihara/shared/domain/repository"
)

var (
	// ErrAreadyExist is error data is found in datasource when requests to crate
	ErrAreadyExist = repository.ErrAreadyExist
	// ErrNotExist is error data is not found in datasource
	ErrNotExist = repository.ErrNotExist
	// ErrNoPermission ...
	ErrNoPermission = repository.ErrNoPermission
)
