package export

import (
	"bytes"
	"github.com/Skactor/mitmproxy/logger"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Exporter interface {
	Open(cfg interface{}) error
	WriteBytes(data []byte) error
	WriteInterface(data interface{}) error
	Close() error
}

type OutputRequest struct {
	Url               string            `json:"url"`
	Method            string            `json:"method"`
	StatusCode        int               `json:"status_code"`
	RawRequestHeader  string            `json:"raw_request"`
	RawRequestBody    string            `json:"raw_request_body"`
	RawResponseHeader string            `json:"raw_response"`
	RawResponseBody   string            `json:"raw_response_body"`
	ReqHeaders        map[string]string `json:"request_headers"`
	RespHeaders       map[string]string `json:"response_headers"`
}

func HeadersToMap(headers http.Header) map[string]string {
	result := map[string]string{}
	for name, value := range headers {
		result[name] = strings.Join(value, "")
	}
	return result
}

func OutputRequestFromResponse(req *http.Request, resp *http.Response) (*OutputRequest, error) {
	out := OutputRequest{
		Method:      req.Method,
		Url:         req.URL.String(),
		StatusCode:  resp.StatusCode,
		ReqHeaders:  HeadersToMap(req.Header),
		RespHeaders: HeadersToMap(resp.Header),
	}

	rawRequest, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	splitRequest := bytes.SplitN(rawRequest, []byte{13, 10, 13, 10}, 2)
	out.RawRequestHeader, out.RawRequestBody = string(splitRequest[0]), string(splitRequest[1])

	rawResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}
	splitResponse := bytes.SplitN(rawResponse, []byte{13, 10, 13, 10}, 2)
	out.RawResponseHeader, out.RawResponseBody = string(splitResponse[0]), string(splitResponse[1])
	return &out, nil
}
