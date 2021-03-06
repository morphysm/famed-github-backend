package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/phuslu/log"
)

type logMsg struct {
	SendTime string `json:"sendTime"`
	Host     string `json:"host"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   string `json:"status"`
	Error    error  `json:"error"`
	RTT      string `json:"rTT"`
}

type loggingRoundTripper struct {
	rT http.RoundTripper
}

func (lRT loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	sendTime := time.Now()
	res, err := lRT.rT.RoundTrip(req)
	receiveTime := time.Now()

	msg := logMsg{
		SendTime: sendTime.Format(time.RFC3339Nano),
		Host:     req.Host,
		Method:   req.Method,
		RTT:      fmt.Sprintf("%d ms", receiveTime.Sub(sendTime).Milliseconds()),
	}

	if err != nil {
		msg.Error = err
	}

	if res != nil {
		msg.Status = res.Status
	}

	if req.URL != nil {
		msg.Path = req.URL.Path
	}

	log.Info().
		Str("sendTime", msg.SendTime).
		Str("host", msg.Host).
		Str("method", msg.Method).
		Str("path", msg.Path).
		Str("status", msg.Status).
		Err(msg.Error).Msg("Request:")

	return res, err
}

func AddLogging(client *http.Client) *http.Client {
	client.Transport = loggingRoundTripper{client.Transport}
	return client
}
