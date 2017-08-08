package evmn

import (
	"errors"
	"expvar"
	"sort"
	"strings"
)

var (
	// ErrUnknownCmd is given for an unimplemented command
	ErrUnknownCmd = errors.New("Unknown command")

	// ErrUnknownSvc is the response when fetching a nonexistent service
	ErrUnknownSvc = errors.New("Unknown service")
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
		fetch(key, func(k, v string) {
			lines = append(lines, kk(k)+".value "+v)
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

		fetch(key, func(k0, v string) {
			k := kk(k0)
			lines = append(lines,
				k+".label "+k0,
				k+".min 0",
				k+".type DERIVE",
			)
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
