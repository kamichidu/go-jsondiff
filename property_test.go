package jsondiff

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePathTester(t *testing.T) {
	t.Run("", func(t *testing.T) {
		v := createPathTester([]string{"$"})
		assert.Equal(t, struct{}{}, v)
	})
	t.Run("", func(t *testing.T) {
		v := createPathTester([]string{"$", "hoge", "fuga"})
		assert.Equal(t, map[string]interface{}{
			"hoge": map[string]interface{}{
				"fuga": struct{}{},
			},
		}, v)
	})
	t.Run("", func(t *testing.T) {
		v := createPathTester([]string{"$", "hoge", "fuga[1]"})
		assert.Equal(t, map[string]interface{}{
			"hoge": map[string]interface{}{
				"fuga": []interface{}{
					nil,
					struct{}{},
				},
			},
		}, v)
	})
}
