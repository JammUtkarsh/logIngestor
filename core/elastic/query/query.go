package query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dyte-submissions/november-2023-hiring-JammUtkarsh/core/elastic"
	"github.com/maximelamure/elasticsearch"
)

/*
CLI Perspective:
Flags based system. If true then inclue that query.
sde query "some text" # return all the document with "some text"
sde query --level "abcd" "some text" # return all the document with level "abcd" and "some text" | whole document
sde query --level "abcd" # return all the document with level "abcd" | whole document
sde query --only --level "abcd" --commit "efas12"  # return all the document with level "abcd" | only commit and level

Flags | Text
0       |    1            -> some text
1	     |    0            -> level:abcd
*/
func ElasticSearch(level, message, resourceID, traceID, spanID, commit, parentResourceID, timestamp string) ([]elastic.DataModel, error) {
	query := BuildSearchQuery(level, message, resourceID, traceID, timestamp, spanID, commit, parentResourceID)
	var response elasticsearch.SearchResult

	req, err := http.NewRequest(http.MethodPost, elastic.URL+elastic.Index+"_search?pretty=true", strings.NewReader(query))
	if err != nil {
		return []elastic.DataModel{}, err
	}
	req.SetBasicAuth(elastic.Username, elastic.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return []elastic.DataModel{}, err
	}
	defer resp.Body.Close()

	// unmarshal resp to response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []elastic.DataModel{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []elastic.DataModel{}, err
	}

	// transform response to DataModel
	var data []elastic.DataModel
	for _, hit := range response.Hits.Hits {
		var d elastic.DataModel
		err = json.Unmarshal(hit.Source, &d)
		if err != nil {
			return []elastic.DataModel{}, err
		}
		data = append(data, d)
	}
	return data, nil
}

func BuildSearchQuery(level, message, resourceID, traceID, spanID, commit, parentResourceID, timestamp string) string {
	queryParts := make([]string, 0)
	if level != "" {
		queryParts = append(queryParts, SearchByLevelQuery(level))
	}
	if message != "" {
		queryParts = append(queryParts, SearchByMessageQuery(message))
	}
	if resourceID != "" {
		queryParts = append(queryParts, SearchByResourceIDQuery(resourceID))
	}
	if traceID != "" {
		queryParts = append(queryParts, SearchByTraceIDQuery(traceID))
	}
	if spanID != "" {
		queryParts = append(queryParts, SearchBySpanIDQuery(spanID))
	}
	if commit != "" {
		queryParts = append(queryParts, SearchByCommitQuery(commit))
	}
	if parentResourceID != "" {
		queryParts = append(queryParts, SearchByParentResourceIDQuery(parentResourceID))
	}
	if timestamp != "" {
		queryParts = append(queryParts, SearchByTimestampRangeQuery(timestamp))
	}
	return fmt.Sprintf(`{ "query": { "bool": {"must": [ %s ] } } }`, strings.Join(queryParts, ","))
}

func SearchByLevelQuery(level string) string {
	return fmt.Sprintf(`{"match": {"level": "%s"}}`, level)
}

func SearchByMessageQuery(message string) string {
	return fmt.Sprintf(`{"match": {"message": "%s"}}`, message)
}

func SearchByResourceIDQuery(resourceID string) string {
	return fmt.Sprintf(`{"match": {"resourceId": "%s"}}`, resourceID)
}

func SearchByTraceIDQuery(traceID string) string {
	return fmt.Sprintf(`{"match": {"traceId": "%s"}}`, traceID)
}

func SearchByTimestampRangeQuery(timestamp string) string {
	startDate, endDate, err := elastic.IsValidTimeRange(timestamp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return fmt.Sprintf(`{"range": {"timestamp": {"gte": "%s", "lte": "%s"}}}`, startDate, endDate)
}

func SearchBySpanIDQuery(spanID string) string {
	return fmt.Sprintf(`{"match": {"spanId": "%s"}}`, spanID)
}

func SearchByCommitQuery(commit string) string {
	return fmt.Sprintf(`{"match": {"commit": "%s"}}`, commit)
}

func SearchByParentResourceIDQuery(parentResourceID string) string {
	return fmt.Sprintf(`{"match": {"parentResourceId": "%s"}}`, parentResourceID)
}
