package core

import "github.com/pkg/errors"

var (
	ErrThisFeatureIsNotSupported = errors.New("this feature is not supported by this coin")
)
