package generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/clients/gotenberg"
	"github.com/Freelance-launchpad/backend-interview/common/jcal"
	"github.com/Freelance-launchpad/backend-interview/common/jclock"
	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/js3/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jtemplate"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	usersclient "github.com/Freelance-launchpad/backend-interview/services/users/pkg/client"
)

type PDFGenerator struct {
	template    *template.Template
	usersClient usersclient.Client
	converter   gotenberg.Client
	clock       jclock.Clock
}

func NewPDFGenerator(coopActivityReportTemplateHTML []byte, usersClient usersclient.Client, converter gotenberg.Client) (*PDFGenerator, error) {
	template, err := template.New("coop_activity_report").Parse(string(coopActivityReportTemplateHTML))
	if err != nil {
		return nil, fmt.Errorf("template.New.ParseFiles has failed: %w", err)
	}

	return &PDFGenerator{
		template:    template,
		converter:   converter,
		usersClient: usersClient,
		clock:       jclock.RealClock{},
	}, nil
}

type templateUnit string

const (
	templateUnitHour  templateUnit = "heure"
	templateUnitHours templateUnit = "heures"
	templateUnitDay   templateUnit = "jour"
	templateUnitDays  templateUnit = "jours"
)

func getTemplateUnit(duration float64, unit domain.Unit) templateUnit {
	if unit == domain.UnitDays {
		if duration < 2 {
			return templateUnitDay
		}
		return templateUnitDays
	}
	if duration < 2 {
		return templateUnitHour
	}
	return templateUnitHours
}

type templateData struct {
	Logo       template.HTML
	Signature  template.HTML
	Stamp      template.HTML
	Stylesheet template.CSS
	Entity     jentity.Entity

	Female   bool
	FullName string

	ReportMonthPreposition string
	ReportMonth            string
	CurrentDate            string

	Work                  string
	WorkUnit              templateUnit
	Prospection           string
	ProspectionUnit       templateUnit
	Formation             string
	FormationUnit         templateUnit
	Vacation              string
	VacationUnit          templateUnit
	Away                  string
	AwayUnit              templateUnit
	WeeklyWorkedHours     string
	WeeklyWorkedHoursUnit templateUnit

	OtherActivity          *string
	OtherActivityUnit      *templateUnit
	OtherActivityFrequency *string

	WeekType     string
	ExtraRestDay *string
}

func (p *PDFGenerator) GenerateActivityReport(ctx context.Context, report domain.Report) (_ js3.FileOutput, err error) {
	defer jerror.Wrap(&err)

	const fileName = "rapport_activité_%s_%d.pdf"

	if !report.ActivityType.IsLife() {
		return js3.FileOutput{}, fmt.Errorf("invalid activity type: %s", report.ActivityType)
	}

	data := templateData{
		Logo:       template.HTML(jtemplate.LogoJump),
		Signature:  template.HTML(jtemplate.SignatureNicolasFayon),
		Stylesheet: template.CSS(jtemplate.Stylesheet),
		Entity:     report.ActivityType,
	}

	switch report.ActivityType {
	case jentity.EntityBlue:
		data.Stamp = template.HTML(jtemplate.StampJumpBlue)
	case jentity.EntityGreen:
		data.Stamp = template.HTML(jtemplate.StampJumpGreen)
	case jentity.EntityRealty:
		data.Stamp = template.HTML(jtemplate.StampJumpRealty)
	case jentity.EntityFormation:
		data.Stamp = template.HTML(jtemplate.StampJumpFormation)
	}

	user, err := p.usersClient.GetUserByOfferID(ctx, report.OfferID)
	if err != nil {
		return js3.FileOutput{}, err
	}

	data.FullName = fmt.Sprintf("%s %s", user.LastName, user.FirstName)

	if user.SocialSecurityNumber == nil {
		return js3.FileOutput{}, errors.New("user has no social security number")
	}

	if (*user.SocialSecurityNumber)[0:1] == "2" {
		data.Female = true
	}

	data.ReportMonthPreposition = "de "
	monthString := jcal.MonthToString(time.Month(report.Month), "FRA")
	if monthString[0:1] == "A" || monthString[0:1] == "O" {
		data.ReportMonthPreposition = "d'"
	}
	data.ReportMonth = fmt.Sprintf("%s %d", monthString, report.Year)
	data.CurrentDate = time.Now().Format("02/01/2006")

	data.Work = fmt.Sprintf("%.01f", report.WorkDuration)
	data.WorkUnit = getTemplateUnit(report.WorkDuration, report.Unit)
	data.Prospection = fmt.Sprintf("%.01f", report.ProspectionDuration)
	data.ProspectionUnit = getTemplateUnit(report.ProspectionDuration, report.Unit)
	data.Formation = fmt.Sprintf("%.01f", report.FormationDuration)
	data.FormationUnit = getTemplateUnit(report.FormationDuration, report.Unit)
	data.Vacation = fmt.Sprintf("%.0f", report.VacationDays)
	data.VacationUnit = getTemplateUnit(report.VacationDays, domain.UnitDays)
	data.Away = fmt.Sprintf("%.0f", report.DaysAway)
	data.AwayUnit = getTemplateUnit(report.DaysAway, domain.UnitDays)

	if report.ActivityType != jentity.EntityBlue {
		weeklyWorkedHours := (report.WorkDuration * 12) / 52
		data.WeeklyWorkedHours = fmt.Sprintf("%.01f", weeklyWorkedHours)
		data.WeeklyWorkedHoursUnit = getTemplateUnit(weeklyWorkedHours, domain.UnitHours)
	}

	if report.OtherActivity != nil && report.OtherActivityFrequency != nil {
		data.OtherActivity = new(fmt.Sprintf("%.01f", *report.OtherActivity))
		data.OtherActivityUnit = new(getTemplateUnit(*report.OtherActivity, domain.UnitHours))
		if *report.OtherActivityFrequency == "daily" {
			data.OtherActivityFrequency = new("jour")
		} else {
			data.OtherActivityFrequency = new("semaine")
		}
	}

	if report.WeekType == nil {
		return js3.FileOutput{}, errors.New("week type is nil")
	}
	if *report.WeekType == pkg.WeekTypeMonFri {
		data.WeekType = "lundi au vendredi"
	} else {
		data.WeekType = "mardi au samedi"
	}
	if report.FirstDayAsExtraRest == nil {
		return js3.FileOutput{}, errors.New("first day as extra rest is nil")
	}
	if *report.FirstDayAsExtraRest {
		data.ExtraRestDay = new(jdate.New(report.Year, time.Month(report.Month), 1).Format(jdate.FrenchDateFormat))
	}

	length, pdf, err := p.generateTemplatedPDF(ctx, data)
	if err != nil {
		return js3.FileOutput{}, err
	}

	fileOutput := js3.FileOutput{
		ContentLength: length,
		ContentType:   "application/pdf",
		Reader:        io.NopCloser(pdf),
		FileName:      fmt.Sprintf(fileName, monthString, report.Year),
	}

	return fileOutput, nil
}

func (p *PDFGenerator) generateTemplatedPDF(ctx context.Context, data templateData) (int64, io.ReadSeeker, error) {
	const errorMsg = "PDFGenerator.generateTemplatedPDF has failed"

	var buffer bytes.Buffer
	if err := p.template.Execute(&buffer, data); err != nil {
		return 0, nil, fmt.Errorf("%s: %w", errorMsg, err)
	}

	length, pdf, err := p.converter.HTMLToPDF(ctx, &buffer, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("%s: %w", errorMsg, err)
	}

	return length, pdf, nil
}
