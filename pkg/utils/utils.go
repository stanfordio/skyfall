package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

// From https://stackoverflow.com/questions/17863821/how-to-read-last-lines-from-a-big-file-with-go-every-10-secs
func GetLastLine(filepath string) (string, error) {
	fileHandle, err := os.Open(filepath)

	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	return line, nil
}

func RetryingHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 15
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 15 * time.Minute
	retryClient.CheckRetry = XRPCRetryPolicy
	client := retryClient.StandardClient()
	client.Timeout = 30 * time.Second

	return client
}

func XRPCRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Do not retry network errors, since these are usually because the PDS is dead
	if err != nil {
		if _, ok := err.(*net.OpError); !ok {
			log.Warnf("Not retrying on error: %s", err.Error())
			return false, err
		}
	}

	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
