package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"

	exporter "github.com/yudppp/notion-csv-exporter"
)

func main() {
	token := flag.String("token", "", "API token for authentication")
	databaseID := flag.String("databaseID", "", "Database ID for the operation")

	flag.Parse()

	if *token == "" || *databaseID == "" {
		fmt.Println("Error: Both token and databaseID are required.")
		flag.Usage()
		os.Exit(1)
	}

	client := exporter.NewExporter(*token)
	w := &bytes.Buffer{}
	err := client.ExportDatabase(context.Background(), *databaseID, w)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(w.String())
}
