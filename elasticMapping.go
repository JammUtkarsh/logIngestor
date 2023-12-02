//go:build ignore

package main

import (
	"bytes"
	"net/http"
	"os"
)

const (
	URL   = "http://localhost:9200/"
	Index = "dyte-sde/"
)

var (
	username string = os.Getenv("ELASTIC_USERNAME")
	password string = os.Getenv("ELASTIC_PASSWORD")
)

func init() {
	if username == "" || password == "" {
		panic("ELASTIC_USERNAME or ELASTIC_PASSWORD not set")
	}
}

func main() {
	if err := Ping(); err != nil {
		panic(err)
	}
	if err := CreateIndexWithMapping(); err != nil {
		panic(err)
	}
	println("Index created")
}

func Ping() error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
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

func CreateIndexWithMapping() error {
	req, err := http.NewRequest(http.MethodPut, URL+Index, bytes.NewBuffer([]byte(mapping)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
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

// Placed at last because it consumes a lot of space
const mapping = `{
    "mappings": {
      "properties": {
        "level": {
          "type": "keyword"
        },
        "message": {
          "type": "text"
        },
        "resourceId": {
          "type": "keyword"
        },
        "timestamp": {
          "type": "date"
        },
        "traceId": {
          "type": "keyword"
        },
        "spanId": {
          "type": "keyword"
        },
        "commit": {
          "type": "text"
        },
        "metadata": {
          "properties": {
            "parentResourceId": {
              "type": "keyword"
            }
          }
        }
      }
    }
  }`
