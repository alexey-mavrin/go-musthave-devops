package grpcint

import (
	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	pb "github.com/alexey-mavrin/go-musthave-devops/internal/grpcint/proto"
)

// MetricsToPb converts common.Metrics to pb.Metrics
// Fow now it ignores the possibility of errors in passed structure
func MetricsToPb(m common.Metrics) *pb.Metrics {
	p := pb.Metrics{}

	p.Id = m.ID

	if m.MType == common.NameGauge {
		p.Mtype = pb.Metrics_GAUGE
		p.Value = *m.Value
	}
	if m.MType == common.NameCounter {
		p.Mtype = pb.Metrics_COUNTER
		p.Delta = *m.Delta
	}
	return &p
}
