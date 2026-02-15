package template

import "errors"

var (
	ErrTemplateNotFound = errors.New("template not found")
	ErrSMSTooLong       = errors.New("SMS body exceeds maximum character limit")
	ErrInvalidImageURL  = errors.New("image_url must be a valid https URL")
	ErrMissingImageKey  = errors.New("image_key is required when image_url is set")
	ErrUploadDisabled   = errors.New("image upload is not configured")
)
