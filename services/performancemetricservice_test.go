package services

import "testing"

func TestLookupPerformanceMetricCatalog(t *testing.T) {
	item, ok := LookupPerformanceMetricCatalog("cpu.usage_percent")
	if !ok || item.Category != "cpu" || item.ValueKind != PerformanceValueKindPercent {
		t.Fatalf("unexpected catalog item: %+v ok=%v", item, ok)
	}
	_, ok = LookupPerformanceMetricCatalog("unknown.metric")
	if ok {
		t.Fatal("expected unknown metric to be missing")
	}
}

func TestQryPerformanceMetricValidation(t *testing.T) {
	_, err := QryPerformanceMetric("", "cpu.usage_percent", 1, 2)
	if err != ErrPerformanceParamInvalid {
		t.Fatalf("expected param invalid, got %v", err)
	}
	_, err = QryPerformanceMetric("node-a", "", 1, 2)
	if err != ErrPerformanceParamInvalid {
		t.Fatalf("expected param invalid, got %v", err)
	}
	_, err = QryPerformanceMetric("node-a", "bad.metric", 1, 2)
	if err != ErrPerformanceMetricTypeInvalid {
		t.Fatalf("expected metric type invalid, got %v", err)
	}
	_, err = QryPerformanceMetric("node-a", "cpu.usage_percent", 10, 10)
	if err != ErrPerformanceParamInvalid {
		t.Fatalf("expected param invalid for equal range, got %v", err)
	}
	_, err = QryPerformanceMetric("node-a", "cpu.usage_percent", 0, PerformanceMetricsMaxRangeMs()+1)
	if err != ErrPerformanceRangeTooLarge {
		t.Fatalf("expected range too large, got %v", err)
	}
}

func TestQryPerformanceCatalogCount(t *testing.T) {
	items := QryPerformanceCatalog().Items
	if len(items) != 31 {
		t.Fatalf("expected 31 catalog items, got %d", len(items))
	}
}
