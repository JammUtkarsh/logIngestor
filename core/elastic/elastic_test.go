package elastic

import (
	"reflect"
	"testing"
	"time"
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

func TestCleanOutput(t *testing.T) {
	apiData := []DataModel{
		{
			Level:      "error",
			Message:    "failed",
			ResourceId: "123",
			Timestamp:  time.Now(),
			TraceId:    "trace123",
			SpanId:     "span456",
			Commit:     "commit789",
			Metadata: Metadata{
				ParentResourceId: "parent123",
			},
		},
		// Add more test cases for other scenarios
	}

	// Define test cases
	testCases := []struct {
		testName       string
		flag           DataModel
		timestamp      string
		expectedResult []DataModel
	}{
		{
			testName: "Two flags are set",
			flag: DataModel{
				Level:   "error",
				Message: "failed",
			},
			timestamp: "",
			expectedResult: []DataModel{
				{
					Level:      "error",
					Message:    "failed",
					ResourceId: "",
					Timestamp:  time.Time{},
					TraceId:    "",
					SpanId:     "",
					Commit:     "",
					Metadata:   Metadata{},
				},
			},
		},
		// Add more test cases for other scenarios
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			CleanOutput(tc.flag, tc.timestamp, &apiData)
			if !reflect.DeepEqual(apiData, tc.expectedResult) {
				t.Errorf("Test %s:\nExpected\n%v, got\n%v", tc.testName, tc.expectedResult, apiData)
			}
		})
	}
}
