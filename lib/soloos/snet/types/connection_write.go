package types

import (
	"io"
)

func (p *Connection) WriteAcquire() {
	p.writeMutex.Lock()
}

func (p *Connection) WriteRelease() {
	p.writeMutex.Unlock()
}

func (p *Connection) afterWriteHeader(bodySize uint32) error {
	p.LastWriteLimit = bodySize
	if p.LastWriteLimit == 0 {
		p.WriteRelease()
	}
	return nil
}

func (p *Connection) Write(b []byte) (int, error) {
	if p.LastWriteLimit == 0 {
		return 0, io.EOF
	}

	if len(b) == 0 {
		return 0, nil
	}

	var (
		n   int
		err error
	)

	if p.LastWriteLimit > uint32(len(b)) {
		n, err = p.NetConn.Write(b)
	} else {
		n, err = p.NetConn.Write(b[:int(p.LastWriteLimit)])
	}

	if err != nil {
		p.WriteRelease()
		return n, err
	}

	p.LastWriteLimit -= uint32(n)
	if p.LastWriteLimit == 0 {
		p.WriteRelease()
	}

	return n, nil
}

func (p *Connection) WriteAll(b []byte) error {
	var (
		n   int
		err error
	)

	for off := 0; off < len(b); off += n {
		n, err = p.Write(b[off:])
		if err != nil {
			return err
		}
	}

	return nil
}
