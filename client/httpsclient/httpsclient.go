package httpsclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var client *http.Client

func init() {
	cliCrt, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Fatalln(err)
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{
		Certificates:       []tls.Certificate{cliCrt},
		InsecureSkipVerify: true,
	}}

	client = &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}
}

func Get(url string) (*http.Response, error) {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth("pay", "Pay2020$")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func DownloadFile(file, url string) error {
	response, err := Get(url)
	if err != nil {
		return fmt.Errorf("Request Get failed %s ", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status %d ", response.StatusCode)
	}
	out, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Create file error %s ", err.Error())
	}

	defer out.Close()
	defer response.Body.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("Save data to file error： %s ", err.Error())
	}
	return nil
}
