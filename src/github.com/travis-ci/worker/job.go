package worker

import (
	"io"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/travis-ci/worker/backend"
	"golang.org/x/net/context"
)

// JobPayload is the payload we receive over RabbitMQ.
type JobPayload struct {
	Type       string                 `json:"type"`
	Job        JobJobPayload          `json:"job"`
	Build      BuildPayload           `json:"source"`
	Repository RepositoryPayload      `json:"repository"`
	UUID       string                 `json:"uuid"`
	Config     map[string]interface{} `json:"config"`
	Timeouts   TimeoutsPayload        `json:"timeouts,omitempty"`
}

// JobJobPayload contains information about the job.
type JobJobPayload struct {
	ID     uint64 `json:"id"`
	Number string `json:"number"`
}

// BuildPayload contains information about the build.
type BuildPayload struct {
	ID     uint64 `json:"id"`
	Number string `json:"number"`
}

// RepositoryPayload contains information about the repository.
type RepositoryPayload struct {
	ID   uint64 `json:"id"`
	Slug string `json:"slug"`
}

// TimeoutsPayload contains information about any custom timeouts. The timeouts
// are given in seconds, and a value of 0 means no custom timeout is set.
type TimeoutsPayload struct {
	HardLimit  uint64 `json:"hard_limit"`
	LogSilence uint64 `json:"log_silence"`
}

// FinishState is the state that a job finished with (such as pass/fail/etc.).
// You should not provide a string directly, but use one of the FinishStateX
// constants defined in this package.
type FinishState string

// Valid finish states for the FinishState type
const (
	FinishStatePassed    FinishState = "passed"
	FinishStateFailed    FinishState = "failed"
	FinishStateErrored   FinishState = "errored"
	FinishStateCancelled FinishState = "cancelled"
)

// A Job ties togeher all the elements required for a build job
type Job interface {
	Payload() *JobPayload
	RawPayload() *simplejson.Json
	StartAttributes() *backend.StartAttributes

	Received() error
	Started() error
	Error(context.Context, string) error
	Requeue() error
	Finish(FinishState) error

	LogWriter(context.Context) (LogWriter, error)
}

// A LogWriter is primarily an io.Writer that will send all bytes to travis-logs
// for processing, and also has some utility methods for timeouts and log length
// limiting. Each LogWriter is tied to a given job, and can be gotten by calling
// the LogWriter() method on a Job.
type LogWriter interface {
	io.WriteCloser
	WriteAndClose([]byte) (int, error)
	SetTimeout(time.Duration)
	Timeout() <-chan time.Time
	SetMaxLogLength(int)
}
