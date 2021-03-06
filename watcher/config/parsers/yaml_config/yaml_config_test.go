package yamlconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/watcher"
)

type mockLoader struct {
	Source []byte
}

func (l *mockLoader) Load() ([]byte, error) {
	return l.Source, nil
}

func TestParse(t *testing.T) {
	loader := &mockLoader{
		Source: []byte(`---
- host: "amqp://guest:guest@localhost"
  exchange: "/"
  routing_key: "#"
- host: "amqp://guest:guest@localhost"
  exchange: "/other-exchange"
  exchange_type: "direct"
  routing_key: "key.*"`),
	}
	parser := New(loader)
	configs, err := parser.Parse()

	assert.Nil(t, err)
	expected := []watcher.Config{
		watcher.Config{
			Host:         "amqp://guest:guest@localhost",
			Exchange:     "/",
			ExchangeType: "topic",
			RoutingKey:   "#",
		},
		watcher.Config{
			Host:         "amqp://guest:guest@localhost",
			Exchange:     "/other-exchange",
			ExchangeType: "direct",
			RoutingKey:   "key.*",
		},
	}
	assert.Equal(t, expected, configs)
}

func TestParse_InvalidConfig(t *testing.T) {
	loader := &mockLoader{
		Source: []byte(`---
- host: "amqp://guest:guest@localhost"
  routing_key: "#"`),
	}
	parser := New(loader)
	_, err := parser.Parse()

	assert.NotNil(t, err)
}
