package packet

import (
    "errors"
	"math"
)

func (p *Packet) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
}

func (p *Packet) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.buffer) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.buffer[p.pos]
	p.pos++
	return
}

func (p *Packet) ReadBytes() (ret []byte, err error) {
	if p.pos+4 > len(p.buffer) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadUint32()
	if p.pos+int(size) > len(p.buffer) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = p.buffer[p.pos : p.pos+int(size)]
	p.pos += int(size)
	return
}

func (p *Packet) ReadString() (ret string, err error) {
	if p.pos+4 > len(p.buffer) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadUint32()
	if p.pos+int(size) > len(p.buffer) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.buffer[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}

func (p *Packet) ReadUint8() (ret uint8, err error) {
	if p.pos+1 > len(p.buffer) {
		err = errors.New("read uint8 failed")
		return
	}
	buf := p.buffer[p.pos : p.pos+1]
	ret = uint8(buf[0])
	p.pos += 1
	return
}

func (p *Packet) ReadInt8() (ret int8, err error) {
	_ret, _err := p.ReadUint8()
	ret = int8(_ret)
	err = _err
	return
}

func (p *Packet) ReadUint16() (ret uint16, err error) {
	if p.pos+2 > len(p.buffer) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.buffer[p.pos : p.pos+2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.pos += 2
	return
}

func (p *Packet) ReadInt16() (ret int16, err error) {
	_ret, _err := p.ReadUint16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *Packet) ReadUint32() (ret uint32, err error) {
	if p.pos+4 > len(p.buffer) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.buffer[p.pos : p.pos+4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	p.pos += 4
	return
}

func (p *Packet) ReadInt32() (ret int32, err error) {
	_ret, _err := p.ReadUint32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *Packet) ReadUint64() (ret uint64, err error) {
	if p.pos+8 > len(p.buffer) {
		err = errors.New("read uint64 failed")
		return
	}

	ret = 0
	buf := p.buffer[p.pos : p.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}
	p.pos += 8
	return
}

func (p *Packet) ReadInt64() (ret int64, err error) {
	_ret, _err := p.ReadUint64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *Packet) ReadFloat32() (ret float32, err error) {
	bits, _err := p.ReadUint32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *Packet) ReadFloat64() (ret float64, err error) {
	bits, _err := p.ReadUint64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}

func Reader(buffer []byte) *Packet {
	return &Packet{buffer: buffer}
}