package main

import (
	"fmt"
	"os"

	cli "github.com/jammutkarsh/elasticlogs/core/CLI"
	"github.com/jammutkarsh/elasticlogs/core/elastic"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "logs",
		Short: "logs(log search) is a CLI tool to search logs",
		Example: `logs serve
logs query [flags]`,
		ValidArgs: []string{"serve", "query"},
		Version:   "v0.1.0",
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please provide a valid subcommand")
			fmt.Println("For more information, use --help")
		},
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server at port 3000",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := elastic.Ping(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		Run: cli.Serve,
	}
	rootCmd.AddCommand(serveCmd)

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Search documents",
		Example: `
logs query --level=info --message "some message"
logs query --level=error --time "2021-10-10T10:10:10Z 2023-11-10T10:10:10Z"
logs query --level=error --message="Failed" --time="2021-10-10T10:10:10Z 2022-11-10T10:10:10Z" --all`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := elastic.Ping(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		Run: cli.Query,
	}
	queryCmd.Flags().String("level", "", "Filter by level")
	queryCmd.Flags().String("message", "", "Filter by message")
	queryCmd.Flags().String("resource", "", "Filter by resource ID")
	queryCmd.Flags().String("time", "", `Range from timestamp in RFC3339 format (ISO 8601) "{from} {to}"`)
	queryCmd.Flags().String("trace", "", "Filter by trace ID")
	queryCmd.Flags().String("span", "", "Filter by span ID")
	queryCmd.Flags().String("commit", "", "Filter by commit")
	queryCmd.Flags().String("parentresource", "", "Filter by parent resource ID")
	queryCmd.Flags().Bool("all", false, "Show all fields")
	rootCmd.AddCommand(queryCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
