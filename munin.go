package evmn

import (
	"expvar"
	"strings"
	"unicode"
)

func kk(k string) string {
	if k == "" {
		return k
	}
	k = strings.Replace(k, ".", "_", -1)
	if !unicode.IsLetter(rune(k[0])) {
		return "_" + k
	}
	return k
}

func fetch(name string, f func(k, v string)) {
	expvar.Do(func(kv expvar.KeyValue) {
		if kv.Key != name && !strings.HasPrefix(kv.Key, name+":") {
			return
		}
		lst := strings.Split(kv.Key, ":")
		if len(lst) == 1 {
			f(lst[0], kv.Value.String())
			return
		}
		f(lst[1], kv.Value.String())
	})
}
