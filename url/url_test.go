package url

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestUrl(t *testing.T) {
	baseLink := "http://www.baidu.com/news/"
	link := "./a.htm"
	expectedLink := "http://www.baidu.com/news/a.htm"
	realLink, err := ResolveReference(baseLink, link)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedLink, realLink)

	host := "www.baidu.com"
	realHost, err := ExtractHost(realLink)
	assert.Equal(t, nil, err)
	assert.Equal(t, host, realHost)

	domain := "baidu.com"
	realDomain, err := ExtractDomain(host)
	assert.Equal(t, nil, err)
	assert.Equal(t, domain, realDomain)
}
