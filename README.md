# Skyfall

Skyfall is a tool for downloading data from Bluesky. It can:

* Download every (known) Bluesky user's CAR file and store it on the disk
* Attach to the Bluesky firehose and output structured, hydrated data (with backfill)
* Turn a folder filled with CAR files into structured, hydrated data (in the _same format_ as the Bluesky firehose output)
* Output hydrated data into JSONL or BigQuery

## Development

Skyfall is a fairly simple Go CLI tool. Just download the Go toolchain and run `go build` to build the binary. You can also run the binary directly with `go run cmd/main.go`.

## Usage

```
NAME:
   skyfall - A simple CLI for Bluesky data ingest

USAGE:
   skyfall [global options] command [command options] [arguments...]

VERSION:
   prerelease

COMMANDS:
   stream    Sip from the firehose
   repodump  Dump everyone's repos (as CAR) into a folder
   hydrate   Hydrate CAR pulls into the same format as the stream
   help      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cache-size value  maximum size of the cache, in bytes (default: 4294967296)
   --handle value      handle to authenticate with, e.g., miles.land or det.bsky.social
   --password value    password to authenticate with
   --help, -h          show help
   --version, -v       print the version
```

Most commands require authentication. You can authenticate with the `--handle` and `--password` flags. Typically, that will look like: `go run cmd/main.go --handle <handle> --password <password> command ...`.

### Stream

```
NAME:
   skyfall stream - Sip from the firehose

USAGE:
   skyfall stream [command options] [arguments...]

OPTIONS:
   --worker-count value     number of workers to scale to (default: 32)
   --output-file value      file to write output to (if specified, will attempt to backfill from the most recent event in the file) (default: "output.jsonl")
   --stringify-full         whether to stringify the full event in file output (if true, the JSON will be stringified; this is helpful when you want output to match what would be sent to BigQuery) (default: false)
   --output-bq-table value  name of a BigQuery table to output to in ID form (e.g., dgap_bsky.example_table)
   --backfill-seq value     seq to backfill from (if specified, will override the seqno extracted from the output file/bigquery table) (default: 0)
   --autorestart            automatically restart the stream if it dies (default: true)
   --help, -h               show help
```

Example usage:

```
go run cmd/main.go --handle <handle> --password <password> stream --output-file output.jsonl
go run cmd/main.go --handle <handle> --password <password> stream --output-bq-table dgap_bsky.example_table
```

### Dump repos

```
NAME:
   skyfall repodump - Dump everyone's repos (as CAR) into a folder

USAGE:
   skyfall repodump [command options] [arguments...]

OPTIONS:
   --output-folder value  folder to write repos to (default: "output")
   --help, -h             show help
```

The folder will be broken up into subfolders based on the DID to prevent any one folder from getting too large.

Example usage:

```
go run cmd/main.go --handle <handle> --password <password> repodump --output-folder repos
```

### Hydrate

```
NAME:
   skyfall hydrate - Hydrate CAR pulls into the same format as the stream

USAGE:
   skyfall hydrate [command options] [arguments...]

OPTIONS:
   --input value            folder or file to read data from
   --worker-count value     number of workers to scale to (default: 32)
   --output-file value      file to write output to (if specified, will attempt to backfill from the most recent event in the file) (default: "output.jsonl")
   --output-bq-table value  name of a BigQuery table to output to in ID form (e.g., dgap_bsky.example_table)
   --help, -h               show help
```

Example usage:

```
go run cmd/main.go --handle <handle> --password <password> hydrate --input repos --output-file output.jsonl
go run cmd/main.go --handle <handle> --password <password> hydrate --input repos --output-bq-table dgap_bsky.example_table
```

## BigQuery

Skyfall can output to BigQuery. To do so, you'll need to authenticate to Google using the `GOOGLE_APPLICATION_CREDENTIALS` environment variable. You can set this to the path of a service account JSON file.

## License

Skyfall is licensed under the Apache 2.0 license. See [LICENSE](LICENSE) for more details.
