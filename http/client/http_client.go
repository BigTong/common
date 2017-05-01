package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BigTong/common/dns"
	"github.com/BigTong/common/http/util"
)

const (
	HTTP_GET  = "GET"
	HTTP_POST = "POST"
	HTTP_PUT  = "PUT"
)

type HttpClient struct {
	// request timeout threshold
	timeout int

	// local ip. If this field is setted, dial remote address with this ip.
	ip string

	// support use proxy
	proxy string

	// use cookieJar
	needCookie bool

	// re-use tcp connection
	keepAlive bool

	// dns client interface, resolve host and ip address
	dnsClient dns.DnsClientInterface

	// http.Client
	httpClient *http.Client
}

func (h *HttpClient) TimeOut() int {
	return h.timeout
}

func (h *HttpClient) Ip() string {
	return h.ip
}

func (h *HttpClient) Proxy() string {
	return h.proxy
}

func (h *HttpClient) NeedCookie() bool {
	return h.needCookie
}

func (h *HttpClient) KeepAlive() bool {
	return h.keepAlive
}

func (h *HttpClient) HttpClient() *http.Client {
	return h.httpClient
}

// simple get method, use this method without
// specific option, use more easily
func (h *HttpClient) SimpleGet(url string,
	header map[string]string) (*HttpResponse, error) {
	return h.Get(url, header, nil, nil, nil)
}

// GET method, only use url encoded params
func (h *HttpClient) Get(url string, headers map[string]string,
	params map[string]string, cookies []*http.Cookie,
	option *FetchOption) (*HttpResponse, error) {
	var httpResponse *HttpResponse = NewHttpResponse()

	request, err := http.NewRequest(HTTP_GET, url, nil)

	if err != nil {
		httpResponse.Status = NewRequestErr
		return httpResponse, errors.New(
			fmt.Sprintf("new request get error: %s", err.Error()))
	}

	request.URL.RawQuery += util.UrlencodedFrom(params)

	return h.do(request, httpResponse, headers, cookies, option)
}

// simple post method, use url encoded type as params
func (h *HttpClient) UrlEncodedPost(url string,
	header map[string]string,
	params map[string]string) (*HttpResponse, error) {
	return h.setHeaderContentTypeAndDo(url, HTTP_POST,
		header, strings.NewReader(util.UrlencodedFrom(params)),
		FORM_URLENCODED, nil, nil)
}

// simple post method, use json type as params
func (h *HttpClient) JsonPost(url string, header map[string]string,
	body []byte) (*HttpResponse, error) {
	return h.setHeaderContentTypeAndDo(url, HTTP_POST,
		header, bytes.NewReader(body), JSON, nil, nil)
}

func (h *HttpClient) JsonPut(url string, header map[string]string,
	body []byte) (*HttpResponse, error) {
	return h.setHeaderContentTypeAndDo(url, HTTP_PUT, header,
		bytes.NewReader(body), JSON, nil, nil)
}

// accept body as post params, auto add Content-Type in header
// according to FormContentType
func (h *HttpClient) setHeaderContentTypeAndDo(url, method string,
	headers map[string]string, body io.Reader,
	contentType FormContentType, cookies []*http.Cookie,
	option *FetchOption) (*HttpResponse, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers["Content-Type"]; !ok {
		if len(contentType) == 0 {
			contentType = FORM_URLENCODED
		}
		headers["Content-Type"] = string(contentType)
	}

	var httpResponse *HttpResponse = NewHttpResponse()

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		httpResponse.Status = NewRequestErr
		return httpResponse, errors.New(
			fmt.Sprintf("new request get error: %s", err.Error()))
	}

	return h.do(request, httpResponse, headers, cookies, option)
}

func (h *HttpClient) do(request *http.Request,
	httpResponse *HttpResponse, headers map[string]string,
	cookies []*http.Cookie, option *FetchOption) (*HttpResponse, error) {

	if option == nil {
		option = GetDefaultFetchOption()
	}

	setHeader(request, headers, option)
	setCookies(request, cookies)

	var (
		response *http.Response
		err      error
	)

	if option.GetFollowRedirect() {
		response, err = h.httpClient.Do(request)
	} else {
		response, err = h.httpClient.Transport.RoundTrip(request)
	}
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		httpResponse.Status = DoRequestErr
		return httpResponse, errors.New(
			fmt.Sprintf("do request get error: %s", err.Error()))
	}

	httpResponse.Status = response.StatusCode
	httpResponse.Cookies = response.Cookies()
	httpResponse.Header = response.Header
	if httpResponse.Status != StatusOK {
		return httpResponse, errors.New(
			fmt.Sprintf("http status is not StatusOK: %d",
				httpResponse.Status))
	}

	htmlByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		httpResponse.Status = ReadTimeOut
		return httpResponse, errors.New(
			fmt.Sprintf("read body get error: %s, len %d",
				err.Error(), len(htmlByte)))
	}

	unCompressedContent, err := unCompress(httpResponse, htmlByte)
	if err != nil {
		httpResponse.Status = UncompressErr
		return httpResponse, errors.New(
			fmt.Sprintf("unCompress get error: %s",
				err.Error()))
	}
	if len(unCompressedContent) == 0 {
		httpResponse.Status = UpexpectedErr
	}

	httpResponse.Body = unCompressedContent
	if option.GetTransformToUtf8() {
		httpResponse.Body = utf8Converter.ToUTF8(unCompressedContent)
	} else {
		httpResponse.Body = unCompressedContent
	}

	return httpResponse, nil
}

func setCookies(request *http.Request, cookies []*http.Cookie) {
	if cookies == nil {
		return
	}
	for _, ck := range cookies {
		request.AddCookie(ck)
	}
}

func setHeader(req *http.Request,
	headers map[string]string, option *FetchOption) {
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	req.Header.Set("User-Agent", GetHeaderUA(option.GetUseWapUserAgent()))

	// use Header.Set to override default header
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
