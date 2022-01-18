package httpclient

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type NameValuePair struct {
	Key, Value string
}

var client *http.Client

const xWWWFormUrlEncoded = "application/x-www-form-urlencoded"

func init() {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client = &http.Client{Timeout: 60 * time.Second, Transport: t}
}

func Post(params []NameValuePair, url string) ([]byte, error) {
	response, err := client.Post(url, xWWWFormUrlEncoded, strings.NewReader(buildParameter(params)))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {

	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(body)
	}
	return body, nil
}

func Get() {

}

func buildParameter(params []NameValuePair) string {
	var sb strings.Builder
	for i, v := range params {
		sb.WriteString(fmt.Sprintf("%s=%s", v.Key, v.Value))
		if i < len(params)-1 {
			sb.WriteString("&")
		}
	}
	return sb.String()
}
