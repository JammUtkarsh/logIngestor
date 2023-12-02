package query

import "testing"

func TestBuildSearchQuery(t *testing.T) {
	tests := []struct {
		name         string
		inputQueries []string
		expected     string
	}{
		{
			name:         "Single query",
			inputQueries: []string{"info", "", "", "", "", "", "", "", ""},
			expected:     `{ "query": { "bool": {"must": [ {"match": {"level": "info"}} ] } } }`,
		},
		{
			name:         "All fields, except time",
			inputQueries: []string{"0", "1", "2", "3", "4", "5", "6", ""},
			expected:     `{ "query": { "bool": {"must": [ {"match": {"level": "0"}},{"match": {"message": "1"}},{"match": {"resourceId": "2"}},{"match": {"traceId": "3"}},{"match": {"spanId": "4"}},{"match": {"commit": "5"}},{"match": {"parentResourceId": "6"}} ] } } }`,
		},
		{
			name:         "time field",
			inputQueries: []string{"info", "", "", "", "", "", "", "2021-10-10T10:10:10Z 2022-11-10T10:10:10Z"},
			expected:     `{ "query": { "bool": {"must": [ {"match": {"level": "info"}},{"range": {"timestamp": {"gte": "2021-10-10T10:10:10Z", "lte": "2021-10-10T10:10:10Z"}}} ] } } }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BuildSearchQuery(test.inputQueries[0], test.inputQueries[1], test.inputQueries[2], test.inputQueries[3], test.inputQueries[4], test.inputQueries[5], test.inputQueries[6], test.inputQueries[7])
			if result != test.expected {
				t.Errorf("Test '%s' failed. Expected:\n%s\nGot:\n%s", test.name, test.expected, result)
			}
		})
	}
}
