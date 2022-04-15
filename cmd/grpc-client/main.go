package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/alexey-mavrin/go-musthave-devops/internal/common"
	"github.com/alexey-mavrin/go-musthave-devops/internal/grpcint"
	pb "github.com/alexey-mavrin/go-musthave-devops/internal/grpcint/proto"
	"google.golang.org/grpc"
)

type cfg struct {
	m      common.Metrics
	server string
}

func parseCmdLine() (cfg, error) {
	var c cfg
	gauge := flag.Float64("g", 0, "gauge value")
	counter := flag.Int64("c", 0, "counter value")
	name := flag.String("n", "", "name")
	server := flag.String("s", ":3200", "name")
	flag.Parse()

	isG := common.IsFlagPassed("g")
	isC := common.IsFlagPassed("c")
	isN := common.IsFlagPassed("n")

	if (isG && isC) || (!isG && !isC) {
		return c, errors.New("set either 'g' or 'c' value")
	}

	if !isN {
		return c, errors.New("name is not set")
	}

	c.m.ID = *name
	if isG {
		c.m.MType = common.NameGauge
		c.m.Value = gauge
	}

	if isC {
		c.m.MType = common.NameCounter
		c.m.Delta = counter
	}
	c.server = *server

	return c, nil
}

func main() {
	config, err := parseCmdLine()
	if err != nil {
		log.Fatal(err)
	}

	p := grpcint.MetricsToPb(config.m)
	pList := make([](*pb.Metrics), 0)
	pList = append(pList, p)

	conn, err := grpc.Dial(config.server, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mc := pb.NewMetricesClient(conn)

	req := pb.UpdateMetricesRequest{
		Count:    1,
		Metrices: pList,
	}
	resp, _ := mc.UpdateMetrices(context.Background(), &req)
	if resp.Error != "" {
		log.Fatal(resp.Error)
	}
}
