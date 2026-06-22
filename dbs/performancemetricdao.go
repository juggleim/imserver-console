package dbs

import "github.com/juggleim/imserver-console/commons/dbcommons"

type PerformanceMetricDao struct {
	ID          int64   `gorm:"column:id;primaryKey"`
	NodeName    string  `gorm:"column:node_name"`
	CollectTime int64   `gorm:"column:collect_time"`
	MetricType  string  `gorm:"column:metric_type"`
	MetricValue float64 `gorm:"column:metric_value"`
}

func (metric PerformanceMetricDao) TableName() string {
	return "performance_metrics"
}

func (metric PerformanceMetricDao) ListNodeNames() []string {
	var items []string
	err := dbcommons.GetDb().Model(&PerformanceMetricDao{}).
		Distinct("node_name").
		Order("node_name asc").
		Pluck("node_name", &items).Error
	if err == nil {
		return items
	}
	return []string{}
}

func (metric PerformanceMetricDao) QryMetric(nodeName, metricType string, start, end int64) []*PerformanceMetricDao {
	var items []*PerformanceMetricDao
	err := dbcommons.GetDb().
		Where("node_name=? and metric_type=? and collect_time>=? and collect_time<=?", nodeName, metricType, start, end).
		Order("collect_time asc").
		Limit(2000).
		Find(&items).Error
	if err == nil {
		return items
	}
	return []*PerformanceMetricDao{}
}
