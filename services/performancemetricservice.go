package services

import (
	"errors"

	"github.com/juggleim/imserver-console/dbs"
)

var (
	ErrPerformanceParamInvalid      = errors.New("performance param invalid")
	ErrPerformanceMetricTypeInvalid = errors.New("performance metric type invalid")
	ErrPerformanceRangeTooLarge     = errors.New("performance range too large")
)

const performanceMetricsMaxRangeMs int64 = 24 * 60 * 60 * 1000

type PerformanceValueKind string

const (
	PerformanceValueKindPercent PerformanceValueKind = "percent"
	PerformanceValueKindBytes   PerformanceValueKind = "bytes"
	PerformanceValueKindCount   PerformanceValueKind = "count"
)

type PerformanceMetricCatalogItem struct {
	MetricType string               `json:"metric_type"`
	Category   string               `json:"category"`
	ValueKind  PerformanceValueKind `json:"value_kind"`
	LabelKey   string               `json:"label_key"`
}

type PerformanceMetricPoint struct {
	CollectTime int64   `json:"collect_time"`
	MetricValue float64 `json:"metric_value"`
}

type PerformanceMetricResult struct {
	NodeName   string                     `json:"node_name"`
	MetricType string                     `json:"metric_type"`
	Category   string                     `json:"category"`
	ValueKind  PerformanceValueKind       `json:"value_kind"`
	Points     []PerformanceMetricPoint   `json:"points"`
}

type PerformanceNodesResult struct {
	Items []string `json:"items"`
}

type PerformanceCatalogResult struct {
	Items []PerformanceMetricCatalogItem `json:"items"`
}

var performanceMetricCatalog = []PerformanceMetricCatalogItem{
	{MetricType: "cpu.usage_percent", Category: "cpu", ValueKind: PerformanceValueKindPercent, LabelKey: "monitor.metrics.cpu.usage_percent"},
	{MetricType: "load.load1", Category: "cpu", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.cpu.load1"},
	{MetricType: "load.load5", Category: "cpu", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.cpu.load5"},
	{MetricType: "load.load15", Category: "cpu", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.cpu.load15"},

	{MetricType: "memory.total_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.total_bytes"},
	{MetricType: "memory.used_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.used_bytes"},
	{MetricType: "memory.free_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.free_bytes"},
	{MetricType: "memory.available_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.available_bytes"},
	{MetricType: "memory.usage_percent", Category: "memory", ValueKind: PerformanceValueKindPercent, LabelKey: "monitor.metrics.memory.usage_percent"},
	{MetricType: "memory.swap_total_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.swap_total_bytes"},
	{MetricType: "memory.swap_used_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.swap_used_bytes"},
	{MetricType: "memory.swap_free_bytes", Category: "memory", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.memory.swap_free_bytes"},
	{MetricType: "memory.swap_usage_percent", Category: "memory", ValueKind: PerformanceValueKindPercent, LabelKey: "monitor.metrics.memory.swap_usage_percent"},

	{MetricType: "disk.total_bytes", Category: "disk", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.disk.total_bytes"},
	{MetricType: "disk.used_bytes", Category: "disk", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.disk.used_bytes"},
	{MetricType: "disk.free_bytes", Category: "disk", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.disk.free_bytes"},
	{MetricType: "disk.usage_percent", Category: "disk", ValueKind: PerformanceValueKindPercent, LabelKey: "monitor.metrics.disk.usage_percent"},

	{MetricType: "go_runtime.goroutine_count", Category: "go_runtime", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.go_runtime.goroutine_count"},
	{MetricType: "go_runtime.gomaxprocs", Category: "go_runtime", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.go_runtime.gomaxprocs"},
	{MetricType: "go_runtime.cgo_call_count", Category: "go_runtime", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.go_runtime.cgo_call_count"},
	{MetricType: "go_runtime.alloc_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.alloc_bytes"},
	{MetricType: "go_runtime.total_alloc_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.total_alloc_bytes"},
	{MetricType: "go_runtime.sys_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.sys_bytes"},
	{MetricType: "go_runtime.heap_alloc_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.heap_alloc_bytes"},
	{MetricType: "go_runtime.heap_sys_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.heap_sys_bytes"},
	{MetricType: "go_runtime.heap_inuse_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.heap_inuse_bytes"},
	{MetricType: "go_runtime.stack_inuse_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.stack_inuse_bytes"},
	{MetricType: "go_runtime.next_gc_bytes", Category: "go_runtime", ValueKind: PerformanceValueKindBytes, LabelKey: "monitor.metrics.go_runtime.next_gc_bytes"},
	{MetricType: "go_runtime.last_gc_time_unix_nano", Category: "go_runtime", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.go_runtime.last_gc_time_unix_nano"},
	{MetricType: "go_runtime.num_gc", Category: "go_runtime", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.go_runtime.num_gc"},
	{MetricType: "go_runtime.pause_total_ns", Category: "go_runtime", ValueKind: PerformanceValueKindCount, LabelKey: "monitor.metrics.go_runtime.pause_total_ns"},
}

var performanceMetricCatalogMap map[string]PerformanceMetricCatalogItem

func init() {
	performanceMetricCatalogMap = make(map[string]PerformanceMetricCatalogItem, len(performanceMetricCatalog))
	for _, item := range performanceMetricCatalog {
		performanceMetricCatalogMap[item.MetricType] = item
	}
}

func PerformanceMetricsMaxRangeMs() int64 {
	return performanceMetricsMaxRangeMs
}

func LookupPerformanceMetricCatalog(metricType string) (PerformanceMetricCatalogItem, bool) {
	item, ok := performanceMetricCatalogMap[metricType]
	return item, ok
}

func QryPerformanceNodes() *PerformanceNodesResult {
	dao := dbs.PerformanceMetricDao{}
	return &PerformanceNodesResult{Items: dao.ListNodeNames()}
}

func QryPerformanceCatalog() *PerformanceCatalogResult {
	items := make([]PerformanceMetricCatalogItem, len(performanceMetricCatalog))
	copy(items, performanceMetricCatalog)
	return &PerformanceCatalogResult{Items: items}
}

func QryPerformanceMetric(nodeName, metricType string, start, end int64) (*PerformanceMetricResult, error) {
	if nodeName == "" || metricType == "" {
		return nil, ErrPerformanceParamInvalid
	}
	catalogItem, ok := LookupPerformanceMetricCatalog(metricType)
	if !ok {
		return nil, ErrPerformanceMetricTypeInvalid
	}
	if end <= start {
		return nil, ErrPerformanceParamInvalid
	}
	if end-start > performanceMetricsMaxRangeMs {
		return nil, ErrPerformanceRangeTooLarge
	}

	dao := dbs.PerformanceMetricDao{}
	rows := dao.QryMetric(nodeName, metricType, start, end)
	points := make([]PerformanceMetricPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, PerformanceMetricPoint{
			CollectTime: row.CollectTime,
			MetricValue: row.MetricValue,
		})
	}
	return &PerformanceMetricResult{
		NodeName:   nodeName,
		MetricType: metricType,
		Category:   catalogItem.Category,
		ValueKind:  catalogItem.ValueKind,
		Points:     points,
	}, nil
}
