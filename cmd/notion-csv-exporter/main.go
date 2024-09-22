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

	flag.Parse()

	if *token == "" || *databaseID == "" {
		log.Println("Error: Both token and databaseID are required.")
		flag.Usage()
		os.Exit(1)
	}
	client := exporter.NewExporter(*token)
	err := client.ExportDatabase(context.Background(), *databaseID, os.Stdout)
	if err != nil {
		log.Fatalf("Failed to export database: %v", err)
	}
}
