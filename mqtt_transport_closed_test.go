// The race detector reports a close concurrent with a send, which is the case under test.
//go:build !race

package fimpgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMqttTransport_ChannelClosedDuringWaitIsRecovered(t *testing.T) {
	mqtt := NewMqttTransport(covBroker, "cov_closed", "", "", true, 1, 1, nil)
	mqtt.SetReceiveChTimeout(5)
	payload, err := NewStringMessage("evt.a.b", "svc", "v", nil, nil, nil).SerializeToJson()
	require.NoError(t, err)

	closed := make(MessageCh)
	mqtt.RegisterChannel("closed", closed)
	time.AfterFunc(100*time.Millisecond, func() { close(closed) })

	mqtt.handleIncomingMessage(fakeMessage{topic: "pt:j1/mt:evt/rt:app/rn:n/ad:1", payload: payload})

	mqtt.UnregisterChannel("closed")
}
