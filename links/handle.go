package links

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Out struct {
	ok        bytes.Buffer
	errs      bytes.Buffer
	code400   bytes.Buffer
	code401   bytes.Buffer
	code403   bytes.Buffer
	code404   bytes.Buffer
	code500   bytes.Buffer
	codeOther bytes.Buffer
	sync.RWMutex
}

func (r Out) OutPut() {
	_ = os.Mkdir("result", os.ModeDir)
	_ = ioutil.WriteFile("result/ok.md", r.ok.Bytes(), 0644)
	_ = ioutil.WriteFile("result/err.md", r.errs.Bytes(), 0644)

	// 404
	data := bytes.Buffer{}
	s404s := strings.Split(r.code404.String(), "\n")
	for _, v := range s404s {
		data.WriteString(v + "\n")
	}
	_ = ioutil.WriteFile("result/source.md", data.Bytes(), 0644)

	// other
	data.Reset()
	data.WriteString(r.code400.String())
	data.WriteString(r.code401.String())
	data.WriteString(r.code403.String())
	data.WriteString(r.code500.String())
	data.WriteString(r.codeOther.String())
	_ = ioutil.WriteFile("result/fails.md", data.Bytes(), 0644)
}

var (
	out Out
)

// load:en+path
// output:zh+path
func handle(pathname, URL string, wg *sync.WaitGroup, isTLS bool) {
	defer wg.Done()
	//s:=fmt.Sprintf("%v\n",URL)
	//out.Lock()
	//defer out.Unlock()
	//out.code400.WriteString(s)
	//return

	_, err := url.Parse(URL)
	if err != nil {
		fmt.Println(URL)
		return
	}

	client := http.Client{
		Timeout: 15 * time.Second,
	}
	if isTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	req.Close = true
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	out.Lock()
	defer out.Unlock()
	if err != nil {
		s := fmt.Sprintf("- ERR:%v,URL:%v,FILE:%v\n", err, URL, pathname)
		out.errs.WriteString(s)
		return
	}
	defer resp.Body.Close()

	s := fmt.Sprintf("path:%v, [%v](%v)\n", pathname, path.Base(URL), URL)
	//s := fmt.Sprintf("- CODE:%v ,[%v](%v) ,FILE:%v\n", resp.StatusCode, path.Base(URL), URL, pathname)
	if resp.StatusCode < 400 {
		out.ok.WriteString(s)
	} else {
		switch resp.StatusCode {
		case 400:
			out.code400.WriteString(s)
		case 401:
			out.code401.WriteString(s)
		case 403:
			out.code403.WriteString(s)
		case 404:
			out.code404.WriteString(s)
		case 500:
			out.code500.WriteString(s)
		default:
			out.codeOther.WriteString(s)
		}
	}
}
