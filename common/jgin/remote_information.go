package jgin

import (
	"github.com/gin-gonic/gin"
)

// RemoteInformation represents information about a remote user.
type RemoteInformation struct {
	RemoteUserID    *string
	RemoteUserAgent *string
	RemoteAddress   *string
}

// GetRemoteInformation retrieves remote informations from gin's context.
func GetRemoteInformation(c *gin.Context) RemoteInformation {
	return RemoteInformation{
		RemoteUserID:    new(GetUserIDFromContext(c.Request.Context())),
		RemoteUserAgent: new(c.Request.UserAgent()),
		RemoteAddress:   new(c.ClientIP()),
	}
}
