package main

import (
	"context"
	"iter"
	"log"
	"net/http"
	"os"
	"time"

	pu "github.com/takanoriyanagitani/go-postreqs2usock"
)

var spat pu.StreamPath = pu.StreamPath(os.Getenv("ENV_UNIX_SOCK_PATH"))

var timeout time.Duration = time.Second

var preqs pu.PartialReqs = func(
	yield func(pu.HTTPPostRequestBasicPartial, error) bool,
) {
	var ok bool

	ok = yield(
		pu.HTTPPostRequestBasicPartial{
			URLPath: "/api/v1",
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			PostBody: []byte(`{
				"helo":"wrld"
			}`),
		},
		nil,
	)
	if !ok {
		return
	}

	ok = yield(
		pu.HTTPPostRequestBasicPartial{
			URLPath: "/api/v1",
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			PostBody: []byte(`{
				"hello":"world"
			}`),
		},
		nil,
	)

	if !ok {
		return
	}
}

func sub(ctx context.Context) error {
	var hreqs iter.Seq2[*http.Request, error] = preqs.ToReqs(ctx)
	return pu.HTTPRequests(hreqs).ToUnixSocket(
		ctx,
		spat,
		timeout,
	)
}

func main() {
	err := sub(context.Background())
	if nil != err {
		log.Printf("%v\n", err)
	}
}
