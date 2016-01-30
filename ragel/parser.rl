package httpstring

import (
	"strings"
	"time"
)


%% machine cache_control_parser;
%% write data;

// NewCacheControl is a struct type used for storing parsed tokens
func NewCacheControl(data []byte) *CacheControl {
	
	cs, p, pe := 0, 0, len(data)
	eof := pe

	str_beg := 0
	str_end := 0
	num, neg := 0, 1
	val := false


	cc := &CacheControl{
		strings: make(map[string]string),
		durations: make(map[string]time.Duration),
	}

	setbool := func(s string) {cc.setflag(s) }

	setstrings := func(k string) { 
		setbool(k)
		cc.strings[strings.ToLower(k)] = string(data[str_beg:str_end])
	}
	setduration := func(k string) { 
		setbool(k)
		v := time.Duration(int64(neg * num))
		cc.durations[k] = time.Second * v
		num = 0
		neg = 1
		val = false
	}

%%{
	action private          { setstrings("private") }
	action no_cache         { setstrings("no-cache") }

	action no_store         { setbool("no-store"); }
	action public           { setbool("public"); }

	action no_transform     { setbool("no-transform") }
	action only_if_cached   { setbool("only-if-cached") }
	action must_revalidate  { setbool("must-revalidate") }
	action proxy_revalidate { setbool("proxy-revalidate") }

	action max_age          { setduration("max-age") }
	action s_maxage         { setduration("s-maxage") }
	action min_fresh        { setduration("min-fresh") }


	action max_stale {
		setbool("max-stale");
		if val {
			setduration("max-stale")
		}
	}

	action str_init { str_beg = p }
	action str_end  { str_end = p }
	action sec_neg  { neg = -1 }
	action sec_incr { num = num * 10 + (int(fc)-'0') }

	tchar = [!#-'*-+\-.0-:@A-Z^-`a-z|~];
	value_token = tchar+ >str_init %str_end ;
	value_string = '"' ( [^"]* >str_init %str_end ) '"' ;
	value = '=' ( value_token | value_string ) ;
	integer = '0' | [1-9] digit* ;
	seconds = '=' ( ( "-" @sec_neg )? integer >{ val = true; } @sec_incr ) ;
	extension = ( tchar+ ) ( value )? ;

	def = "public"i %public
		| "private"i value? %private
		| "no-cache"i value? %no_cache
		| "no-store"i %no_store
		| "max-age"i seconds %max_age
		| "s-maxage"i seconds %s_maxage
		| "max-stale"i ( seconds )? %max_stale
		| "min-fresh"i seconds %min_fresh
		| "no-transform"i %no_transform
		| "only-if-cached"i %only_if_cached
		| "must-revalidate"i %must_revalidate
		| "proxy-revalidate"i %proxy_revalidate
		;

	ext = extension - def ;
	dir = ( def | ext );
	main := space* dir ( space* [,;] space* dir )* ;
	write init;
	write exec;
}%%

	return cc

}
