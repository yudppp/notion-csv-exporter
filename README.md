# Notion CSV Exporter

This is a simple tool to export data from a Notion database into a CSV file using the Notion API.

## Features

- Exports data from a specified Notion database to a CSV file.
- Supports customization of Notion API tokens and database IDs via command line arguments.

## Requirements

Before you start, make sure you have:

- A valid Notion Integration Token
- The Notion Database ID that you want to export
- Go installed on your machine

## Installation

You can install `notion-csv-exporter` using `go install`:

```bash
go install github.com/yudppp/notion-csv-exporter/cmd/notion-csv-exporter@latest
```

This will download and build the tool, placing the binary in your `$GOBIN` directory (by default `$HOME/go/bin`).

If you want to install a specific version, you can use a tagged version like this:

```bash
go install github.com/yudppp/notion-csv-exporter/cmd/notion-csv-exporter@v1.0.0
```

Make sure your `$GOBIN` is added to your system's `$PATH`, so you can run the tool from anywhere in your terminal.

## Usage

Once installed, you can use the tool from anywhere in your terminal:

```bash
notion-csv-exporter -token={NOTION_API_TOKEN} -databaseID={NOTION_DATABASE_ID}
```

Replace `{NOTION_API_TOKEN}` with your Notion API token and `{NOTION_DATABASE_ID}` with your database ID.

### Example

Here’s an example command to export data:

```bash
notion-csv-exporter -token=secret_abc12345 -databaseID=1234567890abcdef1234567890abcdef
```

This command will export the contents of the specified Notion database to a CSV file named `output.csv` in the current directory.

## How to Get Notion API Token and Database ID

### API Token

1. Go to [Notion Integrations](https://www.notion.so/my-integrations).
2. Create a new integration.
3. Copy the "Internal Integration Token."

### Database ID

1. Open the Notion database in your browser.
2. The database ID is the part of the URL between `/` and `?` or the last part after `/` if there is no query string.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
