package output

import (
	"context"

	"github.com/stanfordio/skyfall/pkg/output/bq"
	"github.com/stanfordio/skyfall/pkg/output/outfile"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

type Output interface {
	Setup() error
	GetBackfillSeqno() (int64, error)
	StreamOutput(context.Context) error
}

func NewOutput(cctx *cli.Context, outputChannel chan map[string]interface{}) (Output, error) {
	if cctx.String("output-bq-table") != "" {
		log.Infof("output-bq-table specified, so writing output to BigQuery table: %s", cctx.String("output-bq-table"))
		bq, err := bq.New(cctx.Context, cctx.String("output-bq-table"), outputChannel)
		if err != nil {
			log.Fatalf("Failed to create BigQuery output: %+v", err)
			return nil, err
		}
		return bq, nil
	}

	if cctx.String("output-file") != "" {
		log.Infof("output-file specified, so writing output to file: %s", cctx.String("output-file"))
		return outfile.Outfile{
			OutputFilePath: cctx.String("output-file"),
			OutputChannel:  outputChannel,
			StringifyFull:  cctx.Bool("stringify-full"),
		}, nil
	}

	return nil, nil
}
