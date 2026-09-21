package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UKExpenseAdditionalInformation_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		info    UKExpenseAdditionalInformation
		wantErr error
	}{
		"OK": {
			info: UKExpenseAdditionalInformation{
				ApprovedInAssignment: true,
				RelationToAssignment: "related to assignment",
			},
		},
		"empty relation to assignment": {
			info:    UKExpenseAdditionalInformation{},
			wantErr: ErrMissingRelationToAssignment,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, tt.info.Validate(), tt.wantErr)
		})
	}
}
