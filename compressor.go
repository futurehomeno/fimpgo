package fimpgo

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
)

type MsgCompressor struct {
	compressor        *gzip.Writer
	compressionBuffer bytes.Buffer
	mux               sync.Mutex
}

func NewMsgCompressor(alg, compLevel string) *MsgCompressor {
	var err error
	comp := &MsgCompressor{}
	comp.compressor, err = gzip.NewWriterLevel(&comp.compressionBuffer, gzip.BestCompression)
	if err != nil {
		log.Error("Compressor can't be initiated .Err:", err)
	}

	return comp
}

// CompressBinMsg - compresses binary message and return compressed byte array.
func (c *MsgCompressor) CompressBinMsg(msg []byte) ([]byte, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.compressor.Reset(&c.compressionBuffer)

	_, err := c.compressor.Write(msg)
	if err != nil {
		return nil, err
	}

	if err := c.compressor.Flush(); err != nil {
		return nil, err
	}

	if err := c.compressor.Close(); err != nil {
		return nil, err
	}

	cp := c.compressionBuffer.Bytes()
	c.compressionBuffer.Reset()
	return cp, nil
}

func (c *MsgCompressor) DecompressBinMsg(binMsg []byte) ([]byte, error) {
	var err error
	var decompressorBuffer bytes.Buffer
	decompressorBuffer.Write(binMsg)

	decompressor, err := gzip.NewReader(&decompressorBuffer)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(decompressor)
	if err := decompressor.Close(); err != nil {
		return nil, err
	}

	decompressorBuffer.Reset()
	return response, err
}

func (c *MsgCompressor) CompressFimpMsg(msg *FimpMessage) ([]byte, error) {
	binMsg, err := msg.SerializeToJson()
	if err != nil {
		return nil, err
	}

	return c.CompressBinMsg(binMsg)
}

func (c *MsgCompressor) DecompressFimpMsg(compBinMsg []byte) (*FimpMessage, error) {
	binMsg, err := c.DecompressBinMsg(compBinMsg)
	if err != nil {
		return nil, err
	}

	return NewMessageFromBytes(binMsg)
}
