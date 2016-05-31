package main

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type StatsType int

const (
	Consume StatsType = iota
	ConsumeFail
	Publish
	PublishFail
	SlowConsumerAlert
)

// Stats stores various stats infomation
type Stats struct {
	Consume       uint64 `json:"consume"`
	ConsumePerSec uint64 `json:"consume_per_sec"`
	ConsumeFail   uint64 `json:"consume_fail"`

	Publish       uint64 `json:"publish"`
	PublishPerSec uint64 `json:"publish_per_sec"`
	PublishFail   uint64 `json:"publish_fail"`

	SlowConsumerAlert uint64 `json:"slow_consumer_alert"`

	// Delay is Consume - Pulish
	// This indicate how slow publish to kafka
	Delay uint64 `json:"delay"`
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Json() ([]byte, error) {
	s.Delay = s.Consume - s.Publish
	return json.Marshal(s)
}

func (s *Stats) PerSec() {
	lastConsume, lastPublish := uint64(0), uint64(0)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			s.ConsumePerSec = s.Consume - lastConsume
			s.PublishPerSec = s.Publish - lastPublish

			lastConsume = s.Consume
			lastPublish = s.Publish
		}
	}
}

func (s *Stats) Inc(statsType StatsType) {
	switch statsType {
	case Consume:
		atomic.AddUint64(&s.Consume, 1)
	case ConsumeFail:
		atomic.AddUint64(&s.ConsumeFail, 1)
	case Publish:
		atomic.AddUint64(&s.Publish, 1)
	case PublishFail:
		atomic.AddUint64(&s.PublishFail, 1)
	case SlowConsumerAlert:
		atomic.AddUint64(&s.SlowConsumerAlert, 1)
	}
}