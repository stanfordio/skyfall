package bq

import (
	"context"
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

func New(ctx context.Context, tablePath string) (*BQ, error) {
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
		Context:     ctx,
		Client:      client,
		OutputTable: table,
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
				log.Errorf("Failed to convert existing schema to JSON: %+v", err1)
				return err1
			}

			desiredSchema, err2 := schema.ToJSONFields()
			if err2 != nil {
				log.Errorf("Failed to convert desired schema to JSON: %+v", err2)
				return err2
			}

			log.Errorf("Found schema: %+v", foundSchema)
			log.Errorf("Desired schema: %+v", desiredSchema)
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
	query := fmt.Sprintf("SELECT MAX(_Seq) as max_seq FROM `%s.%s.%s`", bq.OutputTable.ProjectID, bq.OutputTable.DatasetID, bq.OutputTable.TableID)
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
	// Recursively replace all fields whose keys start with "$" with "_"
	for k, v := range value {
		if strings.HasPrefix(k, "$") {
			value["_"+k[1:]] = v
			delete(value, k)
		}
		if vMap, ok := v.(map[string]interface{}); ok {
			value[k] = cleanOutput(vMap)
		}
	}
	return value
}

func (bq BQ) StreamOutput(ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)
	inserter := bq.OutputTable.Inserter()
	inserter.IgnoreUnknownValues = true

	for {
		e := <-bq.OutputChannel
		cleaned := cleanOutput(e)
		if err := inserter.Put(ctx, cleaned); err != nil {
			log.Errorf("Failed to write output: %+v", err)
			cancel()
		} else {
			log.Infof("Wrote output: %+v", cleaned)
		}
	}
}
