package client

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/BigTong/common/encodings"
)

// all field and method is package private in this file, easy to refactor util.go

var utf8Converter *encodings.Utf8Converter = encodings.NewUtf8Converter()

// support gzip and deflate uncompress
func unCompress(resp *HttpResponse, content []byte) ([]byte, error) {
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Encoding")), "gzip") {
		unCompressedContent, err := unGzipHtml(content)
		if err != nil {
			return content, err
		}
		return unCompressedContent, nil
	}

	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Encoding")), "deflate") {
		unCompressedContent, err := deflate(content)
		if err != nil {
			return content, err
		}
		return unCompressedContent, nil
	}

	return content, nil
}

func unGzipHtml(html []byte) ([]byte, error) {
	var b bytes.Buffer
	b.Write(html)

	ungz, err := gzip.NewReader(&b)
	if err != nil {
		return html, err
	}
	defer ungz.Close()

	content, err := ioutil.ReadAll(ungz)
	if err != nil {
		return html, err
	}

	if len(content) == 0 {
		return html, errors.New("content length is 0")
	}
	return content, nil
}

func deflate(html []byte) ([]byte, error) {
	// for special server data to deflate
	// see http://stackoverflow.com/questions/29513472/golang-compress-flate-module-cant-decompress-valid-deflate-compressed-http-b
	html = bytes.TrimPrefix(html, []byte("\x78\x9c"))
	reader := flate.NewReader(bytes.NewReader(html))
	defer reader.Close()

	enflated, err := ioutil.ReadAll(reader)
	if err != nil {
		return html, err
	}

	if len(enflated) == 0 {
		return html, errors.New("content length is 0")
	}

	return enflated, nil
}
