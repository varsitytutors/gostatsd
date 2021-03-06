package statsd

import (
	"bytes"
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/varsitytutors/gostatsd"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagStripMergesCounters(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, []Filter{
		{DropTags: gostatsd.StringMatchList{gostatsd.NewStringMatch("key2:*")}},
	})
	mm := gostatsd.NewMetricMap()
	mm.Receive(&gostatsd.Metric{Type: gostatsd.COUNTER, Name: "metric", Timestamp: 20, Tags: gostatsd.Tags{"key:value", "key2:value2"}, Value: 1, Rate: 0.1})
	mm.Receive(&gostatsd.Metric{Type: gostatsd.COUNTER, Name: "metric", Timestamp: 10, Tags: gostatsd.Tags{"key:value"}, Value: 20, Rate: 1})
	th.DispatchMetricMap(context.Background(), mm)

	expected := gostatsd.NewMetricMap()
	expected.Counters["metric"] = map[string]gostatsd.Counter{
		"key:value": {Timestamp: 20, Value: 30, Tags: gostatsd.Tags{"key:value"}},
	}
	require.EqualValues(t, expected, tch.mm[0])
}

func TestTagStripMergesGauges(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, []Filter{
		{DropTags: gostatsd.StringMatchList{gostatsd.NewStringMatch("key2:*")}},
	})
	mm := gostatsd.NewMetricMap()
	mm.Receive(&gostatsd.Metric{Type: gostatsd.GAUGE, Name: "metric", Timestamp: 10, Tags: gostatsd.Tags{"key:value", "key2:value2"}, Value: 10})
	mm.Receive(&gostatsd.Metric{Type: gostatsd.GAUGE, Name: "metric", Timestamp: 20, Tags: gostatsd.Tags{"key:value"}, Value: 20})
	th.DispatchMetricMap(context.Background(), mm)

	expected := gostatsd.NewMetricMap()
	expected.Gauges["metric"] = map[string]gostatsd.Gauge{
		"key:value": {Timestamp: 20, Value: 20, Tags: gostatsd.Tags{"key:value"}},
	}
	require.EqualValues(t, expected, tch.mm[0])
}

func TestTagStripMergesTimers(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, []Filter{
		{DropTags: gostatsd.StringMatchList{gostatsd.NewStringMatch("key2:*")}},
	})
	mm := gostatsd.NewMetricMap()
	mm.Receive(&gostatsd.Metric{Type: gostatsd.TIMER, Name: "metric", Timestamp: 10, Tags: gostatsd.Tags{"key:value", "key2:value2"}, Value: 10, Rate: 1})
	mm.Receive(&gostatsd.Metric{Type: gostatsd.TIMER, Name: "metric", Timestamp: 20, Tags: gostatsd.Tags{"key:value"}, Value: 20, Rate: 1})
	th.DispatchMetricMap(context.Background(), mm)

	// Make sure the actual values are deterministic
	sort.Float64s(tch.mm[0].Timers["metric"]["key:value"].Values)

	expected := gostatsd.NewMetricMap()
	expected.Timers["metric"] = map[string]gostatsd.Timer{
		"key:value": {Timestamp: 20, Values: []float64{10, 20}, Tags: gostatsd.Tags{"key:value"}, SampledCount: 2},
	}
	require.EqualValues(t, expected, tch.mm[0])
}

func TestTagStripMergesSets(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, []Filter{
		{DropTags: gostatsd.StringMatchList{gostatsd.NewStringMatch("key2:*")}},
	})
	mm := gostatsd.NewMetricMap()
	mm.Receive(&gostatsd.Metric{Type: gostatsd.SET, Name: "metric", Timestamp: 10, Tags: gostatsd.Tags{"key:value", "key2:value2"}, StringValue: "abc"})
	mm.Receive(&gostatsd.Metric{Type: gostatsd.SET, Name: "metric", Timestamp: 20, Tags: gostatsd.Tags{"key:value"}, StringValue: "def"})
	th.DispatchMetricMap(context.Background(), mm)

	expected := gostatsd.NewMetricMap()
	expected.Sets["metric"] = map[string]gostatsd.Set{
		"key:value": {Timestamp: 20, Values: map[string]struct{}{"abc": {}, "def": {}}, Tags: gostatsd.Tags{"key:value"}},
	}
	require.EqualValues(t, expected, tch.mm[0])
}

func TestFilterPassesNoFilters(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	m := &gostatsd.Metric{
		Name: "name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	expected := []*gostatsd.Metric{
		{
			Name: "name",
			Tags: gostatsd.Tags{
				"foo:bar",
				"host:baz",
			},
			Hostname: "baz",
		},
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, expected, tch.m)
}

func TestFilterPassesEmptyFilters(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{}
	m := &gostatsd.Metric{
		Name: "name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	expected := []*gostatsd.Metric{
		{
			Name: "name",
			Tags: gostatsd.Tags{
				"foo:bar",
				"host:baz",
			},
			Hostname: "baz",
		},
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, expected, tch.m)
}

func TestFilterKeepNonMatch(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{
		{
			MatchMetrics: gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.name")},
			DropMetric:   true,
		},
	}
	m := &gostatsd.Metric{
		Name: "good.name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	expected := []*gostatsd.Metric{
		{
			Name: "good.name",
			Tags: gostatsd.Tags{
				"foo:bar",
				"host:baz",
			},
			Hostname: "baz",
		},
	}
	assert.Equal(t, expected, tch.m)
}

func TestFilterDropsBadName(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{
		{
			MatchMetrics: gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.name")},
			DropMetric:   true,
		},
	}
	m := &gostatsd.Metric{
		Name: "bad.name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 0, len(tch.m))
}

func TestFilterDropsBadPrefix(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{
		{
			MatchMetrics: gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.*")},
			DropMetric:   true,
		},
	}
	m := &gostatsd.Metric{
		Name: "bad.name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 0, len(tch.m))
}

func TestFilterKeepsWhitelist(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{
		{
			MatchMetrics:   gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.*")},
			ExcludeMetrics: gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.good")},
			DropMetric:     true,
		},
	}

	m := &gostatsd.Metric{
		Name: "bad.name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})

	m = &gostatsd.Metric{
		Name: "bad.good",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})

	expected := []*gostatsd.Metric{
		{
			Name: "bad.good",
			Tags: gostatsd.Tags{
				"foo:bar",
				"host:baz",
			},
			Hostname: "baz",
		},
	}
	assert.Equal(t, expected, tch.m)
}

func TestFilterDropsTag(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{
		{
			MatchMetrics: gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.name")},
			DropTags:     gostatsd.StringMatchList{gostatsd.NewStringMatch("foo:*")},
		},
	}

	m := &gostatsd.Metric{
		Name: "bad.name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})

	expected := []*gostatsd.Metric{
		{
			Name: "bad.name",
			Tags: gostatsd.Tags{
				"host:baz",
			},
			Hostname: "baz",
		},
	}
	assert.Equal(t, expected, tch.m)
}

func TestFilterDropsHost(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	th.filters = []Filter{
		{
			MatchMetrics: gostatsd.StringMatchList{gostatsd.NewStringMatch("bad.name")},
			DropHost:     true,
		},
	}

	m := &gostatsd.Metric{
		Name: "bad.name",
		Tags: gostatsd.Tags{
			"foo:bar",
			"host:baz",
		},
		Hostname: "baz",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})

	expected := []*gostatsd.Metric{
		{
			Name: "bad.name",
			Tags: gostatsd.Tags{
				"foo:bar",
				"host:baz",
			},
			Hostname: "",
		},
	}
	assert.Equal(t, expected, tch.m)
}

func TestNewTagHandlerFromViper(t *testing.T) {
	var data = []byte(`
filters='drop-noisy-metric drop-noisy-metric-with-tag drop-noisy-tag drop-noisy-keep-quiet-metric drop-host'

[filter.drop-noisy-metric]
match-metrics='noisy.*'
drop-metric=true

[filter.drop-noisy-metric-with-tag]
match-metrics='noisy.*'
match-tags='noisy-tag:*'
drop-metric=true

[filter.drop-noisy-tag]
match-metrics='noisy.*'
drop-tags='noisy-tag:*'

[filter.drop-noisy-keep-quiet-metric]
match-metrics='noisy.*'
exclude-metrics='noisy.quiet.* noisy.ok.*'
drop-metric=true

[filter.drop-host]
match-metrics='global.*'
drop-host=true
drop-tags='host:*'
`)

	v := viper.New()
	v.SetConfigType("toml")
	err := v.ReadConfig(bytes.NewBuffer(data))
	assert.NoError(t, err)
	if err != nil {
		return
	}

	nh := &nopHandler{}
	th := NewTagHandlerFromViper(v, nh, nil)

	empty := gostatsd.StringMatchList{}

	expected := []Filter{
		{MatchMetrics: toStringMatch([]string{"noisy.*"}), ExcludeMetrics: empty, MatchTags: empty, DropTags: empty, DropMetric: true, DropHost: false},
		{MatchMetrics: toStringMatch([]string{"noisy.*"}), ExcludeMetrics: empty, MatchTags: toStringMatch([]string{"noisy-tag:*"}), DropTags: empty, DropMetric: true, DropHost: false},
		{MatchMetrics: toStringMatch([]string{"noisy.*"}), ExcludeMetrics: empty, MatchTags: empty, DropTags: toStringMatch([]string{"noisy-tag:*"}), DropMetric: false, DropHost: false},
		{MatchMetrics: toStringMatch([]string{"noisy.*"}), ExcludeMetrics: toStringMatch([]string{"noisy.quiet.*", "noisy.ok.*"}), DropTags: empty, MatchTags: empty, DropMetric: true, DropHost: false},
		{MatchMetrics: toStringMatch([]string{"global.*"}), ExcludeMetrics: empty, MatchTags: empty, DropTags: toStringMatch([]string{"host:*"}), DropMetric: false, DropHost: true},
	}
	assert.Equal(t, expected, th.filters)
}

func assertHasAllTags(t *testing.T, actual gostatsd.Tags, expected ...string) {
	assert.Equal(t, len(expected), len(actual))
	seenActual := map[string]struct{}{}
	for _, actualTag := range actual {
		seenActual[actualTag] = struct{}{}
	}
	assert.Equal(t, len(actual), len(seenActual), "found duplicates in actual")
	for _, expectedTag := range expected {
		if _, ok := seenActual[expectedTag]; !ok {
			assert.Fail(
				t,
				"missing tag",
				"have tags: [%s], expected tags: [%s], missing tag: %v",
				strings.Join(actual, ","),
				strings.Join(expected, ","),
				expectedTag,
			)
		}
	}

	seenExpected := map[string]struct{}{}
	for _, expectedTag := range expected {
		seenExpected[expectedTag] = struct{}{}
	}
	assert.Equal(t, len(expected), len(seenExpected), "found duplicates in expected")
	for _, actualTag := range actual {
		if _, ok := seenExpected[actualTag]; !ok {
			assert.Fail(
				t,
				"extra tag",
				"have tags: [%s], expected tags: [%s], extra tag: %s",
				strings.Join(actual, ","),
				strings.Join(expected, ","),
				actualTag,
			)
		}
	}
}

func TestTagMetricHandlerAddsNoTags(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	m := &gostatsd.Metric{}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 1, len(tch.m)) // Metric tracked
	assertHasAllTags(t, tch.m[0].Tags)
	assert.Equal(t, "", tch.m[0].Hostname) // No hostname added
}

func TestTagMetricHandlerAddsSingleTag(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{"tag1"}, nil)
	m := &gostatsd.Metric{}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 1, len(tch.m)) // Metric tracked
	assertHasAllTags(t, tch.m[0].Tags, "tag1")
	assert.Equal(t, "", tch.m[0].Hostname) // No hostname added
}

func TestTagMetricHandlerAddsMultipleTags(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{"tag1", "tag2"}, nil)
	m := &gostatsd.Metric{}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 1, len(tch.m)) // Metric tracked
	assertHasAllTags(t, tch.m[0].Tags, "tag1", "tag2")
	assert.Equal(t, "", tch.m[0].Hostname) // No hostname added
}

func TestTagMetricHandlerAddsHostname(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	m := &gostatsd.Metric{
		SourceIP: "1.2.3.4",
	}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 1, len(tch.m))                // Metric tracked
	assert.Equal(t, 0, len(tch.m[0].Tags))        // No tags added
	assert.Equal(t, "1.2.3.4", tch.m[0].Hostname) // Hostname injected
}

func TestTagMetricHandlerAddsDuplicateTags(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{"tag1", "tag2", "tag2", "tag3", "tag1"}, nil)
	m := &gostatsd.Metric{}
	th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	assert.Equal(t, 1, len(tch.m)) // Metric tracked
	assertHasAllTags(t, tch.m[0].Tags, "tag1", "tag2", "tag3")
	assert.Equal(t, "", tch.m[0].Hostname) // No hostname added
}

func TestTagEventHandlerAddsNoTags(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	e := &gostatsd.Event{}
	th.DispatchEvent(context.Background(), e)
	assert.Equal(t, 1, len(tch.e)) // Metric tracked
	assertHasAllTags(t, tch.e[0].Tags)
	assert.Equal(t, "", tch.e[0].Hostname) // No hostname added
}

func TestTagEventHandlerAddsSingleTag(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{"tag1"}, nil)
	e := &gostatsd.Event{}
	th.DispatchEvent(context.Background(), e)
	assert.Equal(t, 1, len(tch.e)) // Metric tracked
	assertHasAllTags(t, tch.e[0].Tags, "tag1")
	assert.Equal(t, "", tch.e[0].Hostname) // No hostname added
}

func TestTagEventHandlerAddsMultipleTags(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{"tag1", "tag2"}, nil)
	e := &gostatsd.Event{}
	th.DispatchEvent(context.Background(), e)
	assert.Equal(t, 1, len(tch.e)) // Metric tracked
	assertHasAllTags(t, tch.e[0].Tags, "tag1", "tag2")
	assert.Equal(t, "", tch.e[0].Hostname) // No hostname added
}

func TestTagEventHandlerAddsHostname(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{}, nil)
	e := &gostatsd.Event{
		SourceIP: "1.2.3.4",
	}
	th.DispatchEvent(context.Background(), e)
	assert.Equal(t, 1, len(tch.e)) // Metric tracked
	assertHasAllTags(t, tch.e[0].Tags)
	assert.Equal(t, "1.2.3.4", tch.e[0].Hostname) // Hostname injected
}

func TestTagEventHandlerAddsDuplicateTags(t *testing.T) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{"tag1", "tag2", "tag2", "tag3", "tag1"}, nil)
	e := &gostatsd.Event{}
	th.DispatchEvent(context.Background(), e)
	assert.Equal(t, 1, len(tch.e)) // Metric tracked
	assertHasAllTags(t, tch.e[0].Tags, "tag1", "tag2", "tag3")
	assert.Equal(t, "", tch.e[0].Hostname) // No hostname added
}

func BenchmarkTagMetricHandlerAddsDuplicateTagsSmall(b *testing.B) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc:cccccccccccccccccccccccccccccccc",
	}, nil)

	b.ReportAllocs()
	b.ResetTimer()

	baseTags := gostatsd.Tags{
		"cccccccccccccccccccccccccccccccc:cccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddd:dddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}

	for n := 0; n < b.N; n++ {
		metricTags := make(gostatsd.Tags, 0, len(baseTags)+th.EstimatedTags())
		metricTags = append(metricTags, baseTags...)
		m := &gostatsd.Metric{
			Tags: metricTags,
		}
		th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	}
}

func BenchmarkTagMetricHandlerAddsDuplicateTagsLarge(b *testing.B) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc:cccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddd:dddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"ffffffffffffffffffffffffffffffff:ffffffffffffffffffffffffffffffff",
		"gggggggggggggggggggggggggggggggg:gggggggggggggggggggggggggggggggg",
		"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh:hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh",
		"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii:iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii",
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj:jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj",
	}, nil)

	b.ReportAllocs()
	b.ResetTimer()

	baseTags := gostatsd.Tags{
		"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh:hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh",
		"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii:iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii",
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj:jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj",
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk:kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk",
		"llllllllllllllllllllllllllllllll:llllllllllllllllllllllllllllllll",
	}

	for n := 0; n < b.N; n++ {
		metricTags := make(gostatsd.Tags, 0, len(baseTags)+th.EstimatedTags())
		metricTags = append(metricTags, baseTags...)
		m := &gostatsd.Metric{
			Tags: metricTags,
		}
		th.DispatchMetrics(context.Background(), []*gostatsd.Metric{m})
	}
}

func BenchmarkTagEventHandlerAddsDuplicateTagsSmall(b *testing.B) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc:cccccccccccccccccccccccccccccccc",
	}, nil)

	eventTags := gostatsd.Tags{
		"cccccccccccccccccccccccccccccccc:cccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddd:dddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		e := &gostatsd.Event{
			Tags: eventTags.Copy(),
		}
		th.DispatchEvent(context.Background(), e)
	}
}

func BenchmarkTagEventHandlerAddsDuplicateTagsLarge(b *testing.B) {
	tch := &capturingHandler{}
	th := NewTagHandler(tch, gostatsd.Tags{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc:cccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddd:dddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"ffffffffffffffffffffffffffffffff:ffffffffffffffffffffffffffffffff",
		"gggggggggggggggggggggggggggggggg:gggggggggggggggggggggggggggggggg",
		"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh:hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh",
		"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii:iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii",
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj:jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj",
	}, nil)

	eventTags := gostatsd.Tags{
		"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh:hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh",
		"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii:iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii",
		"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj:jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj",
		"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk:kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk",
		"llllllllllllllllllllllllllllllll:llllllllllllllllllllllllllllllll",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		e := &gostatsd.Event{
			Tags: eventTags.Copy(),
		}
		th.DispatchEvent(context.Background(), e)
	}
}
