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
	k = strings.Replace(k, " ", "_", -1)
	if !unicode.IsLetter(rune(k[0])) {
		return "_" + k
	}
	return k
}

func fetch(name string, f func(k string, v interface{})) {
	expvar.Do(func(kv expvar.KeyValue) {
		if kv.Key != name && !strings.HasPrefix(kv.Key, name+":") {
			return
		}

		if kv.Key == name {
			f(kv.Key, kv.Value)
			return
		}

		lst := strings.Split(kv.Key, ":")
		f(lst[1], kv.Value)
	})
}
