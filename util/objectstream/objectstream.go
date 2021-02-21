package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "htttp://"+server+"/objects/"+object, reader)
		client := http.Client{}
		r, err := client.Do(request)
		if err != nil && r.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return htp code: %d", r.StatusCode)
		}
		c <- err
	}()
	return &PutStream{
		writer: writer,
		c:      c,
	}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}

type GetStream struct {
	reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
	reader, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if reader.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", reader.StatusCode)
	}
	return &GetStream{reader: reader.Body}, nil
}
func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	return newGetStream("http://" + server + "/objects/" + object)
}
func (read *GetStream) Read(p []byte) (n int, err error) {
	return read.reader.Read(p)
}
