package statsd

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/varsitytutors/gostatsd"
	"github.com/varsitytutors/gostatsd/pkg/stats"

	log "github.com/sirupsen/logrus"
)

// MetricFlusher periodically flushes metrics from all Aggregators to Senders.
type MetricFlusher struct {
	// Counter fields below must be read/written only using atomic instructions.
	// 64-bit fields must be the first fields in the struct to guarantee proper memory alignment.
	// See https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	lastFlush      int64 // Last time the metrics where aggregated. Unix timestamp in nsec.
	lastFlushError int64 // Time of the last flush error. Unix timestamp in nsec.

	flushInterval      time.Duration // How often to flush metrics to the sender
	aggregateProcesser AggregateProcesser
	backends           []gostatsd.Backend
}

// NewMetricFlusher creates a new MetricFlusher with provided configuration.
func NewMetricFlusher(flushInterval time.Duration, aggregateProcesser AggregateProcesser, backends []gostatsd.Backend) *MetricFlusher {
	return &MetricFlusher{
		flushInterval:      flushInterval,
		aggregateProcesser: aggregateProcesser,
		backends:           backends,
	}
}

// Run runs the MetricFlusher.
func (f *MetricFlusher) Run(ctx context.Context) {
	statser := stats.FromContext(ctx)

	flushTicker := time.NewTicker(f.flushInterval)
	defer flushTicker.Stop()

	lastFlush := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case thisFlush := <-flushTicker.C: // Time to flush to the backends
			flushDelta := thisFlush.Sub(lastFlush)
			if f.aggregateProcesser != AggregateProcesser(nil) {
				f.flushData(ctx, flushDelta, statser)
			}
			statser.NotifyFlush(flushDelta)
			lastFlush = thisFlush
		}
	}
}

func (f *MetricFlusher) flushData(ctx context.Context, flushInterval time.Duration, statser stats.Statser) {
	var sendWg sync.WaitGroup
	timerTotal := statser.NewTimer("flusher.total_time", nil)
	processWait := f.aggregateProcesser.Process(ctx, func(workerId int, aggr Aggregator) {
		// This is in the flusher, but it's an aggregator action, so put it in that space.
		tags := gostatsd.Tags{fmt.Sprintf("aggregator_id:%d", workerId)}

		timerFlush := statser.NewTimer("aggregator.aggregation_time", tags)
		aggr.Flush(flushInterval)
		timerFlush.SendGauge()

		timerProcess := statser.NewTimer("aggregator.process_time", tags)
		aggr.Process(func(m *gostatsd.MetricMap) {
			f.sendMetricsAsync(ctx, &sendWg, m)
		})
		timerProcess.SendGauge()

		timerReset := statser.NewTimer("aggregator.reset_time", tags)
		aggr.Reset()
		timerReset.SendGauge()
	})
	processWait() // Wait for all workers to execute function
	sendWg.Wait() // Wait for all backends to finish sending
	timerTotal.SendGauge()
}

func (f *MetricFlusher) sendMetricsAsync(ctx context.Context, wg *sync.WaitGroup, m *gostatsd.MetricMap) {
	wg.Add(len(f.backends))
	for _, backend := range f.backends {
		backend.SendMetricsAsync(ctx, m, func(errs []error) {
			defer wg.Done()
			f.handleSendResult(errs)
		})
	}
}

func (f *MetricFlusher) handleSendResult(flushResults []error) {
	timestampPointer := &f.lastFlush
	for _, err := range flushResults {
		if err != nil {
			timestampPointer = &f.lastFlushError
			if err != context.DeadlineExceeded && err != context.Canceled {
				log.Errorf("Sending metrics to backend failed: %v", err)
			}
		}
	}
	atomic.StoreInt64(timestampPointer, time.Now().UnixNano())
}
