package main

import (
	"context"
	"flag"
	"log"
	"os"

	exporter "github.com/yudppp/notion-csv-exporter"
)

func main() {
	token := flag.String("token", "", "API token for authentication")
	databaseID := flag.String("databaseID", "", "Database ID for the operation")
	sortKey := flag.String("sortKey", "", "Key to sort the database entries")
	order := flag.String("order", "descending", "Order to sort the entries (ascending or descending)")

	flag.Parse()

	if *token == "" || *databaseID == "" {
		log.Println("Error: Both token and databaseID are required.")
		flag.Usage()
		os.Exit(1)
	}
	if *order != "" && *order != "ascending" && *order != "descending" {
		log.Println("Error: Order must be either 'ascending' or 'descending'.")
		flag.Usage()
		os.Exit(1)
	}

	client := exporter.NewExporter(*token)
	err := client.ExportDatabase(
		context.Background(),
		*databaseID,
		exporter.Options{
			SortKey: *sortKey,
			Order:   *order,
		},
		os.Stdout,
	)
	if err != nil {
		log.Fatalf("Failed to export database: %v", err)
	}
}
