package tutorials

import (
	"crypto/tls"
	"net/http"
	"time"
)

func send(sign, msg, url string) ([]byte, error) {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Timeout: 60 * time.Second, Transport: t}
	client.Post(url, "")

	return nil, nil

}
