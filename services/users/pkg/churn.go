package pkg

import (
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/google/uuid"
)

// ChurnType represents a churn type.
type ChurnType string

const (
	// ChurnTypeResignation represents a resignation churn type.
	ChurnTypeResignation ChurnType = "resignation"

	// ChurnTypeMutualAgreement represents a mutual agreement churn type.
	ChurnTypeMutualAgreement ChurnType = "mutual_agreement"

	// ChurnTypeTrialPeriodTerminationAtEmployerInitiative represents a trial period termination at employer initiative churn type.
	ChurnTypeTrialPeriodTerminationAtEmployerInitiative ChurnType = "trial_period_termination_at_employer_initiative"

	// ChurnTypeTrialPeriodTerminationAtEmployeeInitiative represents a trial period termination at employee initiative churn type.
	ChurnTypeTrialPeriodTerminationAtEmployeeInitiative ChurnType = "trial_period_termination_at_employee_initiative"

	// ChurnTypeDeath represents a death churn type.
	ChurnTypeDeath ChurnType = "death"

	// ChurnTypeDismissal represents a dismissal ("licenciement") churn type.
	ChurnTypeDismissal ChurnType = "dismissal"

	// ChurnTypePresumedResignation represents a presumed resignation churn type.
	ChurnTypePresumedResignation ChurnType = "presumed_resignation"
)

// IsValid checks if the churn type is valid.
func (c ChurnType) IsValid() bool {
	switch c {
	case ChurnTypeResignation,
		ChurnTypeMutualAgreement,
		ChurnTypeTrialPeriodTerminationAtEmployerInitiative,
		ChurnTypeTrialPeriodTerminationAtEmployeeInitiative,
		ChurnTypeDeath,
		ChurnTypeDismissal,
		ChurnTypePresumedResignation:
		return true
	default:
		return false
	}
}

// ChurnStatus represents a churn status.
type ChurnStatus string

const (
	// ChurnStatusNotReady represents a churn when it just has been declared by the Care team.
	ChurnStatusNotReady ChurnStatus = "not_ready"

	// ChurnStatusReady represents a churn when the 1st step of churning is completed by the Care team.
	ChurnStatusReady ChurnStatus = "ready"

	// ChurnStatusWaitingForPayment represents a churn when the churned has been validated but the final payment is yet to be done.
	ChurnStatusWaitingForPayment ChurnStatus = "waiting_for_payment"

	// ChurnStatusChurned represents a finalized churn.
	ChurnStatusChurned ChurnStatus = "churned"

	// ChurnStatusCanceled represents a canceled churn.
	ChurnStatusCanceled ChurnStatus = "canceled"
)

func (c ChurnStatus) IsValid() bool {
	return c == ChurnStatusNotReady ||
		c == ChurnStatusReady ||
		c == ChurnStatusWaitingForPayment ||
		c == ChurnStatusChurned ||
		c == ChurnStatusCanceled
}

// ChurnReason represents a churn reason.
type ChurnReason string

const (
	// ChurnReasonFoundCDICDD represents the churn reason 'found_cdi_cdd'.
	ChurnReasonFoundCDICDD ChurnReason = "found_cdi_cdd"

	// ChurnReasonStartingMyCompany represents the churn reason 'starting_my_company'.
	ChurnReasonStartingMyCompany ChurnReason = "starting_my_company"

	// ChurnReasonNoMoreMission represents the churn reason 'no_more_mission'.
	ChurnReasonNoMoreMission ChurnReason = "no_more_mission"

	// ChurnReasonTraining represents the churn reason 'training'.
	ChurnReasonTraining ChurnReason = "training"

	// ChurnReasonRetirement represents the churn reason 'retirement'.
	ChurnReasonRetirement ChurnReason = "retirement"

	// ChurnReasonMovingAbroad represents the churn reason 'moving_abroad'.
	ChurnReasonMovingAbroad ChurnReason = "moving_abroad"

	// ChurnReasonCouldNotContractualizeNewMission represents the churn reason 'could_not_contractualize_new_mission'.
	ChurnReasonCouldNotContractualizeNewMission ChurnReason = "could_not_contractualize_new_mission"

	// ChurnReasonNotSatisfiedWithJump represents the churn reason 'not_satisfied_with_jump'.
	ChurnReasonNotSatisfiedWithJump ChurnReason = "not_satisfied_with_jump"

	// ChurnReasonContractualizationProblem represents the churn reason 'contractualization_problem'.
	ChurnReasonContractualizationProblem ChurnReason = "contractualization_problem"

	// ChurnReasonMissingFeature represents the churn reason 'missing_feature'.
	ChurnReasonMissingFeature ChurnReason = "missing_feature"

	// ChurnReasonGoalReached represents the churn reason 'goal_reached'.
	ChurnReasonGoalReached ChurnReason = "goal_reached"

	// ChurnReasonOther represents the churn reason 'other'.
	ChurnReasonOther ChurnReason = "other"
)

func (c ChurnReason) IsValid() bool {
	return c == ChurnReasonFoundCDICDD ||
		c == ChurnReasonStartingMyCompany ||
		c == ChurnReasonNoMoreMission ||
		c == ChurnReasonTraining ||
		c == ChurnReasonRetirement ||
		c == ChurnReasonMovingAbroad ||
		c == ChurnReasonCouldNotContractualizeNewMission ||
		c == ChurnReasonNotSatisfiedWithJump ||
		c == ChurnReasonContractualizationProblem ||
		c == ChurnReasonMissingFeature ||
		c == ChurnReasonGoalReached ||
		c == ChurnReasonOther
}

// Churn represents a churn.
type Churn struct {
	ID         uuid.UUID      `json:"id"`
	UserID     string         `json:"user_id"`
	ContractID uuid.UUID      `json:"contract_id"`
	OfferID    uuid.UUID      `json:"offer_id"`
	UserType   jentity.Entity `json:"user_type"`
	FullName   string         `json:"full_name"`
	ChurnType  ChurnType      `json:"churn_type"`
	Reason     ChurnReason    `json:"reason"`
	Status     ChurnStatus    `json:"status"`
	ChurnDate  jdate.Date     `json:"churn_date"`
	OwnerEmail *string        `json:"owner_email"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at"`
}

func (c Churn) IsValid(today jdate.Date) bool {
	return c.ChurnType.IsValid() && c.Reason.IsValid() && !c.ChurnDate.Before(today)
}
