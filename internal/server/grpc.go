package server

import (
	"context"
	"fmt"

	pb "github.com/alexey-mavrin/go-musthave-devops/internal/grpcint/proto"
)

// MetricesServer is to serve grps requests
type MetricesServer struct {
	pb.UnimplementedMetricesServer
}

func pbToStatReq(p *pb.Metrics) statReq {
	var req statReq
	req.name = p.Id
	switch p.Mtype {
	case pb.Metrics_GAUGE:
		req.statType = statTypeGauge
		req.valueGauge = p.Value
	case pb.Metrics_COUNTER:
		req.statType = statTypeCounter
		req.valueCounter = p.Delta
	}
	return req
}

// UpdateMetrices get the sequence of mertices and store them in the server
func (s *MetricesServer) UpdateMetrices(
	ctx context.Context,
	in *pb.UpdateMetricesRequest,
) (*pb.UpdateMetricesResponse, error) {
	var ret pb.UpdateMetricesResponse
	for i := int32(0); i < in.Count; i++ {
		err := updateStatStorage(pbToStatReq(in.Metrices[i]))
		if err != nil {
			ret.Error = fmt.Sprintf("%v", err)
			return &ret, err
		}
	}
	return &ret, nil
}
