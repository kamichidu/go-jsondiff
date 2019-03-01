package jsondiff

import (
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

type Option func(*state)

func WithLogger(v *log.Logger) Option {
	return func(st *state) {
		st.Logger = v
	}
}

func WithIgnorePath(v string) Option {
	if _, err := jsonpath.New(v); err != nil {
		panic(err)
	}
	return func(st *state) {
		st.IgnorePaths = append(st.IgnorePaths, v)
	}
}

func WithSetPath(v string) Option {
	if _, err := jsonpath.New(v); err != nil {
		panic(err)
	}
	return func(st *state) {
		st.SetPaths = append(st.SetPaths, v)
	}
}

type state struct {
	// NOTE: format only similar to $.property or $.array[0]
	Path string

	SetPaths []string

	IgnorePaths []string

	Logger *log.Logger
}

func (st state) PushState(suffix string) state {
	st.Path = st.Path + suffix
	return st
}

func (st state) matchAny(paths ...string) bool {
	tester := createPathTester(strings.Split(st.Path, "."))
	for _, p := range paths {
		v, err := jsonpath.Get(p, tester)
		// err means invalid jsonpath or "unknown key xxx"
		// if invalid jsonpath error rejects by WithIgnorePath
		// then this err means only "unkwnon key xxx"
		if err != nil {
			continue
		}
		switch val := v.(type) {
		case []interface{}:
			if len(val) > 0 {
				return true
			}
		default:
			if val != nil {
				return true
			}
		}
	}
	return false
}

func (st state) IsIgnored() bool {
	return st.matchAny(st.IgnorePaths...)
}

func (st state) IsSet() bool {
	return st.matchAny(st.SetPaths...)
}
