package web

const (
	ErrorCodeBlockedUser    = 1
	ErrorMessageBlockedUser = "blocked user"

	ErrorCodeInvalidPassword    = 2
	ErrorMessageInvalidPassword = "invalid password"

	ErrorCodeInvalidToken    = 3
	ErrorMessageInvalidToken = "invalid token"

	ErrorCodeInternal    = 10
	ErrorMessageInternal = "internal error"

	ErrorCodeInvalidBody    = 20
	ErrorMessageInvalidBody = "invalid body"

	ErrorCodeSave    = 21
	ErrorMessageSave = "error on saving"

	ErrorCodeSendingMessage    = 23
	ErrorMessageSendingMessage = "sending message"

	ErrorCodeSearch    = 24
	ErrorMessageSearch = "error on search"
)
