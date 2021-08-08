package core

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakeEventHandler implements our own fake here, since with counterfeiter we would get an import cycle.
type fakeEventHandler struct {
	err error
}

func (h fakeEventHandler) Handle(source EventSource) EventResult {
	return ToResult(source, h.err)
}

func TestFireEvent(t *testing.T) {
	gitUrl := newUrl()
	tests := map[string]struct {
		givenName      EventName
		givenSource    EventSource
		givenHandler   EventHandler
		expectedResult EventResult
	}{
		"GivenNonExistingName_ThenExpectError": {
			givenSource: EventSource{
				Url: gitUrl,
			},
			expectedResult: EventResult{
				Url:   gitUrl,
				Error: errors.New("no event handler exists for ''"),
			},
		},
		"GivenRegisteredHandler_WhenSuccessful_ThenExpectResult": {
			givenSource: EventSource{
				Url: gitUrl,
			},
			givenName:    "handler1",
			givenHandler: fakeEventHandler{},
			expectedResult: EventResult{
				Url: gitUrl,
			},
		},
		"GivenRegisteredHandler_WhenFailed_ThenExpectResultWithError": {
			givenSource: EventSource{
				Url: gitUrl,
			},
			givenName:    "handler2",
			givenHandler: fakeEventHandler{err: errors.New("failed")},
			expectedResult: EventResult{
				Url:   gitUrl,
				Error: errors.New("failed"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.givenName != "" {
				RegisterHandler(tt.givenName, tt.givenHandler)
			}
			result := <-FireEvent(tt.givenName, tt.givenSource)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func newUrl() *GitURL {
	u, _ := url.Parse("https://github.com/ccremer/greposync.git")
	return FromURL(u)
}
