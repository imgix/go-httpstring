package httpstring

import (
	"strings"
	"testing"
	"time"
)

func TestCacheControl(t *testing.T) {
	var (
		cc *CacheControl
		e  error
	)

	//cc = NewCacheControl([]byte("max-age=100, no-cache"))
	cc, e = NewCacheControl(
		[]byte("max-age=100, must-revalidate"),
	)
	if e != nil {
		t.Fatal(e)
	}
	if !cc.IsSet("must-revalidate") {
		t.Errorf("'must-revalidate' not set")

	}
	if d, set := cc.GetDuration("max-age"); !set {
		t.Errorf("'max-age bit not set")

	} else if d != time.Second*100 {
		t.Errorf("'max-age not parsed: %v", d)
	}

	cc, e = NewCacheControl([]byte("max-age=100, s-maxage=200"))
	if e != nil {
		t.Fatal(e)
	}
	if d, set := cc.GetDuration("s-maxage"); set {
		if d != time.Second*200 {
			t.Errorf("'s-maxage not parsed: %v", d)
		}

	} else {
		t.Errorf("'s-maxage not found")
	}

	many := []string{
		"public", "private", "no-cache", "no-store",
		"max-stale", "no-transform", "only-if-cached",
		"must-revalidate", "proxy-revalidate",
	}
	cc, e = NewCacheControl([]byte(strings.Join(many, ",")))
	if e != nil {
		t.Fatal(e)
	}

	for _, k := range many {
		if !cc.IsSet(k) {
			t.Errorf("%v not set on CacheControl", k)
		}
	}
	cc, e = NewCacheControl([]byte("public, private, no-cache, no-store, max-age=300"))
	if e != nil {
		t.Fatal(e)
	}
	for _, k := range []string{"public", "private", "no-cache", "no-store"} {
		if !cc.IsSet(k) {
			t.Errorf("%v not set on CacheControl", k)
		}

	}
	if d, set := cc.GetDuration("max-age"); !set {
		t.Errorf("'max-age' bit not set")

	} else if d != time.Second*300 {
		t.Errorf("'max-age' not parsed: %v", d)
	}
	cc, e = NewCacheControl([]byte("public, private, no-cache, no-store, max-age=f00"))
	if e != nil {
		t.Fatal(e)
	}
	if _, set := cc.GetDuration("max-age"); set {
		t.Errorf("malformed 'max-age' extracted")

	}
	cc, e = NewCacheControl([]byte("private=\"foo\""))
	if e != nil {
		t.Fatal(e)
	}
	if set, v := cc.GetVal("private"); set {
		if v != "foo" {
			t.Errorf("strings map stored wrong val: %s", v)
		}
	} else {
		t.Errorf("strings map didn't store 'private'")
	}

	cc, e = NewCacheControl([]byte("lprivate=\"foo\""))

	if e != nil {
		t.Fatal(e)
	}
	if set, v := cc.GetVal("lprivate"); set {
		if v != "foo" {
			t.Errorf("strings map stored wrong val: %s", v)
		}
	} else {
		t.Errorf("strings map didn't store 'private'")
	}

	cc, e = NewCacheControl([]byte(""))
	if e != nil {
		t.Fatal(e)
	}
	if cc.IsSet("private") {
		t.Errorf("empty cache-control string set 'private' flag")

	}

	cc, e = NewCacheControl([]byte("public, private, no-cache, no-store, max-age=300"))
	if e != nil {
		t.Fatal(e)
	}
	cc.ClearVal("public")
	if cc.IsSet("public") {
		t.Errorf("failed to clear 'public' flag'")
	}

	cc.SetDuration("max-age", 30*time.Second)
	if dur, ok := cc.GetDuration("max-age"); !ok {
		t.Errorf("lost max-age setting in SetDuration")

	} else if dur != 30*time.Second {
		t.Errorf("failed to set new duration")
	}

	cc.ClearVal("max-age")
	if _, ok := cc.GetDuration("max-age"); ok {
		t.Errorf("failed to clear 'max-age' flag'")
	}

	cc, e = NewCacheControl([]byte(""))
	if e != nil {
		t.Fatal(e)
	}
	cc.SetDuration("max-age", time.Second*30)
	if dur, ok := cc.GetDuration("max-age"); !ok {
		t.Errorf("fresh object not received max-age")

	} else if dur != 30*time.Second {
		t.Errorf("failed to set duration")
	}
	cc, e = NewCacheControl([]byte(
		"public, private, no-cache, no-store, max-age=300"),
	)
	if e != nil {
		t.Fatal(e)
	}

	cc.ClearVal("public")
	cc.SetDuration("max-age", time.Hour*100)
	s := cc.String()
	// if !strings.Contains(s, "private") || !strings.Contains(s, "no-cache") || !strings.Contains(s, "no-nostore") || !strings.Contains(s, "max-age=360000") {
	if !strings.Contains(s, "private") || !strings.Contains(s, "no-cache") || !strings.Contains(s, "no-store") || !strings.Contains(s, "max-age=360000") {

		t.Errorf("to string busted")
		t.Errorf("'%s' != '%s'", cc.String(),
			"max-age=360000; private; no-cache; no-store")
	}

}

func ExampleNewCacheControl() {

	if cc, e := NewCacheControl([]byte("max-age=100, must-revalidate")); e != nil {
		panic(e)
	} else {
		if v, set := cc.GetDuration("max-age"); set {
			time.Sleep(v)

		}
	}

}

func TestDegenerativeCase(t *testing.T) {
	_, e := NewCacheControl([]byte(
		"max-age=100, no-cache, no-store"),
	)
	if e != nil {
		t.Fatal(e)
	}

}

func TestFlagNotDuration(t *testing.T) {

	if cc, e := NewCacheControl([]byte("max-stale")); e != nil {
	} else {
		if d, ok := cc.GetDuration("max-stale"); ok {
			t.Errorf("non duration max-stale, got duration: %v", d)

		}
		if !cc.IsSet("max-stale") {
			t.Errorf("max-stale not detected")
		}
	}
	/*

		cc = NewCacheControl([]byte("max-stale=10"))
		if d, ok := cc.GetDuration("max-stale"); !ok {
			t.Errorf("duration max-stale not detected")
		} else if d != time.Second*10 {
			t.Errorf("duration max-stale incorrect: %v", d)

		}
	*/
}
