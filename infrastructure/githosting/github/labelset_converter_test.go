package github

import (
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"
)

func TestLabelSetConverter_ConvertToEntity(t *testing.T) {
	tests := map[string]struct {
		givenLabels []*github.Label
	}{
		"GivenNilSlices_WhenConverting_ThenReturnEmpty": {
			givenLabels: nil,
		},
		"GivenEmptySlices_WhenConverting_ThenReturnEmpty": {
			givenLabels: []*github.Label{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			converter := LabelSetConverter{}
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
