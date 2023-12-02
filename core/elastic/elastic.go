package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type DataModel struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ResourceId string    `json:"resourceId"`
	Timestamp  time.Time `json:"timestamp"`
	TraceId    string    `json:"traceId"`
	SpanId     string    `json:"spanId"`
	Commit     string    `json:"commit"`
	Metadata   Metadata  `json:"metadata"`
}

type Metadata struct {
	ParentResourceId string `json:"parentResourceId"`
}

type SearchResult struct {
	Took     uint64 `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits         ResultHits      `json:"hits"`
	Aggregations json.RawMessage `json:"aggregations"`
}

type ResultHits struct {
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	MaxScore float32 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Hit struct {
	Index     string              `json:"_index"`
	Type      string              `json:"_type"`
	ID        string              `json:"_id"`
	Score     float32             `json:"_score"`
	Source    json.RawMessage     `json:"_source"`
	Highlight map[string][]string `json:"highlight,omitempty"`
}


type TimeRange struct {
	From string
	To   string
}

var (
	URL      string = os.Getenv("ELASTIC_URL")
	Username string = os.Getenv("ELASTIC_USERNAME")
	Password string = os.Getenv("ELASTIC_PASSWORD")
	Index string = os.Getenv("ELASTIC_INDEX")
)

func Ping() error {
	if Username == "" || Password == "" {
		return fmt.Errorf("ELASTIC_USERNAME or ELASTIC_PASSWORD not set")
	}
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(Username, Password)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

// The function `IsValidTimeRange` checks if a given timestamp string represents a valid time range and
// returns the formatted timestamps if valid.
func IsValidTimeRange(timestamp string) (string, string, error) {
	arr := strings.Split(timestamp, " ")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("please provide a valid time range")
	}
	ts0, err := time.Parse(time.RFC3339, arr[0])
	if err != nil {
		return "", "", fmt.Errorf("{from} has invalid timestamp format")
	}
	ts1, err := time.Parse(time.RFC3339, arr[1])
	if err != nil {
		return "", "", fmt.Errorf("{to} has invalid timestamp format")
	}

	if ts0.After(ts1) {
		return "", "", fmt.Errorf("`to` must be after `from`")
	}
	return ts0.Format(time.RFC3339), ts0.Format(time.RFC3339), nil
}

// The CleanOutput function removes the fields whose flags are not set in the CLI.
func CleanOutput(flag DataModel, timestamp string, m *[]DataModel) {
	for i := range *m {
		if flag.Level == "" {
			(*m)[i].Level = ""
		}
		if flag.Message == "" {
			(*m)[i].Message = ""
		}
		if flag.ResourceId == "" {
			(*m)[i].ResourceId = ""
		}
		if flag.TraceId == "" {
			(*m)[i].TraceId = ""
		}
		if flag.SpanId == "" {
			(*m)[i].SpanId = ""
		}
		if flag.Commit == "" {
			(*m)[i].Commit = ""
		}
		if flag.Metadata.ParentResourceId == "" {
			(*m)[i].Metadata.ParentResourceId = ""
		}
		if timestamp == "" {
			(*m)[i].Timestamp = time.Time{}
		}
	}
}
