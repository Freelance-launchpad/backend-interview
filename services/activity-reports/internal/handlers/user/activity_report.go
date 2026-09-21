package userhandlers

import (
	"net/http"

	"github.com/Freelance-launchpad/backend-interview/common/jcal"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jgin"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ActivityReportHandler struct {
	activityReportStore   domain.ActivityReportStore
	activityReportService domain.ActivityReportService
}

func NewActivityReportHandler(activityReportStore domain.ActivityReportStore, activityReportService domain.ActivityReportService) *ActivityReportHandler {
	return &ActivityReportHandler{
		activityReportStore:   activityReportStore,
		activityReportService: activityReportService,
	}
}

// @Description Get missing activity reports for the current offer.
// @Router /v3/activity-reports/missing [get]
// @Param X-Jump-Offer-ID header string true "Offer ID"
// @Produce json
// @Success 200 {array} domain.MissingReportResponse
// @Failure 403 {object} jgin.Error "INVALID_OFFER_ID"
// @Failure 500 {object} jgin.Error "INTERNAL_SERVER_ERROR"
// @Tags Activity Reports, User
func (h *ActivityReportHandler) GetMissingActivityReports(c *gin.Context) {
	ctx := c.Request.Context()

	offerID := jgin.GetOfferIDFromContext(ctx)
	if offerID == uuid.Nil {
		jgin.HandleError(c, jerror.NewPermissionDeniedError("INVALID_OFFER_ID"))
		return
	}

	periods, err := h.activityReportService.GetMissingActivityReports(ctx, offerID)
	if err != nil {
		jgin.HandleError(c, err)
		return
	}

	response := make([]domain.MissingReportResponse, 0, len(periods))
	for _, p := range periods {
		response = append(response, domain.MissingReportResponse{
			Month: int(p.Month),
			Year:  p.Year,
			Days:  jcal.Default().MaxWorkingDays(p.Year, p.Month),
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Description Create activity reports for the current offer.
// @Router /v3/activity-reports [post]
// @Param X-Jump-Offer-ID header string true "Offer ID"
// @Param body body []domain.CreateReportInput true "Activity reports"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} jgin.Error "INVALID_PAYLOAD"
// @Failure 403 {object} jgin.Error "INVALID_OFFER_ID | INVALID_MISSION_OWNER"
// @Failure 500 {object} jgin.Error "INTERNAL_SERVER_ERROR"
// @Tags Activity Reports, User
func (h *ActivityReportHandler) CreateActivityReports(c *gin.Context) {
	ctx := c.Request.Context()

	offerID := jgin.GetOfferIDFromContext(ctx)
	if offerID == uuid.Nil {
		jgin.HandleError(c, jerror.NewPermissionDeniedError("INVALID_OFFER_ID"))
		return
	}

	var body []domain.CreateReportInput
	if err := c.BindJSON(&body); err != nil {
		jgin.HandleError(c, jerror.NewValidationError("INVALID_PAYLOAD"))
		return
	}

	if len(body) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	reports := make([]domain.Report, 0, len(body))
	for _, report := range body {
		reports = append(reports, report.ToReport())
	}

	if err := h.activityReportService.CreateActivityReportsV3(ctx, offerID, reports); err != nil {
		jgin.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
