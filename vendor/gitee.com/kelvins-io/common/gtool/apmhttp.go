package gtool

import (
	"context"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"net/http"
	"time"
)

var (
	netClient = &http.Client{
		Timeout: time.Second * 90,
	}

	tracingClient = apmhttp.WrapClient(netClient)
)

type TracingClientIface interface {
	HttpGet(ctx context.Context, url string) (*http.Response, error)
	HttpPost(ctx context.Context, url string, bodyType string, body io.Reader) (*http.Response, error)
	HttpHead(ctx context.Context, url string) (*http.Response, error)
	HttpDo(ctx context.Context, r *http.Request) (*http.Response, error)
}

type TracingClient struct {
	HttpClient *http.Client
}

func NewTracingClient(client *http.Client) TracingClientIface {
	return &TracingClient{HttpClient: apmhttp.WrapClient(client)}
}

func (t *TracingClient) HttpGet(ctx context.Context, url string) (*http.Response, error) {
	resp, err := ctxhttp.Get(ctx, t.HttpClient, url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TracingClient) HttpPost(ctx context.Context, url string, bodyType string, body io.Reader) (*http.Response, error) {
	resp, err := ctxhttp.Post(ctx, t.HttpClient, url, bodyType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TracingClient) HttpHead(ctx context.Context, url string) (*http.Response, error) {
	resp, err := ctxhttp.Head(ctx, t.HttpClient, url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TracingClient) HttpDo(ctx context.Context, r *http.Request) (*http.Response, error) {
	resp, err := ctxhttp.Do(ctx, t.HttpClient, r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func HttpGet(ctx context.Context, url string) (*http.Response, error) {
	resp, err := ctxhttp.Get(ctx, tracingClient, url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func HttpPost(ctx context.Context, url string, bodyType string, body io.Reader) (*http.Response, error) {
	resp, err := ctxhttp.Post(ctx, tracingClient, url, bodyType, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func HttpHead(ctx context.Context, url string) (*http.Response, error) {
	resp, err := ctxhttp.Head(ctx, tracingClient, url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
