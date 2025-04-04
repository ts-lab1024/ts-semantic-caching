package main

// This file lifted wholesale from mountainflux by Mark Rushakoff.

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	httpClientName        = "tsbs_load_influx"
	headerContentEncoding = "Content-Encoding"
	headerGzip            = "gzip"
)

var (
	errBackoff          = fmt.Errorf("backpressure is needed")
	backoffMagicWords0  = []byte("engine: cache maximum memory size exceeded")
	backoffMagicWords1  = []byte("write failed: hinted handoff queue not empty")
	backoffMagicWords2a = []byte("write failed: read message type: read tcp")
	backoffMagicWords2b = []byte("i/o timeout")
	backoffMagicWords3  = []byte("write failed: engine: cache-max-memory-size exceeded")
	backoffMagicWords4  = []byte("timeout")
	backoffMagicWords5  = []byte("write failed: can not exceed max connections of 500")
)

type HTTPWriterConfig struct {
	Host      string
	Database  string
	DebugInfo string
}

type HTTPWriter struct {
	client fasthttp.Client

	c   HTTPWriterConfig
	url []byte
}

func NewHTTPWriter(c HTTPWriterConfig, consistency string) *HTTPWriter {
	return &HTTPWriter{
		client: fasthttp.Client{
			Name: httpClientName,
		},

		c:   c,
		url: []byte(c.Host + "/write?consistency=" + consistency + "&db=" + url.QueryEscape(c.Database)),
	}
}

var (
	methodPost = []byte("POST")
	textPlain  = []byte("text/plain")
)

func (w *HTTPWriter) initializeReq(req *fasthttp.Request, body []byte, isGzip bool) {
	req.Header.SetContentTypeBytes(textPlain)
	req.Header.SetMethodBytes(methodPost)
	req.Header.SetRequestURIBytes(w.url)
	if isGzip {
		req.Header.Add(headerContentEncoding, headerGzip)
	}
	req.SetBody(body)
}

func (w *HTTPWriter) executeReq(req *fasthttp.Request, resp *fasthttp.Response) (int64, error) {
	start := time.Now()
	err := w.client.Do(req, resp)
	lat := time.Since(start).Nanoseconds()
	if err == nil {
		sc := resp.StatusCode()
		if sc == 500 && backpressurePred(resp.Body()) {
			err = errBackoff
		} else if sc != fasthttp.StatusNoContent {
			err = fmt.Errorf("[DebugInfo: %s] Invalid write response (status %d): %s", w.c.DebugInfo, sc, resp.Body())
		}
	}
	return lat, err
}

func (w *HTTPWriter) WriteLineProtocol(body []byte, isGzip bool) (int64, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	w.initializeReq(req, body, isGzip)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	return w.executeReq(req, resp)
}

func backpressurePred(body []byte) bool {
	if bytes.Contains(body, backoffMagicWords0) {
		return true
	} else if bytes.Contains(body, backoffMagicWords1) {
		return true
	} else if bytes.Contains(body, backoffMagicWords2a) && bytes.Contains(body, backoffMagicWords2b) {
		return true
	} else if bytes.Contains(body, backoffMagicWords3) {
		return true
	} else if bytes.Contains(body, backoffMagicWords4) {
		return true
	} else if bytes.Contains(body, backoffMagicWords5) {
		return true
	} else {
		return false
	}
}
