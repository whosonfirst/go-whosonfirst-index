package ioutil

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
)

type ReadSeekCloser struct {
	io.Reader
	io.Seeker
	io.Closer
	reader bool
	closer bool
	seeker bool
	fh     interface{}
	br     *bytes.Reader
	mu     *sync.RWMutex
}

func (rsc *ReadSeekCloser) Read(p []byte) (n int, err error) {

	log.Println("READ")

	if rsc.seeker {
		return rsc.fh.(io.Reader).Read(p)
	}

	br, err := rsc.bytesReader()

	if err != nil {
		return 0, err
	}

	return br.Read(p)
}

func (rsc *ReadSeekCloser) Close() error {

	log.Println("CLOSE")

	if rsc.closer {
		return rsc.fh.(io.ReadCloser).Close()
	}

	return nil
}

func (rsc *ReadSeekCloser) Seek(offset int64, whence int) (int64, error) {

	log.Println("SEEK")

	if rsc.seeker {
		return rsc.fh.(io.Seeker).Seek(offset, whence)
	}

	br, err := rsc.bytesReader()

	if err != nil {
		return 0, err
	}

	return br.Seek(offset, whence)
}

func (rsc *ReadSeekCloser) bytesReader() (*bytes.Reader, error) {

	log.Println("INFLATE")

	rsc.mu.Lock()
	defer rsc.mu.Unlock()

	if rsc.br != nil {
		return rsc.br, nil
	}

	body, err := io.ReadAll(rsc.fh.(io.Reader))

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)
	rsc.br = br

	return br, nil
}

func NewReadSeekCloser(fh interface{}) (io.ReadSeekCloser, error) {

	reader := true
	seeker := false
	closer := false

	switch fh.(type) {
	case io.ReadSeekCloser:
		return fh.(io.ReadSeekCloser), nil
	case io.Reader:
		// pass
	case io.ReadCloser:
		closer = true
	case io.ReadSeeker:
		seeker = true
	default:
		return nil, fmt.Errorf("Invalid or unsupported type")
	}

	mu := new(sync.RWMutex)

	rsc := &ReadSeekCloser{
		reader: reader,
		seeker: seeker,
		closer: closer,
		fh:     fh,
		mu:     mu,
	}

	return rsc, nil

}
