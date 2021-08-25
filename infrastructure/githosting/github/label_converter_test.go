package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelConverter_convertEntity(t *testing.T) {
	tests := map[string]struct {
		givenLabels []*LabelImpl
	}{
		"GivenActualListWithLabels_WhenConverting_ThenConvertTypes": {
			givenLabels: []*LabelImpl{
				{
					Name:        "label1",
					Description: "active label",
					Color:       "ABABAB",
					Inactive:    false,
				},
				{
					Name:        "label2",
					Description: "dead label",
					Color:       "ABABAB",
					Inactive:    true,
				},
			},
		},
		"GivenNilSlices_WhenConverting_ThenReturnEmpty": {
			givenLabels: nil,
		},
		"GivenEmptySlices_WhenConverting_ThenReturnEmpty": {
			givenLabels: []*LabelImpl{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			converter := LabelConverter{}
			asEntities := converter.ConvertToEntity(tt.givenLabels)
			assert.Len(t, asEntities, len(tt.givenLabels))
			for i := range asEntities {
				assert.Equal(t, tt.givenLabels[i], asEntities[i])
			}
			asConcrete := converter.ConvertFromEntity(asEntities)
			if tt.givenLabels == nil {
				// Keep same starting situation
				asConcrete = nil
			}
			assert.Len(t, asConcrete, len(tt.givenLabels))
			for i := range asConcrete {
				assert.Equal(t, tt.givenLabels[i], asConcrete[i])
			}
		})
	}
}
