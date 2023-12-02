package elastic

import (
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

type TimeRange struct {
	From string
	To   string
}

const (
	URL   = "http://localhost:9200/"
	Index = "dyte-sde/"
)

var (
	Username string = os.Getenv("ELASTIC_USERNAME")
	Password string = os.Getenv("ELASTIC_PASSWORD")
)

func Ping() error {
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

func CleanOutput(level, message, resourceID, traceID, spanID, commit, parentResourceID, timestamp string, m *[]DataModel) {
	for i := range *m {
		if level == "" {
			(*m)[i].Level = ""
		}
		if message == "" {
			(*m)[i].Message = ""
		}
		if resourceID == "" {
			(*m)[i].ResourceId = ""
		}
		if traceID == "" {
			(*m)[i].TraceId = ""
		}
		if spanID == "" {
			(*m)[i].SpanId = ""
		}
		if commit == "" {
			(*m)[i].Commit = ""
		}
		if parentResourceID == "" {
			(*m)[i].Metadata.ParentResourceId = ""
		}
		if timestamp == "" {
			(*m)[i].Timestamp = time.Time{}
		}
	}
}
