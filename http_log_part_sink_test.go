package worker

import (
	"testing"

	gocontext "context"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPLogPartSink(t *testing.T) {
	ctx, cancel := gocontext.WithCancel(gocontext.TODO())
	cancel()
	lps := newHTTPLogPartSink(
		ctx,
		"http://example.org/log-parts/multi",
		uint64(1000))

	assert.NotNil(t, lps)
}

func TestHTTPLogPartSink_flush(t *testing.T) {
	ctx := gocontext.TODO()
	var lps *httpLogPartSink
	for _, lpsValue := range httpLogPartSinksByURL {
		lps = lpsValue
	}
	lps.flush(gocontext.TODO())
	lps.Add(ctx, &httpLogPart{
		JobID:   uint64(4),
		Content: "wat",
		Number:  3,
		Final:   false,
	})

	assert.Len(t, lps.partsBuffer, 1)
	lps.flush(ctx)
	assert.Len(t, lps.partsBuffer, 0)
}
