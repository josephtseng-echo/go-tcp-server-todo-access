package packet

const (
	PACKET_LIMIT = 65535
	PACKET_POOL  = 10000
)

var (
	_pool = make(chan *Packet, PACKET_POOL)
)

type Packet struct {
	pos  int
	buffer []byte
}

func init() {
	go func() {
		for {
			_pool <- &Packet{buffer: make([]byte, 0, 512)}
		}
	}()
}

func (p *Packet) Buffer() []byte {
	return p.buffer
}

func (p *Packet) Length() int {
	return len(p.buffer)
}

func (p *Packet) Pack_Begin(cmd uint16) {
	p.WriteByte(byte('E'))
	p.WriteByte(byte('S'))
	p.WriteUint16(cmd)
	p.WriteUint16(1)
	p.WriteUint16(0)
}

func (p *Packet) Pack_End() {
	size := len(p.buffer) - 8
	p.buffer[6] = byte(size >> 8)
	p.buffer[7] = byte(size)
}

func (p *Packet) Pack_Begin_Rpc() {
	p.WriteByte(byte('R'))
	p.WriteByte(byte('P'))
	p.WriteByte(byte('C'))
	p.WriteUint8(5)
}

func (p *Packet) Pack_End_Rpc() {
	size := len(p.buffer) - 6
	p.buffer[4] = byte(size >> 8)
	p.buffer[5] = byte(size)
}