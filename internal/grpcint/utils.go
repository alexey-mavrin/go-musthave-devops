package grpcint

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

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

// ToString converts pb.Metrics to string
func ToString(p *pb.Metrics) string {

	str := p.Id

	switch p.Mtype {
	case pb.Metrics_GAUGE:
		str += fmt.Sprintf(":%s:%f",
			common.NameGauge,
			p.Value,
		)
	case pb.Metrics_COUNTER:
		str += fmt.Sprintf(":%s:%d",
			common.NameCounter,
			p.Delta,
		)
	}

	return str
}

// ComputeHash computes the hash for the given key and pb.Metrics
func ComputeHash(p *pb.Metrics, key string) (*[]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("no key")
	}
	if p.Id == "" {
		return nil, fmt.Errorf("empty ID field")
	}

	h := hmac.New(sha256.New, []byte(key))
	metricsStr := ToString(p)
	h.Write([]byte(metricsStr))
	hash := h.Sum(nil)

	return &hash, nil
}

// StoreHash stores hash to pb.Metrics struct
func StoreHash(p *pb.Metrics, key string) error {
	if key == "" {
		return nil
	}
	h, err := ComputeHash(p, key)
	if err != nil {
		return err
	}
	p.Hash = hex.EncodeToString(*h)
	return nil
}

// CheckHash checks hash stored into Metrics struct
func CheckHash(p *pb.Metrics, key string) error {
	if key == "" {
		return nil
	}
	h, err := ComputeHash(p, key)
	if err != nil {
		return err
	}
	hashStr := hex.EncodeToString(*h)
	if p.Hash != hashStr {
		return fmt.Errorf("hash value incorrect")
	}
	return nil
}
