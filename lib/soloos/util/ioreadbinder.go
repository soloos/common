package util

import "io"

type IOReadBinder struct {
	readerIndex int
	readers     []io.ReadCloser
}

func (p *IOReadBinder) BindReadCloser(readCloser io.ReadCloser) {
	p.readers = append(p.readers, readCloser)
}

func (p *IOReadBinder) Read(p []byte) (n int, err error) {
	for {
		n, err = p.readers[p.readerIndex].Read(p)
		if err == nil {
			return
		}

		if p.readerIndex >= len(p.readers) {
			return
		}

		p.readerIndex++
	}
}

func (p *IOReadBinder) Close() error {
	var err error
	for k, _ := range p.readers {
		err = p.readers[k].Close()
		if err != nil {
			return err
		}
	}
	return nil
}
