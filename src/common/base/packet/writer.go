package packet

import (
	"math"
)


func (p *Packet) WriteZeros(n int) {
	for i := 0; i < n; i++ {
		p.buffer = append(p.buffer, byte(0))
	}
}

func (p *Packet) WriteBool(v bool) {
	if v {
		p.buffer = append(p.buffer, byte(1))
	} else {
		p.buffer = append(p.buffer, byte(0))
	}
}

func (p *Packet) WriteByte(v byte) {
	p.buffer = append(p.buffer, v)
}

func (p *Packet) WriteBytes(v []byte) {
	p.WriteUint32(uint32(len(v)))
	p.buffer = append(p.buffer, v...)
}

func (p *Packet) WriteRawBytes(v []byte) {
	p.buffer = append(p.buffer, v...)
}

func (p *Packet) WriteString(v string) {
	bytes := []byte(v)
	p.WriteUint32(uint32(len(bytes) + 1))
	p.buffer = append(p.buffer, bytes...)
	p.buffer = append(p.buffer, byte(0))
}

func (p *Packet) WriteUint8(v uint8) {
	p.buffer = append(p.buffer, byte(v), byte(v))
}

func (p *Packet) WriteInt8(v int8) {
    p.WriteUint8(uint8(v))
}

func (p *Packet) WriteUint16(v uint16) {
	p.buffer = append(p.buffer, byte(v>>8), byte(v))
}

func (p *Packet) WriteInt16(v int16) {
	p.WriteUint16(uint16(v))
}

func (p *Packet) WriteUint32(v uint32) {
	p.buffer = append(p.buffer, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (p *Packet) WriteInt32(v int32) {
	p.WriteUint32(uint32(v))
}

func (p *Packet) WriteUint64(v uint64) {
	p.buffer = append(p.buffer, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (p *Packet) WriteInt64(v int64) {
	p.WriteUint64(uint64(v))
}

func (p *Packet) WriteFloat32(f float32) {
	v := math.Float32bits(f)
	p.WriteUint32(v)
}

func (p *Packet) WriteFloat64(f float64) {
	v := math.Float64bits(f)
	p.WriteUint64(v)
}

func Writer() *Packet {
	return <-_pool
}