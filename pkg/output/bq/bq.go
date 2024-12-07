package bq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/bigquery/storage/managedwriter"
	adapt "cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	log "github.com/sirupsen/logrus"
	bq_schema "github.com/stanfordio/skyfall/pkg/output/bq/schema"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/encoding/protojson"   
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"   
	"google.golang.org/protobuf/types/dynamicpb"
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

func sortSchema(schema bigquery.Schema) {
	sort.Slice(schema, func(i, j int) bool {
		return schema[i].Name < schema[j].Name
	})
	for _, field := range schema {
		if field.Type == bigquery.RecordFieldType {
			sortSchema(field.Schema)
		}
	}
}

// Helper function to compare two schemas
func schemasAreEqual(schema1, schema2 bigquery.Schema) bool {
	sortSchema(schema1)
	sortSchema(schema2)
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

	log.SetFormatter(&log.TextFormatter{
		DisableQuote: true,
	})

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

// Clean the map values to handle pointers, nested structures, and key transformations
func cleanOutput(value map[string]interface{}) map[string]interface{} {
	for k, v := range value {
		// Handle pointer dereferencing
		switch v := v.(type) {
		case *string:
			value[k] = derefString(v)
		case *int64:
			value[k] = derefInt64(v)
		case *float64:
			value[k] = derefFloat64(v)
		}

		// Transform keys starting with "$" to "_"
		if strings.HasPrefix(k, "$") {
			newKey := "_" + k[1:]
			value[newKey] = value[k]
			delete(value, k)
		}

		// Recursively clean nested maps
		if subMap, ok := v.(map[string]interface{}); ok {
			value[k] = cleanOutput(subMap)
		}

		// Handle specific fields like timestamps
		if k == "CreatedAt" || k == "PulledTimestamp" || k == "IndexedAt" {
			value[k] = parseTimestamp(value[k])
		}

		// Convert "Full" to JSON string if it's a map
		if k == "Full" {
			value[k] = convertFullToJSON(value[k])
		}
	}
	return value
}

func derefString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func derefInt64(p *int64) int64 {
	if p != nil {
		return *p
	}
	return 0
}

func derefFloat64(p *float64) float64 {
	if p != nil {
		return *p
	}
	return 0.0
}

func parseTimestamp(v interface{}) int64 {
	if t, ok := v.(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, t)
		if err != nil {
			log.Printf("Invalid timestamp format: %v", err)
			return 0
		}
		return parsedTime.Unix()
	}
	return 0
}

func convertFullToJSON(v interface{}) string {
	switch v := v.(type) {
	case map[string]interface{}:
		fullMarshalled, err := json.Marshal(v)
		if err != nil {
			log.Printf("Failed to marshal 'Full' field: %v", err)
			return ""
		}
		return string(fullMarshalled)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func prepareForWrite(value map[string]interface{}) ([]byte, error) {
	cleaned := cleanOutput(value)

	// Remove _Raw field at the root level
	delete(cleaned, "_Raw")

	// Remove ReplyCount field under Projection -> LikedPost
	if projection, ok := cleaned["Projection"].(map[string]interface{}); ok {
		if likedPost, ok := projection["LikedPost"].(map[string]interface{}); ok {
			delete(likedPost, "ReplyCount")
		}
		// Remove ReplyCount field under Projection -> RepostedPost
		if repostedPost, ok := projection["RepostedPost"].(map[string]interface{}); ok {
			delete(repostedPost, "ReplyCount")
		}
	}

	// Marshal the cleaned map into JSON bytes
	return json.Marshal(cleaned)
}

// Streams data to BigQuery
func (bq BQ) StreamOutput(ctx context.Context) error {
	log.Infof("Starting to stream output to BigQuery table: %+v", bq.OutputTable.FullyQualifiedName())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up the schema descriptors
	tableSchema := bq_schema.GetSchema()
	messageDescriptor, descriptor, err := setupDynamicDescriptors(tableSchema)
	if err != nil {
		log.Errorf("Failed to create descriptors: %v", err)
		return err
	}

	client, err := managedwriter.NewClient(ctx, bigquery.DetectProjectID)
	if err != nil {
		log.Errorf("Failed to create managed writer client: %v", err)
		return err
	}
	defer client.Close()

	destinationTable := fmt.Sprintf("projects/%s/datasets/%s/tables/%s", bq.OutputTable.ProjectID, bq.OutputTable.DatasetID, bq.OutputTable.TableID)
	managedStream, err := client.NewManagedStream(ctx,
		managedwriter.WithSchemaDescriptor(descriptor),
		managedwriter.WithDestinationTable(destinationTable),
	)
	if err != nil {
		log.Errorf("Failed to create managed stream: %v", err)
		return err
	}

	// Stream processing loop
	var buffer []map[string]interface{}
	for {
		select {
		case value, ok := <-bq.OutputChannel:
			if !ok {
				log.Info("Channel closed, exiting.")
				return nil
			}
			buffer = append(buffer, value)

			// Flush if buffer size reaches threshold
			if len(buffer) >= 250 {
				if err := bq.flushBuffer(ctx, managedStream, messageDescriptor, buffer); err != nil {
					log.Errorf("Failed to flush buffer: %v", err)
					continue
				} else {
					// Log the number of rows uploaded to BigQuery
					log.Infof("Buffer flushed! Rows uploaded: %d", len(buffer))
				}
				buffer = nil // Clear the buffer after flushing
			}

		case <-ctx.Done():
			log.Warn("Context canceled, stopping.")
			return ctx.Err()
		}
	}
}

func (bq BQ) flushBuffer(ctx context.Context, managedStream *managedwriter.ManagedStream, descriptor protoreflect.MessageDescriptor, buffer []map[string]interface{}) error {
	var encodedRows [][]byte
	for _, value := range buffer {
		preparedValue, err := prepareForWrite(value)
		if err != nil {
			log.Errorf("Failed to prepare value for writing: %+v", err)
			continue
		}

		message := dynamicpb.NewMessage(descriptor)
		if err := protojson.Unmarshal(preparedValue, message); err != nil {
			log.Errorf("Failed to unmarshal into proto message: %+v", err)
			continue
		}

		encodedRow, err := proto.Marshal(message)
		if err != nil {
			log.Errorf("Failed to marshal proto message: %+v", err)
			continue
		}
		encodedRows = append(encodedRows, encodedRow)
	}

	if len(encodedRows) == 0 {
		log.Warn("No rows to write.")
		return nil
	}

	result, err := managedStream.AppendRows(ctx, encodedRows)
	if err != nil {
		return fmt.Errorf("failed to append rows: %w", err)
	}

	_, err = result.GetResult(ctx)
	return err
}

func setupDynamicDescriptors(schema bigquery.Schema) (protoreflect.MessageDescriptor, *descriptorpb.DescriptorProto, error) {
	// Convert BigQuery schema to storage schema
	convertedSchema, err := adapt.BQSchemaToStorageTableSchema(schema)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert BigQuery schema to storage schema: %v", err)
	}

	// Convert storage schema to proto2 descriptor
	descriptor, err := adapt.StorageSchemaToProto2Descriptor(convertedSchema, "root")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert storage schema to proto2 descriptor: %v", err)
	}

	// Ensure the descriptor is a MessageDescriptor
	messageDescriptor, ok := descriptor.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, nil, fmt.Errorf("descriptor is not a MessageDescriptor")
	}

	// Normalize the descriptor
	dp, err := adapt.NormalizeDescriptor(messageDescriptor)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to normalize descriptor: %v", err)
	}

	return messageDescriptor, dp, nil
}
