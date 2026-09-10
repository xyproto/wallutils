package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
)

type bitWriter struct {
    Buffer          *bytes.Buffer
    BitBuffer       uint64
    BitBufferSize   int
}

func (w *bitWriter) writeBits(value uint64, n int) {
    if n < 0 || n > 64 {
        panic("Invalid bit count: must be between 1 and 64")
    }

    if value >= (1 << n) {
        panic("too many bits for the given value")
    }
    
    w.BitBuffer |= (value << w.BitBufferSize)
    w.BitBufferSize += n
    w.writeThrough()
}

func (w *bitWriter) writeBytes(values []byte) {
    for _, v := range values {
        w.writeBits(uint64(v), 8)
    }
}

func (w *bitWriter) writeCode(code huffmanCode) {
    if code.Depth <= 0 {
        return
    }

    value := uint64(code.Bits)
    reversed := uint64(0)
    for i := 0; i < code.Depth; i++ {
        reversed = (reversed << 1) | (value & 1)
        value >>= 1
    }

    w.writeBits(reversed, code.Depth)
}

func (w *bitWriter) alignByte() {
    w.BitBufferSize = (w.BitBufferSize + 7) &^ 7
    w.writeThrough()
}

func (w *bitWriter) writeThrough() {
    for w.BitBufferSize >= 8 {
        w.Buffer.WriteByte(byte(w.BitBuffer & 0xFF))
        w.BitBuffer >>= 8
        w.BitBufferSize -= 8
    }
}