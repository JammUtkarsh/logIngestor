package elastic

import (
	"reflect"
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

func TestCleanOutput(t *testing.T) {
	// APIData is a sample data for testing that will be sent by elasticsearch
	APIData := []DataModel{
		{
			Level:      "error",
			Message:    "Failed to authenticate user",
			ResourceId: "server-9876",
			TraceId:    "jkl-mno-345",
			SpanId:     "span-678",
			Commit:     "9i8u7y6",
			Metadata: Metadata{
				ParentResourceId: "server-5432",
			},
		},
		{
			Level:      "error",
			Message:    "Failed to connect to DB",
			ResourceId: "server-1234",
			TraceId:    "abc-xyz-123",
			SpanId:     "span-456",
			Commit:     "5e5342f",
			Metadata: Metadata{
				ParentResourceId: "server-0987",
			},
		},
		// Add more test cases
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
			expectedResult: []DataModel{
				{
					Level:   "error",
					Message: "Failed to authenticate user",
				},
				{
					Level:   "error",
					Message: "Failed to connect to DB",
				},
			},
		},
		// Add more expected cases
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			testData := make([]DataModel, len(APIData))
			copy(testData, APIData) // Copy APIData to testData
			CleanOutput(tc.flag, tc.timestamp, &testData)
			if !reflect.DeepEqual(testData, tc.expectedResult) {
				t.Errorf("Test %s:\nExpected\n%v, got\n%v", tc.testName, tc.expectedResult, testData)
			}
		})
	}
}
