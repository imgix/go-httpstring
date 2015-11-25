package httpstring

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCacheControl(t *testing.T) {
	var (
		cc *CacheControl
	)

	//cc = NewCacheControl([]byte("max-age=100, no-cache"))
	cc = NewCacheControl([]byte("max-age=100, must-revalidate"))
	if !cc.IsSet("must-revalidate") {
		t.Errorf("'must-revalidate' not set")

	}
	if d, set := cc.GetDuration("max-age"); !set {
		t.Errorf("'max-age bit not set")

	} else if d != time.Second*100 {
		t.Errorf("'max-age not parsed: %v", d)
	}

	cc = NewCacheControl([]byte("max-age=100, s-maxage=200"))
	fmt.Printf("%v\n", cc)
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
	cc = NewCacheControl([]byte(strings.Join(many, ",")))
	for _, k := range many {
		if !cc.IsSet(k) {
			t.Errorf("%v not set on CacheControl", k)
		}
	}
	cc = NewCacheControl([]byte("public, private, no-cache, no-store, max-age=300"))
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
	cc = NewCacheControl([]byte("public, private, no-cache, no-store, max-age=f00"))
	if _, set := cc.GetDuration("max-age"); set {
		t.Errorf("malformed 'max-age' extracted")

	}
	cc = NewCacheControl([]byte("private=\"foo\""))
	if set, v := cc.GetVal("private"); set {
		if v != "foo" {
			t.Errorf("strings map stored wrong val: %s", v)
		}
	} else {
		t.Errorf("strings map didn't store 'private'")
	}

	cc = NewCacheControl([]byte(""))
	if cc.IsSet("private") {
		t.Errorf("empty cache-control string set 'private' flag")

	}

	cc = NewCacheControl([]byte("public, private, no-cache, no-store, max-age=300"))
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

	cc = NewCacheControl([]byte(""))
	cc.SetDuration("max-age", time.Second*30)
	if dur, ok := cc.GetDuration("max-age"); !ok {
		t.Errorf("fresh object not received max-age")

	} else if dur != 30*time.Second {
		t.Errorf("failed to set duration")
	}
	cc = NewCacheControl([]byte(
		"public, private, no-cache, no-store, max-age=300"),
	)
	cc.ClearVal("public")
	cc.SetDuration("max-age", time.Hour*100)
	if cc.String() != "max-age=360000; private; no-cache; no-store" {
		t.Errorf("to string busted")
	}

}

func ExampleNewCacheControl() {

	cc := NewCacheControl([]byte("max-age=100, must-revalidate"))
	if v, set := cc.GetDuration("max-age"); set {
		time.Sleep(v)

	}

}

func TestFlagNotDuration(t *testing.T) {

	cc := NewCacheControl([]byte("max-stale"))
	if d, ok := cc.GetDuration("max-stale"); ok {
		t.Errorf("non duration max-stale, got duration: %v", d)

	}
	if !cc.IsSet("max-stale") {
		t.Errorf("max-stale not detected")
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
