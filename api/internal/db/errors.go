package db

import (
	"errors"
)

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrInternal = errors.New("internal error")

	ErrNoEmployer = errors.New("employer not found")

	ErrInviteTokenNotFound    = errors.New("invite token not found")
	ErrOrgUserAlreadyExists   = errors.New("org user already exists")
	ErrNoOrgUser              = errors.New("org user not found")
	ErrLastActiveAdmin        = errors.New("cannot disable last active admin")
	ErrOrgUserAlreadyDisabled = errors.New("org user already disabled")
	ErrOrgUserNotDisabled     = errors.New("org user not in disabled state")

	ErrDupCostCenterName = errors.New("duplicate cost center name")
	ErrNoCostCenter      = errors.New("cost center not found")

	ErrDupLocationName = errors.New("location name already exists")
	ErrNoLocation      = errors.New("location not found")

	ErrNoOpening = errors.New("opening not found")

	ErrNoRecruiter          = errors.New("recruiter not found")
	ErrNoStateChangeWaiting = errors.New("no state change waiting")

	ErrInvalidLocation = errors.New("one or more invalid locations specified")

	ErrInvalidRecruiter = errors.New(
		"one or more invalid recruiter emails specified",
	)
)
