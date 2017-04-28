package urlutil

import (
	"golang.org/x/net/publicsuffix"
	"net/url"
	// "github.com/BigTong/common/log"
)

func ResolveReference(base, link string) string {
	if len(link) == 0 {
		return ""
	}

	iUrl, err := url.Parse(base)
	if err != nil {
		return link
	}

	ref, err := url.Parse(link)
	if err != nil {
		// log.Error("unescape url get err:%s", link)
		return link
	}

	urlString := iUrl.ResolveReference(ref).String()
	val, err := url.QueryUnescape(urlString)
	if err != nil {
		// log.Error("unescape url get err:%s", link)
		return link
	}
	return val
}

func ExtractDomain(host string) string {
	suffix, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		// log.Warn("parse domain get error:%s", err.Error())
		return host
	}
	return suffix
}

func ExtractHost(link string) string {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		// log.Warn("not parsed url successs:%s", link)
		return link
	}
	return parsedUrl.Host
}
