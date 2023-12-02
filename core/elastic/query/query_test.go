package query

import (
	"testing"

	"github.com/jammutkarsh/elasticlogs/core/elastic"
)

func TestBuildSearchQuery(t *testing.T) {
	tests := []struct {
		name         string
		inputQueries elastic.DataModel
		timestamp    string
		expected     string
	}{
		{
			name:         "Single query",
			inputQueries: elastic.DataModel{Level: "info"},
			expected:     `{ "query": { "bool": {"must": [ {"match": {"level": "info"}} ] } } }`,
		},
		{
			name:         "All fields, except time",
			inputQueries: elastic.DataModel{Level: "0", Message: "1", ResourceId: "2", TraceId: "3", SpanId: "4", Commit: "5", Metadata: elastic.Metadata{ParentResourceId: "6"}},
			expected:     `{ "query": { "bool": {"must": [ {"match": {"level": "0"}},{"match": {"message": "1"}},{"match": {"resourceId": "2"}},{"match": {"traceId": "3"}},{"match": {"spanId": "4"}},{"match": {"commit": "5"}},{"match": {"parentResourceId": "6"}} ] } } }`,
		},
		{
			name:         "time field",
			inputQueries: elastic.DataModel{Level: "info"},
			timestamp:    "2021-10-10T10:10:10Z 2022-11-10T10:10:10Z",
			expected:     `{ "query": { "bool": {"must": [ {"match": {"level": "info"}},{"range": {"timestamp": {"gte": "2021-10-10T10:10:10Z", "lte": "2021-10-10T10:10:10Z"}}} ] } } }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BuildSearchQuery(test.inputQueries, test.timestamp)
			if result != test.expected {
				t.Errorf("Test '%s' failed. Expected:\n%s\nGot:\n%s", test.name, test.expected, result)
			}
		})
	}
}
