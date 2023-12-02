package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/dyte-submissions/november-2023-hiring-JammUtkarsh/core/elastic"
	"github.com/dyte-submissions/november-2023-hiring-JammUtkarsh/core/elastic/query"
	"github.com/spf13/cobra"
)

func Query(cmd *cobra.Command, args []string) {
	// The minimum is set to 2 because a single field alone won't be able to much information to the user.
	if cmd.Flags().NFlag() < 2 || cmd.Flags().NFlag() > 9 {
		fmt.Println(cmd.Flags().FlagUsages())
		fmt.Printf("accepts min 2 and max 9 arg(s), received %d\n", cmd.Flags().NFlag())
		os.Exit(1)
	}

	all, _ := cmd.Flags().GetBool("all")
	level, _ := cmd.Flags().GetString("level")
	message, _ := cmd.Flags().GetString("message")
	resourceID, _ := cmd.Flags().GetString("resource")
	traceID, _ := cmd.Flags().GetString("trace")
	timestamp, _ := cmd.Flags().GetString("time")
	spanID, _ := cmd.Flags().GetString("span")
	commit, _ := cmd.Flags().GetString("commit")
	parentResourceID, _ := cmd.Flags().GetString("parentresource")

	m, err := query.ElasticSearch(level, message, resourceID, traceID, spanID, commit, parentResourceID, timestamp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !all {
		elastic.CleanOutput(level, message, resourceID, traceID, spanID, commit, parentResourceID, timestamp, &m)
	}
	fmt.Println(printLogs(m, all))
}

func printLogs(m []elastic.DataModel, all bool) string {
	out := "Performing query with the following options:\n\n"
	if all {
		for _, v := range m {
			out += fmt.Sprintf("Level: %s\n", v.Level)
			out += fmt.Sprintf("Timestamp: %s\n", v.Timestamp.Format(time.RFC3339))
			out += fmt.Sprintf("Message: %s\n", v.Message)
			out += fmt.Sprintf("Resource ID: %s\n", v.ResourceId)
			out += fmt.Sprintf("Trace ID: %s\n", v.ResourceId)
			out += fmt.Sprintf("Span ID: %s\n", v.SpanId)
			out += fmt.Sprintf("Commit: %s\n", v.Commit)
			out += fmt.Sprintf("Parent Resource ID: %s\n", v.Metadata.ParentResourceId)
			out += "----------------------------------------\n"
		}
	} else {
		for _, v := range m {
			if v.Level != "" {
				out += fmt.Sprintf("Level: %s\n", v.Level)
			}
			if v.Timestamp != (time.Time{}) {
				out += fmt.Sprintf("Timestamp: %s\n", v.Timestamp.Format(time.RFC3339))
			}
			if v.Message != "" {
				out += fmt.Sprintf("Message: %s\n", v.Message)
			}
			if v.ResourceId != "" {
				out += fmt.Sprintf("Resource ID: %s\n", v.ResourceId)
			}
			if v.TraceId != "" {
				out += fmt.Sprintf("Trace ID: %s\n", v.ResourceId)
			}
			if v.SpanId != "" {
				out += fmt.Sprintf("Span ID: %s\n", v.SpanId)
			}
			if v.Commit != "" {
				out += fmt.Sprintf("Commit: %s\n", v.Commit)
			}
			if v.Metadata.ParentResourceId != "" {
				out += fmt.Sprintf("Parent Resource ID: %s\n", v.Metadata.ParentResourceId)
			}
			out += "----------------------------------------\n"
		}
	}
	return out
}
