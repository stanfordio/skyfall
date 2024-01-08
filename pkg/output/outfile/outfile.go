package outfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/stanfordio/skyfall/pkg/utils"
)

type Outfile struct {
	// OutputFile is the path to the file to write output to
	OutputFilePath string
	OutputChannel  chan map[string]interface{}
}

func (outfile Outfile) GetBackfillSeqno() (int64, error) {
	lastLine, err := utils.GetLastLine(outfile.OutputFilePath)
	if err != nil {
		log.Warnf("Unable to read last line of output file for backfill: %+v", err)
		return 0, err
	}

	// Try to parse the last line as JSON, then pull out the last "seq" field
	lastData := make(map[string]interface{})
	err = json.Unmarshal([]byte(lastLine), &lastData)
	if err != nil {
		return 0, errors.New("unable to parse last line as JSON")
	}
	lastSeqFloat := lastData["_Seq"]
	if lastSeqFloat == nil {
		return 0, errors.New("unable to find seq in last line of output file")
	} else {
		lastSeq := int64(lastSeqFloat.(float64))
		return lastSeq, nil
	}
}

func (outfile Outfile) Setup() error {
	if outfile.OutputFilePath == "" {
		return errors.New("output file path is required")
	}
	return nil
}

func (outfile Outfile) StreamOutput(ctx context.Context) error {
	f, err := os.OpenFile(outfile.OutputFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %+v", err)
		return err
	}
	defer f.Close()

	_, cancel := context.WithCancel(ctx)

	for {
		e := <-outfile.OutputChannel
		// JSON encode the event
		marshaled, err := json.Marshal(e)
		if err != nil {
			log.Errorf("Failed to marshal event: %+v", err)
			cancel()
		}
		if _, err := f.Write(append(marshaled, byte('\n'))); err != nil {
			log.Errorf("Failed to write output: %+v", err)
			cancel()
		}
	}
}
