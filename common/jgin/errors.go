package jgin

import (
	"errors"
	"net/http"
	"regexp"

	jerror "github.com/Freelance-launchpad/backend-interview/common/jerror/v2"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error   string           `json:"error"`
	Details []jerror.Details `json:"details,omitempty"`
}

var errInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
var extractErrorCode = regexp.MustCompile("[A-Z][A-Z0-9_]*(?:: [a-zA-Z]+)?$")

// HandleError handles wrapped errors through Errorf, stops the chain, writes the status code and return a JSON body.
func HandleError(c *gin.Context, err error) {
	var (
		baseErr *jerror.Error
	)

	var status int
	var details []jerror.Details

	err = c.Error(err)

	switch {
	case errors.As(err, &baseErr):
		status = baseErr.HTTPStatus
		details = baseErr.Details
	default:
		status = http.StatusInternalServerError
		if !IsServiceContext(c.Request.Context()) {
			err = errInternalServerError
		}
	}

	c.AbortWithStatusJSON(status, Error{
		Error:   normalizeError(err).Error(),
		Details: details,
	})
}

func normalizeError(err error) error {
	if matches := extractErrorCode.FindStringSubmatch(err.Error()); len(matches) > 0 {
		return errors.New(matches[0])
	}

	return err
}
