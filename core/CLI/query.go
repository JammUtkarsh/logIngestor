package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/jammutkarsh/elasticlogs/core/elastic"
	"github.com/jammutkarsh/elasticlogs/core/elastic/query"
	"github.com/spf13/cobra"
)

func Query(cmd *cobra.Command, args []string) {
	// The minimum is set to 2 because a single field alone won't be able to much information to the user.
	if cmd.Flags().NFlag() < 2 || cmd.Flags().NFlag() > 9 {
		fmt.Println(cmd.Flags().FlagUsages())
		fmt.Printf("accepts min 2 and max 9 arg(s), received %d\n", cmd.Flags().NFlag())
		os.Exit(1)
	}
	var cmdFlags elastic.DataModel
	all, _ := cmd.Flags().GetBool("all")
	cmdFlags.Level, _ = cmd.Flags().GetString("level")
	cmdFlags.Message, _ = cmd.Flags().GetString("message")
	cmdFlags.ResourceId, _ = cmd.Flags().GetString("resource")
	cmdFlags.TraceId, _ = cmd.Flags().GetString("trace")
	cmdFlags.SpanId, _ = cmd.Flags().GetString("span")
	cmdFlags.Commit, _ = cmd.Flags().GetString("commit")
	cmdFlags.Metadata.ParentResourceId, _ = cmd.Flags().GetString("parentresource")
	timestamp, _ := cmd.Flags().GetString("time")

	m, err := query.ElasticSearch(cmdFlags, timestamp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !all {
		elastic.CleanOutput(cmdFlags, timestamp, &m)
	}
	fmt.Println(printLogs(m, all))
}

func printLogs(m []elastic.DataModel, all bool) string {
    out := "Performing query with the following options:\n\n"

    for _, v := range m {
        if v.Level != "" || all {
            out += fmt.Sprintf("Level: %s\n", v.Level)
        }
        if v.Timestamp != (time.Time{}) || all {
            out += fmt.Sprintf("Timestamp: %s\n", v.Timestamp.Format(time.RFC3339))
        }
        if v.Message != "" || all {
            out += fmt.Sprintf("Message: %s\n", v.Message)
        }
        if v.ResourceId != "" || all {
            out += fmt.Sprintf("Resource ID: %s\n", v.ResourceId)
        }
        if v.TraceId != "" || all {
            out += fmt.Sprintf("Trace ID: %s\n", v.TraceId)
        }
        if v.SpanId != "" || all {
            out += fmt.Sprintf("Span ID: %s\n", v.SpanId)
        }
        if v.Commit != "" || all {
            out += fmt.Sprintf("Commit: %s\n", v.Commit)
        }
        if v.Metadata.ParentResourceId != "" || all {
            out += fmt.Sprintf("Parent Resource ID: %s\n", v.Metadata.ParentResourceId)
        }
        out += "----------------------------------------\n"
    }

    return out
}
