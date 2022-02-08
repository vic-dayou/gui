package httpsclient

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var client *http.Client

func init() {
	clientKeyPEM := `-----BEGIN PRIVATE KEY-----
Microsoft CSP Name: Microsoft Enhanced Cryptographic Provider v1.0
friendlyName: le-d95503d8-c227-41a9-9261-d96552128e5a
localKeyId: 01000000

MIICXQIBAAKBgQCvJC9MMGRKmxRBI0KMjDtz2KooIc6XOljHPWhTfAamhV3A5v5y
PiZr4haMDpulU08Y0JxsegwDwfbscQrhG7nvilIqIa+HiI1xkfFxjtNUrMN5hpvO
8HUUfwqzb5EdllQcv/C0xxBkeCECIb86JJry7ty4mNBkN2idbGxldMi90QIDAQAB
AoGATvTIIdfbDss06Vyk/smlb8dohmkfQov6Q/AKHUDXmrCbIIDCiuw70/z73y4i
uviAuxYovrqSugryb4tStUMTogmft4methz1/O/083XHwBNKBPnS2fobYDfBxqkX
tH26woCjrEr/O/wngo6iFp7b5yJlyXapN0x+iOF3CShIhAECQQD2gZ6LLYdxSP8i
aRYAPOh10mF5IHt2dl89eOjNiqVGMlkV5aXNT80jAQr/kWGZfIjscb/xkawSKQKs
ovcn99GRAkEAteL02mBrCLfn2idBwXTdil+yeigReAZmRpqQuAfTRZN4RM+5Dw3q
X0IiCkR3oyiwx89n1eGmz1JTZRxoY1AIQQJAWVbQ5xAxLlWOYiJD3wI0Hb+JpCSp
ml18VwMjHJtLGw3US6NXW/m4Fx+hpM5D2STRWyA+uIZbHpnOZlMJ0Gp4gQJBAK38
66JV5y1Q1r2tHc6UHzQ1tMH7wDIjVQSm6FbSTXxZxAt29Rx8gD0dQvi1ZAg0bV7F
fRtwnqPlqZaoJQcTUMECQQD1Dh+Mu3OMb5AHnrtbk9l1qjM3U81QBKdyF0RY+djo
b3cR9I7+hurpqhJmQ7yuvAWe2xWc+YNTQ48FDJTogXlB
-----END PRIVATE KEY-----`
	clientPEM := `-----BEGIN CERTIFICATE-----
localKeyId: 01000000

MIIDrTCCAxagAwIBAgIQWQKhEMePlPB2aPEW+YUIIDANBgkqhkiG9w0BAQUFADAk
MQswCQYDVQQGEwJDTjEVMBMGA1UEChMMQ0ZDQSBURVNUIENBMB4XDTExMDgyNDA3
NDc1MFoXDTEzMDgyNDA3NDc1MFowczELMAkGA1UEBhMCQ04xFTATBgNVBAoTDENG
Q0EgVEVTVCBDQTERMA8GA1UECxMITG9jYWwgUkExFDASBgNVBAsTC0VudGVycHJp
c2VzMSQwIgYDVQQDFBswNDFAWjIwMTEwODI0QFRFU1RAMDAwMDAwMjMwgZ8wDQYJ
KoZIhvcNAQEBBQADgY0AMIGJAoGBAK8kL0wwZEqbFEEjQoyMO3PYqighzpc6WMc9
aFN8BqaFXcDm/nI+JmviFowOm6VTTxjQnGx6DAPB9uxxCuEbue+KUiohr4eIjXGR
8XGO01Ssw3mGm87wdRR/CrNvkR2WVBy/8LTHEGR4IQIhvzokmvLu3LiY0GQ3aJ1s
bGV0yL3RAgMBAAGjggGPMIIBizAfBgNVHSMEGDAWgBRGctwlcp8CTlWDtYD5C9vp
k7P0RTAdBgNVHQ4EFgQUiFLVc/e56LvykZgnvwbiVHMKt0swCwYDVR0PBAQDAgTw
MAwGA1UdEwQFMAMBAQAwOwYDVR0lBDQwMgYIKwYBBQUHAwEGCCsGAQUFBwMCBggr
BgEFBQcDAwYIKwYBBQUHAwQGCCsGAQUFBwMIMIHwBgNVHR8EgegwgeUwT6BNoEuk
STBHMQswCQYDVQQGEwJDTjEVMBMGA1UEChMMQ0ZDQSBURVNUIENBMQwwCgYDVQQL
EwNDUkwxEzARBgNVBAMTCmNybDEyN18xNTcwgZGggY6ggYuGgYhsZGFwOi8vdGVz
dGxkYXAuY2ZjYS5jb20uY246Mzg5L0NOPWNybDEyN18xNTcsT1U9Q1JMLE89Q0ZD
QSBURVNUIENBLEM9Q04/Y2VydGlmaWNhdGVSZXZvY2F0aW9uTGlzdD9iYXNlP29i
amVjdGNsYXNzPWNSTERpc3RyaWJ1dGlvblBvaW50MA0GCSqGSIb3DQEBBQUAA4GB
AFakQbOuB4QHfvewOyDy1mW4iQSRgP2v47QFyExvRk/iOZkUT3tWsYaSLuyRyQV2
eg9lmuMZmB8ITL/0ed7DUsXN7mdoKHmgsBga1Sp8UhR3dznqOSfaAYJqDaIV6gei
TH0Fbj4FTRxcIsf20WzFUN65kkop3hl1ZssxxvA9Asns
-----END CERTIFICATE-----`
	cliCrt, err := tls.X509KeyPair([]byte(clientPEM), []byte(clientKeyPEM))
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

var path = "D:\\Download\\CPCN\\"
var chunk int64 = 1024 * 1000 * 16
var wg sync.WaitGroup

func Download(url string) error {
	urls := strings.Split(url, "/")
	var filename string
	if len(urls) != 0 {
		//https://github.com/yhinan/gui/archive/refs/heads/master.zip
		filename = urls[len(urls)-1]
	} else {
		filename = "temp.log"
	}
	var file *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0666)
	}
	if _, err := os.Stat(path + filename); os.IsNotExist(err) {
		file, err = os.Create(path + filename)
		if err != nil {
			return fmt.Errorf("Create file error: %s ", err.Error())
		}
	} else {
		file, _ = os.Open(path + filename)
	}

	headRequest, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return fmt.Errorf("Head request error %s ", err.Error())
	}

	headResponse, err := client.Do(headRequest)

	size, _ := strconv.ParseInt(headResponse.Header.Get("Content-Length"), 10, 64)

	r := size / chunk
	for i := int64(0); i < r+1; i++ {
		start := i * chunk
		end := (i + 1) * chunk
		if end > size {
			end = size
		}
		request, _ := http.NewRequest("GET", url, nil)
		request.SetBasicAuth("pay", "Pay2020$")
		wg.Add(1)
		go blockSave(file, request, start, end)
	}
	wg.Wait()
	return nil
}

func blockSave(file *os.File, r *http.Request, start, end int64) {
	defer wg.Done()
	response, err := client.Do(r)
	defer response.Body.Close()
	if end != 0 {
		r.Header.Set("Rang", fmt.Sprintf("bytes=%d-%d", start, end))
	}
	if err != nil {
		fmt.Printf("Request Get failed %s ", err.Error())
		return
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusPartialContent {
		fmt.Printf("Response status %d ", response.StatusCode)
		return
	}
	b := make([]byte, end-start+1)
	_, err = response.Body.Read(b)
	if err != nil {
		fmt.Printf("Read bytes from Response.Body error: %s ", err.Error())
		return
	}
	_, err = file.WriteAt(b, int64(start))
	if err != nil {
		fmt.Printf("Save data to file error： %s ", err.Error())
		return
	}
}
