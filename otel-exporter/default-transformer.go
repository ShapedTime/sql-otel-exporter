package otel_exporter

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

type DefaultTransformer struct {
	name             string
	valueColumn      string
	attributeColumns []string
	staticAttributes map[string]string
}

func NewDefaultTransformer(name, valueColumn string, attributeColumns []string, staticAttributes map[string]string) *DefaultTransformer {
	return &DefaultTransformer{
		name:             name,
		valueColumn:      valueColumn,
		attributeColumns: attributeColumns,
		staticAttributes: staticAttributes,
	}
}

func (t *DefaultTransformer) Transform(data []map[string]string) ([]Metric, error) {
	metrics := make([]Metric, 0, len(data))

	for _, row := range data {
		val, err := strconv.ParseFloat(row[t.valueColumn], 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse row %s with value '%s' to float: %w", t.valueColumn, row[t.valueColumn], err)
		}

		metrics = append(metrics, Metric{
			name:  t.name,
			value: val,
			attributes: func() []attribute.KeyValue {
				attrs := make([]attribute.KeyValue, 0, len(t.attributeColumns)+len(t.staticAttributes))

				for _, col := range t.attributeColumns {
					attrs = append(attrs, attribute.String(col, row[col]))
				}

				for k, v := range t.staticAttributes {
					attrs = append(attrs, attribute.String(k, v))
				}

				return attrs
			}(),
		})
	}

	return metrics, nil
}
