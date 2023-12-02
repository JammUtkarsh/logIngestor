package ingestion

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jammutkarsh/elasticlogs/core/elastic"
)

func IngestionData(ctx *gin.Context) {
	bodyBuffer := bytes.Buffer{}
	if _, err := bodyBuffer.ReadFrom(ctx.Request.Body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBuffer.Bytes()))
	if err := ctx.BindJSON(&elastic.DataModel{}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	req, err := http.NewRequest(http.MethodPost, elastic.URL+elastic.Index+"_doc", &bodyBuffer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.SetBasicAuth(elastic.Username, elastic.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error forwarding request to elastic search cluster"})
		return
	}
	defer resp.Body.Close()
	ctx.Status(resp.StatusCode)
}
