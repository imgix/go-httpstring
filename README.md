# go-httpstring
parse http header strings efficiently

### cache-control

```golang
cc = NewCacheControl([]byte("max-age=100; must-revalidate"))
if set, d := cc.GetDuration("max-age"); set {
	time.Sleep(d)
}

if cc.IsSet("must-revalidate") {
	...
}
// case-insensitive
if cc.IsSet("Must-Revalidate") {
	...
}

```

