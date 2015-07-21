# go-httpstring
parse http header strings efficiently

## cache-control

### durations
- max-age
- s-maxage
- min-fresh
- max-stale (duration optional)

parse cache-control, get max-age
```golang

cc = NewCacheControl([]byte("max-age=100; must-revalidate"))
if cc.IsSet("max-age") {
	_, d := cc.GetDuration("max-age")
	time.Sleep(d)
}
```


### flags
- no-store
- no-transform
- must-revalidate
- proxy-revalidate
- private (optional string value)
- no-cache (optional string value)

check for flag
```golang
cc = NewCacheControl([]byte("max-age=100; must-revalidate"))
if cc.IsSet("must-revalidate") {

	...
}

// keys are case-insensitive
if cc.IsSet("Must-Revalidate") {
	...
}
```

get a string value
```golang
cc = NewCacheControl([]byte("private=\"foo\";))
// returns whether flag is set
if cc.IsSet("private") {

}

// check flag set, retreive value one shot
if set, v := cc.GetVal("private"); set {
	v == "foo"
}
```

