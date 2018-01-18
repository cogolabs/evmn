package evmn

import (
	"errors"
	"expvar"
	"reflect"
	"sort"
	"strings"
)

var (
	// ErrUnknownCmd is given for an unimplemented command
	ErrUnknownCmd = errors.New("Unknown command")

	// ErrUnknownSvc is the response when fetching a nonexistent service
	ErrUnknownSvc = errors.New("Unknown service")

	// AllowedTypes is the list of types that are allowed to be exported to munin
	AllowedTypes = []string{"*expvar.Int", "*expvar.Map"}
)

func handler(command string) (r string, err error) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return "", ErrUnknownCmd
	}

	switch fields[0] {

	case "list":
		keys := []string{}
		seen := map[string]bool{}
		expvar.Do(func(kv expvar.KeyValue) {
			// Doing nothing if not an allowed type
			allowed := false
			for _, typ := range AllowedTypes {
				if reflect.TypeOf(kv.Value).String() == typ {
					allowed = true
					break
				}
			}
			if !allowed {
				return
			}
			g := strings.Split(kv.Key, ":")[0]
			f := strings.Split(g, ".")[0]
			if seen[f] {
				return
			}
			keys = append(keys, f)
			seen[f] = true
		})
		sort.Strings(keys)
		return strings.Join(keys, " "), nil

	case "nodes":
		return hostname + "\n.", nil

	case "fetch":
		if len(fields) < 2 {
			return "", ErrUnknownSvc
		}
		key := fields[1]
		lines := []string{}
		fetch(key, func(k string, v interface{}) {
			switch t := v.(type) {
			case *expvar.Int:
				lines = append(lines, kk(k)+".value "+t.String())
			case *expvar.Map:
				t.Do(func(kv expvar.KeyValue) {
					lines = append(lines, kk(kv.Key)+".value "+kv.Value.String())
				})
			}
		})
		return strings.Join(lines, "\n") + "\n.", nil

	case "config":
		if len(fields) < 2 {
			return "", ErrUnknownSvc
		}
		key := fields[1]

		lines := []string{
			"graph_title " + key,
			"graph_category expvar",
			"graph_args --base 1000 --units=si",
		}

		fetch(key, func(k0 string, v interface{}) {
			switch t := v.(type) {
			case *expvar.Int:
				k := kk(k0)
				lines = append(lines,
					k+".label "+k0,
					k+".min 0",
					k+".type DERIVE",
				)
			case *expvar.Map:
				t.Do(func(kv expvar.KeyValue) {
					k := kk(kv.Key)
					lines = append(lines,
						k+".label "+kv.Key,
						k+".min 0",
						k+".type DERIVE",
					)
				})
			}
		})
		return strings.Join(lines, "\n") + "\n.", nil

	case "cap":
		return "multigraph", nil

	case "help", "version":
		fallthrough

	default:
		return "", ErrUnknownCmd

	}
}
