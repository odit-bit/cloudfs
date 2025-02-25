package rpc

import (
	"sync/atomic"
	"time"
)

var grpcMetric = metric{}

type metric struct {
	TotalReceive atomic.Int64
	TotalSend    atomic.Int64
}

type stat struct {
	Receive   int64
	Send      int64
	Timestamp int64
}

func (m *metric) status() *stat {
	return &stat{
		Receive:   m.TotalReceive.Load(),
		Send:      m.TotalSend.Load(),
		Timestamp: time.Now().UTC().Unix(),
	}
}
