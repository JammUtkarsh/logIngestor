package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jammutkarsh/elasticlogs/core/elastic"
	"github.com/jammutkarsh/elasticlogs/core/elastic/ingestion"
	"github.com/spf13/cobra"
)

func Serve(cmd *cobra.Command, args []string) {
	if elastic.Username == "" || elastic.Password == "" {
		panic("ELASTIC_USERNAME or ELASTIC_PASSWORD not set")
	}
	if err := elastic.Ping(); err != nil {
		fmt.Println("Elasticsearch is not running")
		os.Exit(1)
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/", ingestion.IngestionData)
	log.Fatalln(r.Run(":3000"))
}
