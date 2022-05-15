package pgstat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildLikeArgsArray(t *testing.T) {
	cases := []struct {
		name           string
		input          []string
		expectedResult string
		err            error
	}{
		{
			name:           "single arg",
			input:          []string{"select"},
			expectedResult: "{%select%}",
		},
		{
			name:           "a lot of args",
			input:          []string{"select", "update", "delete", "insert"},
			expectedResult: "{%select%,%update%,%delete%,%insert%}",
		},
		{
			name:           "malformed args",
			input:          []string{"select", "", "delete", "insert"},
			expectedResult: "",
			err:            malformedArgsErr,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			args, err := buildLikeArgsArray(testCase.input)
			if testCase.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, malformedArgsErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedResult, args)
		})
	}
}
