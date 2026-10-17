package fimpgo

import (
	"testing"
	"time"

	"github.com/futurehomeno/fimpgo/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgCompressor_Errors(t *testing.T) {
	c := NewMsgCompressor("", "")

	_, err := c.DecompressBinMsg([]byte("not gzip"))
	assert.Error(t, err)

	good, err := c.CompressBinMsg([]byte("hello hello hello"))
	require.NoError(t, err)

	_, err = c.DecompressBinMsg(good[:len(good)-6])
	assert.Error(t, err)

	_, err = c.DecompressFimpMsg([]byte("not gzip"))
	assert.Error(t, err)

	notFimp, err := c.CompressBinMsg([]byte("plain text"))
	require.NoError(t, err)
	_, err = c.DecompressFimpMsg(notFimp)
	assert.Error(t, err)

	_, err = c.CompressFimpMsg(NewObjectMessage("evt.a.b", "s", make(chan int), nil, nil, nil))
	assert.Error(t, err)
}

func TestConnStateT(t *testing.T) {
	var c ConnStateT

	c.OnDone()
	assert.False(t, c.IsConnected())

	c.Init()
	assert.ErrorIs(t, c.WaitConnected(time.Millisecond), utils.ErrTimeout)

	c.OnConnect()
	assert.True(t, c.IsConnected())
	assert.NoError(t, c.WaitConnected(time.Second))

	c.OnDone()
	c.OnDone()
	assert.False(t, c.IsConnected())
	assert.ErrorIs(t, c.WaitConnected(time.Second), utils.ErrConnectionLost)

	select {
	case <-c.DoneC():
	default:
		t.Fatal("done channel not closed")
	}
}
