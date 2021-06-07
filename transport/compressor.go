package transport

import (
	"bytes"
	"compress/gzip"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"sync"
)


type MsgCompressor struct {
	compressor        *gzip.Writer
	decompressor      *gzip.Reader
	compressionBuffer bytes.Buffer
	decompressorBuffer bytes.Buffer
	mux                sync.Mutex
}

func NewMsgCompressor(alg,compLevel string ) *MsgCompressor {
	var err error
	comp := &MsgCompressor{}
	comp.compressor,err = gzip.NewWriterLevel(&comp.compressionBuffer,gzip.BestCompression)
	if err != nil {
		log.Error("Compressor can't be initiated .Err:",err)
	}

	return comp
}

//CompressBinMsg - compresses binary message and return compressed byte array.
func (c *MsgCompressor) CompressBinMsg(msg []byte) ([]byte, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.compressor.Reset(&c.compressionBuffer)
	_ , err := c.compressor.Write(msg)
	if err != nil {
		log.Error("Compression error :",err.Error())
		return nil,err
	}
	c.compressor.Flush()
	c.compressor.Close()
	cp := c.compressionBuffer.Bytes()
	c.compressionBuffer.Reset()
	return cp , nil
}

func (c *MsgCompressor) DecompressBinMsg(binMsg []byte) ([]byte, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.decompressor == nil {
		var err error
		c.decompressorBuffer.Write(binMsg)
		c.decompressor,err = gzip.NewReader(&c.decompressorBuffer)
		if err != nil {
			log.Error("Decompressor can't be initiated .Err:",err)
			return nil, err
		}
	}else {
		c.decompressorBuffer.Reset()
		c.decompressorBuffer.Write(binMsg)
	}
	var resB bytes.Buffer
	_, err := resB.ReadFrom(c.decompressor)
	if err != nil {
		log.Error("Decompression error .Err:",err.Error())
		return nil, err
	}
	return resB.Bytes(),nil
}

func (c *MsgCompressor) CompressFimpMsg(msg *fimpgo.FimpMessage) ([]byte, error) {
	binMsg,err := msg.SerializeToJson()
	if err != nil {
		return nil, err
	}
	return c.CompressBinMsg(binMsg)
}

func (c MsgCompressor) DecompressFimpMsg(compBinMsg []byte) (*fimpgo.FimpMessage, error) {
	binMsg,err := c.DecompressBinMsg(compBinMsg)
	if err != nil {
		return nil, err
	}
	return fimpgo.NewMessageFromBytes(binMsg)
}




