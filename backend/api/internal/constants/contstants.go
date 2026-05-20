package constants

import "time"

const (
	DEFAULT_CONTEXT_TIMEOUT = 3 * time.Second

	// Number of suggested user IDs to return
	MAX_SUGGESTED_IDS = 5
	DEFAULT_PAGE_SIZE = 10

	// Claims keys
	USER_ID_CLAIM = "userID"
	MAIL_CLAIM    = "mail"
)
