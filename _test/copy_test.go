package _test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomtwinkle/go-copy"
)

func TestStringToString(t *testing.T) {
	var src = "sudo i love you"
	var dst = "i love you"
	assert.NoError(t, copy.Copy(&dst, src))
	assert.Equal(t, "sudo i love you", dst)
}

func TestInt64ToInt64(t *testing.T) {
	var src int64 = 64
	var dst int64 = 0
	assert.NoError(t, copy.Copy(&dst, src))
	assert.Equal(t, int64(64), dst)
}

func TestTimeToInt64(t *testing.T) {
	var src = time.Now()
	var dst int64 = 0
	assert.NoError(t, copy.Copy(&dst, src))
	assert.Equal(t, src.Unix(), dst)
}

func TestInt64ToTime(t *testing.T) {
	var src = time.Now().Unix()
	var dst time.Time
	assert.NoError(t, copy.Copy(&dst, src))
	assert.Equal(t, src, dst.Unix())
}

func TestStructToStruct(t *testing.T) {
	now := time.Now()
	type Interface interface{}
	type Struct struct {
		EmbeddingField          string
		DuplicateField          string
		DuplicateFieldDifferent string
	}
	var src = struct {
		Interface
		Struct
		Field1                  int
		Field2                  string
		Field3                  int64
		Field4                  float64
		Field5                  time.Time
		Field6                  interface{}
		DuplicateField          string
		DuplicateFieldDifferent int
	}{
		Struct: Struct{
			EmbeddingField:          "it's a embedding field",
			DuplicateField:          "it's a embedding duplicate field",
			DuplicateFieldDifferent: "it's a embedding duplicate different type field",
		},
		Interface: "it's a embedding interface",
		Field1:                  1,
		Field2:                  "you are a good guy",
		Field3:                  3,
		Field4:                  3.141592654,
		Field5:                  now,
		Field6:                  "it's a interface",
		DuplicateField:          "it's a message to duplicate",
		DuplicateFieldDifferent: 100,
	}
	var dst struct {
		Interface
		Struct
		Field1                  int
		Field2                  string
		Field3                  int64
		Field4                  float64
		Field5                  int64
		Field6                  interface{}
		DuplicateField          string
		DuplicateFieldDifferent int
	}
	cpr := copy.NewCopier()
	assert.NoError(t, cpr.Copy(&dst, src))
	assert.Equal(t, int(1), dst.Field1)
	assert.Equal(t, "you are a good guy", dst.Field2)
	assert.Equal(t, int64(3), dst.Field3)
	assert.Equal(t, float64(3.141592654), dst.Field4)
	assert.Equal(t, now.Unix(), dst.Field5)
	assert.Equal(t, "it's a interface", dst.Field6)
	assert.Equal(t, "it's a message to duplicate", dst.DuplicateField)
	assert.Equal(t, "it's a embedding duplicate field", dst.Struct.DuplicateField)
	assert.Equal(t, 100, dst.DuplicateFieldDifferent)
	assert.Equal(t, "", dst.Struct.DuplicateFieldDifferent) // FIXME
	assert.Equal(t, "it's a embedding field", dst.EmbeddingField)
	assert.Equal(t, "it's a embedding field", dst.Struct.EmbeddingField)
	assert.Equal(t, "it's a embedding interface", dst.Interface)
}
