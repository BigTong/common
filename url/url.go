package url

import (
	"errors"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

const (
	ERR_EMPTY_URL = "empty url"
)

func ResolveReference(base, link string) (string, error) {
	if len(link) == 0 {
		return "", errors.New(ERR_EMPTY_URL)
	}

	iUrl, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	urlString := iUrl.ResolveReference(ref).String()
	val, err := url.QueryUnescape(urlString)
	if err != nil {
		return "", err
	}
	return val, nil
}

func ExtractDomain(host string) (string, error) {
	suffix, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", err
	}
	return suffix, nil
}

func ExtractHost(link string) (string, error) {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return parsedUrl.Host, nil
}
