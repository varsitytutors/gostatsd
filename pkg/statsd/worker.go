package statsd

import (
	"context"
	"fmt"
	"time"

	"github.com/varsitytutors/gostatsd"
	"github.com/varsitytutors/gostatsd/pkg/stats"
)

type processCommand struct {
	f    DispatcherProcessFunc
	done func()
}

type worker struct {
	aggr           Aggregator
	metricsQueue   chan []*gostatsd.Metric
	metricMapQueue chan *gostatsd.MetricMap
	processChan    chan *processCommand
	id             int
}

func (w *worker) work() {
	for {
		select {
		case metrics, ok := <-w.metricsQueue:
			if !ok {
				return
			}
			w.aggr.Receive(metrics...)
		case mm, ok := <-w.metricMapQueue:
			if !ok {
				return
			}
			w.aggr.ReceiveMap(mm)
		case cmd := <-w.processChan:
			w.executeProcess(cmd)
		}
	}
}

func (w *worker) executeProcess(cmd *processCommand) {
	defer cmd.done() // Done with the process command
	cmd.f(w.id, w.aggr)
}

func (w *worker) RunMetrics(ctx context.Context, statser stats.Statser) {
	csw := stats.NewChannelStatsWatcher(
		statser,
		"dispatch_aggregator",
		gostatsd.Tags{fmt.Sprintf("aggregator_id:%d", w.id)},
		cap(w.metricsQueue),
		func() int { return len(w.metricsQueue) },
		1000*time.Millisecond, // TODO: Make configurable
	)
	csw.Run(ctx)
}
