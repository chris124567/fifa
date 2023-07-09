package api

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"

	http "github.com/useflyent/fhttp"

	"github.com/andybalholm/brotli"
)

type compressedReader struct {
	compressed, original io.ReadCloser
}

type newReader func(io.Reader) (io.ReadCloser, error)

func newGzipReader(r io.Reader) (io.ReadCloser, error) {
	compressed, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return compressed, err
}

func newDeflateReader(r io.Reader) (io.ReadCloser, error) {
	return flate.NewReader(r), nil
}

func newBrotliReader(r io.Reader) (io.ReadCloser, error) {
	return io.NopCloser(brotli.NewReader(r)), nil
}

// allows us to close the original reader and the compressed reader
func newCompressedReader(original io.ReadCloser, newReaderFunc newReader) (compressedReader, error) {
	compressed, err := newReaderFunc(original)
	if err != nil {
		return compressedReader{}, err
	}
	return compressedReader{compressed: compressed, original: original}, nil
}

func (r compressedReader) Read(p []byte) (n int, err error) {
	return r.compressed.Read(p)
}

func (r compressedReader) Close() error {
	if err := r.compressed.Close(); err != nil {
		return err
	}
	return r.original.Close()
}

type transportNoReferer struct {
	original http.RoundTripper
}

func (t transportNoReferer) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Del(http.CanonicalHeaderKey("Referer"))
	return t.original.RoundTrip(request)
}

func noReferer(c http.Client) http.Client {
	c.Transport = transportNoReferer{c.Transport}
	return c
}

func noRedirect(c http.Client) http.Client {
	c.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return c
}

func closeBody(body io.ReadCloser) {
	io.Copy(ioutil.Discard, body)
	body.Close()
}
