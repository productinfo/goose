package statsd

import "context"

// Gaugor represents a gauge-type metric, which takes absolute values.
// https://docs.datadoghq.com/developers/metrics/gauges/
type Gaugor Collector

// Gauge takes a float64 and sets the indicated gauge in statsd to this value.
//
// The last parameter is an arbitrary array of tags as maps.
func (g *Gaugor) Gauge(ctx context.Context, n float64, ts ...Tags) {
	tags := loadTags(ctx, g.Tags, ts...)
	Gauge(ctx, g.Name, n, tags, g.Rate.Rate())
}
