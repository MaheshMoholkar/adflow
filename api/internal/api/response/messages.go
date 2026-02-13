package response

// Validation errors
const (
	ErrInvalidRequest   = "ERR_INVALID_REQUEST"
	ErrValidationFailed = "ERR_VALIDATION_FAILED"
	ErrInvalidID        = "ERR_INVALID_ID"
)

// Auth errors
const (
	ErrUnauthorized       = "ERR_UNAUTHORIZED"
	ErrInvalidCredentials = "ERR_INVALID_CREDENTIALS"
	ErrPhoneTaken         = "ERR_PHONE_TAKEN"
	ErrInvalidToken       = "ERR_INVALID_TOKEN"
	ErrExpiredToken       = "ERR_EXPIRED_TOKEN"
	ErrAuthFailed         = "ERR_AUTH_FAILED"
	ErrUserInactive       = "ERR_USER_INACTIVE"
)

// CRUD operation errors
const (
	ErrCreateFailed        = "ERR_CREATE_FAILED"
	ErrUpdateFailed        = "ERR_UPDATE_FAILED"
	ErrDeleteFailed        = "ERR_DELETE_FAILED"
	ErrGetFailed           = "ERR_GET_FAILED"
	ErrListFailed          = "ERR_LIST_FAILED"
	ErrInternalServerError = "ERR_INTERNAL_SERVER_ERROR"
)

// Not found errors
const (
	ErrNotFound = "ERR_NOT_FOUND"
)

// Template errors
const (
	ErrTemplateNotFound = "ERR_TEMPLATE_NOT_FOUND"
	ErrSMSTooLong       = "ERR_SMS_TOO_LONG"
)

// Rule errors
const (
	ErrRuleNotFound = "ERR_RULE_NOT_FOUND"
)

// Plan errors
const (
	ErrPlanRequired     = "ERR_PLAN_REQUIRED"
	ErrChannelNotInPlan = "ERR_CHANNEL_NOT_IN_PLAN"
	ErrForbidden        = "ERR_FORBIDDEN"
)

// Duplicate/Conflict errors
const (
	ErrDuplicate = "ERR_DUPLICATE"
	ErrConflict  = "ERR_CONFLICT"
)
