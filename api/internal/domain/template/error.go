package template

import "errors"

var (
	ErrTemplateNotFound = errors.New("template not found")
	ErrSMSTooLong       = errors.New("SMS body exceeds maximum character limit")
)
