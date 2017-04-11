/*
Deprecated: with colm network pulling ragel's golang target,
continuing use and development of this project doesn't make much sense

use https://godoc.org/github.com/pquerna/cachecontrol instead

Package httpstring provides parsing rules to scan http header strings for the purpose
of extracting flags and values.

Current functionality is limited to parsing a Cache-Control header
*/
package httpstring

import (
	"fmt"
	"strings"
	"time"
)

const (
	good_return = 100
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

var durationflags = map[string]cacheflag{
	"max-stale": maxstale,

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
		return
	}
}
func (c *CacheControl) SetDuration(s string, d time.Duration) bool {
	s = strings.ToLower(s)
	if _, b := durationflags[strings.ToLower(s)]; b {
		c.durations[s] = d
		c.setflag(s)
		return true
	}
	return false
}
func (c *CacheControl) SetString(s, v string) bool {
	s = strings.ToLower(s)
	if _, b := cacheflaglookup[s]; b {
		c.setflag(s)
		c.strings[s] = v
		return true
	}
	return false

}

func (c *CacheControl) SetBool(s string) bool {

	s = strings.ToLower(s)
	if _, b := cacheflaglookup[s]; b {
		c.setflag(s)
		return true
	}
	return false
}

func (c *CacheControl) GetDuration(s string) (time.Duration, bool) {
	if c.IsSet(s) {
		d, b := c.durations[s]
		return d, b
	}
	return time.Second, false
}

func (c *CacheControl) GetVal(s string) (bool, string) {
	if c.IsSet(s) {
		s, p := c.strings[s]
		return p, s
	} else {
		return false, ""
	}
}

func (c *CacheControl) ClearVal(s string) bool {

	if f, b := cacheflaglookup[strings.ToLower(s)]; b {
		c.flags = c.flags &^ f
		return true
	}
	return false

}

// IsSet queries CacheControl object for tokens set after parsing
func (c *CacheControl) IsSet(s string) bool {

	if f, ok := cacheflaglookup[strings.ToLower(s)]; ok {
		return c.flags&f == f
	}
	return false
}

func (c *CacheControl) String() string {
	used := make(map[string]bool)
	tokens := make([]string, 0)

	var (
		p bool
		d time.Duration
	)
	if d, p = c.GetDuration("max-age"); p {
		tokens = append(tokens, durationString("max-age", d))
		used["max-age"] = true
	}
	if d, p = c.GetDuration("s-maxage"); p {
		tokens = append(tokens, durationString("s-maxage", d))
		used["s-maxage"] = true
	}

	if d, p = c.GetDuration("min-fresh"); p {
		tokens = append(tokens, durationString("min-fresh", d))
		used["min-fresh"] = true
	}

	for k, _ := range cacheflaglookup {
		if _, p = used[k]; p {
			continue
		}
		if c.IsSet(k) {
			tokens = append(tokens, k)
		}
	}
	return strings.Join(tokens, "; ")

}

func durationString(t string, d time.Duration) string {
	return fmt.Sprintf("%s=%0.f", t, d.Seconds())
}

// NewCacheControl is a struct type used for storing parsed tokens
func NewCacheControl(data []byte) (*CacheControl, error) {
	cc := &CacheControl{
		strings:   make(map[string]string),
		durations: make(map[string]time.Duration),
	}
	if b := run(data, cc); b != good_return {
		return nil, fmt.Errorf("failed to parse http string")
	}
	return cc, nil
}
