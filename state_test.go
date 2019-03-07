package jsondiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	t.Run("matchAny", func(t *testing.T) {
		cases := []struct {
			Match bool
			Path  string
			Arg   string
		}{
			{true, "$", "$"},
			{true, "$.hoge", "$"},
			{true, "$.hoge.required", "$..required"},
			{true, "$.hoge.required[0]", "$..required"},
			{false, "$.hoge", "$..required"},
			{false, "$.fuga", "$.hoge"},
			{false, "$.fuga.required", "$.hoge.required"},
		}
		for _, c := range cases {
			var st state
			st.Path = c.Path
			assert.Equal(t, c.Match, st.matchAny(c.Arg), "%q with %q", c.Path, c.Arg)
		}
	})
}
