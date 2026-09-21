package pkg

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RequestSaveExpense_ComputeTransportAmount(t *testing.T) {
	t.Parallel()

	validityDate := jdate.New(2026, time.December, 31)

	marshalAdditionalInfo := func(t *testing.T, v any) string {
		t.Helper()
		b, err := json.Marshal(v)
		require.NoError(t, err)
		return string(b)
	}

	tests := map[string]struct {
		input   RequestSaveExpense
		want    RequestSaveExpense
		wantErr bool
	}{
		"non-transport category": {
			input: RequestSaveExpense{Category: ExpenseCategoryNTICGear, Amount: 10000},
			want:  RequestSaveExpense{Category: ExpenseCategoryNTICGear, Amount: 10000},
		},
		"transport subscription - monthly period": {
			input: RequestSaveExpense{
				Category: ExpenseCategoryTransportSubscription,
				Amount:   9000,
				AdditionalInformation: marshalAdditionalInfo(t, AdditionalInformation{
					TransportSubscription: &TransportSubscriptionAdditionalInformation{
						SubCategory:  ExpenseSubCategoryTransportSubscriptionUrban,
						Period:       TransportSubscriptionPeriodMonthly,
						ValidityDate: validityDate,
					},
				}),
			},
			want: RequestSaveExpense{
				Category: ExpenseCategoryTransportSubscription,
				Amount:   9000,
				AdditionalInformation: marshalAdditionalInfo(t, AdditionalInformation{
					TransportSubscription: &TransportSubscriptionAdditionalInformation{
						SubCategory:    ExpenseSubCategoryTransportSubscriptionUrban,
						Period:         TransportSubscriptionPeriodMonthly,
						ValidityDate:   validityDate,
						DeclaredAmount: 9000,
					},
				}),
			},
		},
		"transport subscription - yearly period divides amount by 12": {
			input: RequestSaveExpense{
				Category: ExpenseCategoryTransportSubscription,
				Amount:   12000,
				AdditionalInformation: marshalAdditionalInfo(t, AdditionalInformation{
					TransportSubscription: &TransportSubscriptionAdditionalInformation{
						SubCategory:  ExpenseSubCategoryTransportSubscriptionUrban,
						Period:       TransportSubscriptionPeriodYearly,
						ValidityDate: validityDate,
					},
				}),
			},
			want: RequestSaveExpense{
				Category: ExpenseCategoryTransportSubscription,
				Amount:   1000,
				AdditionalInformation: marshalAdditionalInfo(t, AdditionalInformation{
					TransportSubscription: &TransportSubscriptionAdditionalInformation{
						SubCategory:    ExpenseSubCategoryTransportSubscriptionUrban,
						Period:         TransportSubscriptionPeriodYearly,
						ValidityDate:   validityDate,
						DeclaredAmount: 12000,
					},
				}),
			},
		},
		"transport subscription with invalid additional_information": {
			input:   RequestSaveExpense{Category: ExpenseCategoryTransportSubscription, AdditionalInformation: "not-json"},
			wantErr: true,
		},
		"transport subscription with nil TransportSubscription": {
			input: RequestSaveExpense{
				Category:              ExpenseCategoryTransportSubscription,
				AdditionalInformation: marshalAdditionalInfo(t, AdditionalInformation{}),
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.input
			err := got.ComputeTransportAmount()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_SalaryExpense_GetAmountToRefund(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		se   SalaryExpense
		want int64
	}{
		"round": {
			se: SalaryExpense{
				BaseAmount:        8199,
				ReimbursementRate: 50,
			},
			want: 4100,
		},
		"NTIC and >= 50000": {
			se: SalaryExpense{
				BaseAmount:      70002,
				BaseAmountWoVat: new(int64(66666)),
				Category:        ExpenseCategoryNTICGear,
			},
			want: 1945,
		},
		"NTIC and >= 50000 - last refund": {
			se: SalaryExpense{
				BaseAmount:      70002,
				BaseAmountWoVat: new(int64(66666)),
				Category:        ExpenseCategoryNTICGear,
				Refunded:        70000,
			},
			want: 2,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.se.GetAmountToRefund())
		})
	}
}
