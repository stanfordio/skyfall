package bq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	bq_schema "github.com/stanfordio/skyfall/pkg/output/bq/schema"
	"google.golang.org/api/iterator"
)

type BQ struct {
	Context       context.Context
	Client        *bigquery.Client
	OutputTable   *bigquery.Table
	OutputChannel chan map[string]interface{}
}

func New(ctx context.Context, tablePath string, outputChannel chan map[string]interface{}) (*BQ, error) {
	tablePathComponents := strings.Split(tablePath, ".")

	// The second to last component is the dataset ID, and the last component is the table ID
	if len(tablePathComponents) < 2 {
		return nil, errors.New("expected the bigquery table path to contain at least a dataset ID and a table ID")
	}

	datasetName := tablePathComponents[len(tablePathComponents)-2]
	tableName := tablePathComponents[len(tablePathComponents)-1]
	projectId := bigquery.DetectProjectID
	if len(tablePathComponents) > 2 {
		projectId = tablePathComponents[len(tablePathComponents)-3]
	}

	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}

	table := client.Dataset(datasetName).Table(tableName)

	bq := BQ{
		Context:       ctx,
		Client:        client,
		OutputTable:   table,
		OutputChannel: outputChannel,
	}

	return &bq, nil
}

// Helper function to compare two schemas
func schemasAreEqual(schema1, schema2 bigquery.Schema) bool {
	if len(schema1) != len(schema2) {
		return false
	}
	for i := range schema1 {
		if schema1[i].Name != schema2[i].Name || schema1[i].Type != schema2[i].Type || schema1[i].Repeated != schema2[i].Repeated || schema1[i].Required != schema2[i].Required {
			return false
		}
		// If a record, check that the nested schemas are equal
		if schema1[i].Type == bigquery.RecordFieldType {
			if !schemasAreEqual(schema1[i].Schema, schema2[i].Schema) {
				return false
			}
		}
	}
	return true
}

func (bq BQ) Setup() error {
	schema := bq_schema.GetSchema()

	metadata, err := bq.OutputTable.Metadata(bq.Context)
	if err == nil && metadata != nil {
		log.Infof("Found existing BigQuery table: %+v", bq.OutputTable.FullyQualifiedName())

		// Verify that the schemas are compatible
		if !schemasAreEqual(schema, metadata.Schema) {
			log.Errorf("Found schema mismatch between existing table and output data")
			foundSchema, err1 := metadata.Schema.ToJSONFields()
			if err1 != nil {
				log.Errorf("Failed to convert existing schema to JSON: %s", err1)
				return err1
			}

			desiredSchema, err2 := schema.ToJSONFields()
			if err2 != nil {
				log.Errorf("Failed to convert desired schema to JSON: %s", err2)
				return err2
			}

			log.Errorf("Found schema: %s", foundSchema)
			log.Errorf("Desired schema: %s", desiredSchema)
			return errors.New("the schema of the existing table does not match the schema of the output data; please update the table schema manually")
		} else {
			log.Infof("Schema of table is compatible with output data!")
		}
	} else {
		log.Infof("Table does not exist, so creating new BigQuery table...")
		// Create or update the table
		err := bq.OutputTable.Create(bq.Context, &bigquery.TableMetadata{
			Schema: schema,
			TimePartitioning: &bigquery.TimePartitioning{
				Type: bigquery.MonthPartitioningType,
			},
		})
		if err != nil {
			return err
		}
		log.Infof("Created new BigQuery table: %+v", bq.OutputTable.FullyQualifiedName())
	}

	return nil
}

func (bq BQ) GetBackfillSeqno() (int64, error) {
	query := fmt.Sprintf("SELECT MAX(Seq) as max_seq FROM `%s.%s.%s`", bq.OutputTable.ProjectID, bq.OutputTable.DatasetID, bq.OutputTable.TableID)
	log.Infof("Running query: %s", query)
	result := bq.Client.Query(query)
	it, err := result.Read(context.Background())
	if err != nil {
		return 0, err
	}
	var maxSeq int64
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		if row["max_seq"] != nil {
			maxSeq = row["max_seq"].(int64)
		} else {
			return 0, errors.New("unable to find max_seq in output table")
		}
	}
	return maxSeq, nil
}

func cleanOutput(value map[string]interface{}) map[string]interface{} {
	for k, v := range value {
		// Check if the key starts with "$" and replace it with "_"
		if strings.HasPrefix(k, "$") {
			newKey := "_" + k[1:]
			value[newKey] = v
			delete(value, k)
		}

		// If the value is a map, recursively clean it
		if subMap, ok := v.(map[string]interface{}); ok {
			value[k] = cleanOutput(subMap)
		}
	}
	return value
}

type MapValueSaver struct {
	Data map[string]interface{}
}

func (mvs MapValueSaver) Save() (row map[string]bigquery.Value, insertID string, err error) {
	row = make(map[string]bigquery.Value)
	for k, v := range mvs.Data {
		row[k] = v
	}
	// Generate a unique insertID if needed, or return "" if not
	insertID = ""
	return row, insertID, nil
}

func prepareForWrite(value map[string]interface{}) (out *MapValueSaver, err error) {
	// Set "Full" to the JSON representation of "Full"
	fullMarshalled, err := json.Marshal(value["Full"])
	if err != nil {
		log.Errorf("Failed to marshal event: %+v", err)
		return nil, err
	}
	value["Full"] = string(fullMarshalled)

	cleaned := cleanOutput(value)

	// Insert the "_Raw" field as a string
	json, err := json.Marshal(cleaned)
	if err != nil {
		log.Errorf("Failed to marshal event: %+v", err)
		return nil, err
	}

	cleaned["_Raw"] = string(json)

	mvs := MapValueSaver{Data: cleaned} // Heap allocated, thanks Go.

	return &mvs, nil
}

func (bq BQ) StreamOutput(ctx context.Context) error {
	log.Infof("Streaming output to BigQuery table: %+v", bq.OutputTable.FullyQualifiedName())

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	inserter := bq.OutputTable.Inserter()
	inserter.IgnoreUnknownValues = true

MainLoop:
	for {
		var values []map[string]interface{}

		// Wait for at least one value in the channel
		value, ok := <-bq.OutputChannel
		if !ok {
			break MainLoop // Channel is closed
		}
		values = append(values, value)

		// Collect remaining values from the channel until it is empty
	ChannelCollector:
		for {
			select {
			case value, ok := <-bq.OutputChannel:
				if !ok {
					break MainLoop // Channel is closed
				}

				values = append(values, value)

				// If there are more than 250 values, write them now
				if len(values) >= 250 {
					break ChannelCollector
				}
			default:
				break ChannelCollector // Channel is empty
			}
		}

		// Prepare the values for writing
		var preparedValues []*MapValueSaver
		for _, value := range values {
			preparedValue, err := prepareForWrite(value)
			if err != nil {
				log.Errorf("Failed to prepare value for writing: %+v", err)
				continue
			}
			preparedValues = append(preparedValues, preparedValue)
		}

		// Do the write in a batch (not a BigQuery *batch*, just, like, a
		// semantic batch)
		if err := inserter.Put(ctx, preparedValues); err != nil {
			log.Errorf("Failed to write output: %+v", err)
		}
		log.Infof("Wrote %d rows to BigQuery", len(values))
	}

	return nil
}
