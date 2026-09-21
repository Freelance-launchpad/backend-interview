package adminhandlers

import (
	"net/http"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jgin"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ActivityReportHandler struct {
	activityReportStore domain.ActivityReportStore
}

func NewActivityReportHandler(activityReportStore domain.ActivityReportStore) *ActivityReportHandler {
	return &ActivityReportHandler{
		activityReportStore: activityReportStore,
	}
}

// @Description Get activity reports for the given offer.
// @Router /v3/admin/offers/{offer_id}/activity-reports [get]
// @Param offer_id path string true "Offer ID"
// @Produce json
// @Success 200 {array} domain.ReportResponse
// @Failure 400 {object} jgin.Error "INVALID_OFFER_ID"
// @Failure 500 {object} jgin.Error "INTERNAL_SERVER_ERROR"
// @Tags Activity Reports, Admin
func (h *ActivityReportHandler) GetActivityReports(c *gin.Context) {
	ctx := c.Request.Context()

	offerID, err := uuid.Parse(c.Param("offer_id"))
	if err != nil {
		jgin.HandleError(c, jerror.NewValidationError("INVALID_OFFER_ID"))
		return
	}

	reports, err := h.activityReportStore.GetActivityReportsByOfferID(ctx, offerID)
	if err != nil {
		jgin.HandleError(c, err)
		return
	}

	response := make([]domain.ReportResponse, 0, len(reports))
	for _, r := range reports {
		response = append(response, domain.ReportToResponse(r))
	}

	c.JSON(http.StatusOK, response)
}
