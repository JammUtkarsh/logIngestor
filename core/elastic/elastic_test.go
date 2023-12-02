package elastic

import (
	"testing"
)

func TestIsValidTimeRange(t *testing.T) {
	testCases := []struct {
		timestamp     string
		expectedFrom  string
		expectedTo    string
		expectedError bool
	}{
		{
			timestamp:     "2023-11-19T12:00:00Z 2023-11-19T13:00:00Z",
			expectedFrom:  "2023-11-19T12:00:00Z",
			expectedTo:    "2023-11-19T12:00:00Z",
			expectedError: false,
		},
		{
			timestamp:     "2023-11-19T13:00:00Z 2023-11-19T12:00:00Z",
			expectedFrom:  "",
			expectedTo:    "",
			expectedError: true,
		},
		{
			timestamp:     "2023-11-19",
			expectedFrom:  "",
			expectedTo:    "",
			expectedError: true,
		},
		{
			timestamp:     "2023-11-19T12:00:00Z 2023-11-xx-19T13:00:00Z",
			expectedFrom:  "",
			expectedTo:    "",
			expectedError: true,
		},
	}

	for _, testCase := range testCases {
		from, to, err := IsValidTimeRange(testCase.timestamp)

		if from != testCase.expectedFrom {
			t.Errorf("expected from to be %s, got %s", testCase.expectedFrom, from)
		}

		if to != testCase.expectedTo {
			t.Errorf("expected to to be %s, got %s", testCase.expectedTo, to)
		}

		if (err != nil) != testCase.expectedError {
			t.Errorf("expected error to be %t, got %t", testCase.expectedError, err != nil)
		}

	}
}
