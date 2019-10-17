package jsondiff

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustReadFile(name string) []byte {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return b
}

func TestDiff(t *testing.T) {
	bp := func(s string) *[]byte {
		return toBytesPtr([]byte(s))
	}
	t.Run("case00", func(t *testing.T) {
		hunks, err := Diff(mustReadFile("./testdata/case00.a.json"), mustReadFile("./testdata/case00.b.json"))
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []Hunk{
			{Path: "$.array[1]", Old: bp(`"b"`), New: nil},
			{Path: "$.boolean", Old: bp(`true`), New: bp(`false`)},
			{Path: "$.number1", Old: bp(`3`), New: bp(`3.0`)},
			{Path: "$.number2", Old: bp(`3.0`), New: bp(`3`)},
			{Path: "$.object.b", Old: bp(`"B"`), New: nil},
			{Path: "$.string", Old: bp(`"testing"`), New: bp(`"TESTING"`)},
		}, hunks)
	})
}
