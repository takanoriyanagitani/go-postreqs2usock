package preqs2usock

import (
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"net"
	"net/http"
	"time"
)

type URLPath string

func (p URLPath) AsString() string { return string(p) }

type PostBody []byte

func (b PostBody) AsReader() io.Reader { return bytes.NewReader(b) }

type HTTPPostRequestBasicPartial struct {
	URLPath
	http.Header
	PostBody
}

func (p HTTPPostRequestBasicPartial) ToRequest(
	ctx context.Context,
) (*http.Request, error) {
	var url string = "http://localhost" + p.URLPath.AsString()
	req, err := http.
		NewRequestWithContext(ctx, http.MethodPost, url, p.PostBody.AsReader())
	if nil != err {
		return nil, err
	}

	req.Header = p.Header
	return req, nil
}

type PartialReqs iter.Seq2[HTTPPostRequestBasicPartial, error]

func (p PartialReqs) ToReqs(
	ctx context.Context,
) iter.Seq2[*http.Request, error] {
	return func(yield func(*http.Request, error) bool) {
		for preq, err := range p {
			if nil != err {
				yield(nil, err)
				return
			}

			hreq, err := preq.ToRequest(ctx)
			if !yield(hreq, err) {
				return
			}
		}
	}
}

type HTTPRequests iter.Seq2[*http.Request, error]

type StreamPath string

func (p StreamPath) AsString() string { return string(p) }

func (r HTTPRequests) ToUnixSocket(
	ctx context.Context,
	spat StreamPath,
	timeout time.Duration,
) error {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", spat.AsString())
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	for hreq, err := range r {
		if nil != err {
			return err
		}

		resp, err := client.Do(hreq)
		if nil != err {
			return err
		}
		_, err = io.Copy(io.Discard, resp.Body)
		allErr := errors.Join(err, resp.Body.Close())
		if nil != allErr {
			return allErr
		}
	}
	return nil
}
