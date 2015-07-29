/*

Package httpstring provides parsing rules to scan http header strings for the purpose
of extracting flags and values.

Current functionality is limited to parsing a Cache-Control header
*/
package httpstring

import (
	"strings"
	"time"
)

type cacheflag uint

const (
	private cacheflag = 1 << iota
	nocache
	nostore
	public
	notransform
	onlyifcached
	mustrevalidate
	proxyrevalidate
	maxstale

	maxage
	smaxage
	minfresh
)

var cacheflaglookup = map[string]cacheflag{
	"private":          private,
	"no-cache":         nocache,
	"no-store":         nostore,
	"public":           public,
	"no-transform":     notransform,
	"only-if-cached":   onlyifcached,
	"must-revalidate":  mustrevalidate,
	"proxy-revalidate": proxyrevalidate,
	"max-stale":        maxstale,

	"max-age":   maxage,
	"s-maxage":  smaxage,
	"min-fresh": minfresh,
}

type CacheControl struct {
	flags     cacheflag
	strings   map[string]string
	durations map[string]time.Duration
}

func (c *CacheControl) setflag(s string) {
	if f, ok := cacheflaglookup[strings.ToLower(s)]; ok {
		c.flags = c.flags | f
	}

}

func (c *CacheControl) GetDuration(s string) (time.Duration, bool) {
	t, ok := c.durations[s]
	return t, ok
}

func (c *CacheControl) GetVal(s string) (bool, string) {
	s, p := c.strings[s]
	return p, s
}

// IsSet queries CacheControl object for tokens set after parsing
func (c *CacheControl) IsSet(s string) bool {
	if f, ok := cacheflaglookup[strings.ToLower(s)]; ok {
		return c.flags&f == f
	}
	return false
}
