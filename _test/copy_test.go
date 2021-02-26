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
	t.Run("any type field, embedding field", func(t *testing.T) {
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
			Interface:               "it's a embedding interface",
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
	})

	t.Run("enum copy", func(t *testing.T) {
		type enumA int8
		const (
			enumA1 enumA = iota + 1
			enumA2
			enumA3
		)
		type enumB int32
		const (
			enumB1 enumB = iota + 1
			enumB2
			enumB3
		)

		var src = struct {
			FieldSameEnum      enumA
			FieldDifferentEnum enumA
		}{
			FieldSameEnum:      enumA2,
			FieldDifferentEnum: enumA2,
		}
		var dst struct {
			FieldSameEnum      enumA
			FieldDifferentEnum enumB
		}

		cpr := copy.NewCopier()
		assert.NoError(t, cpr.Copy(&dst, src))

		assert.Equal(t, src.FieldSameEnum, dst.FieldSameEnum)
		assert.Equal(t, int(src.FieldDifferentEnum), int(dst.FieldDifferentEnum))
	})

	t.Run("struct in struct copy", func(t *testing.T) {
		type srcInnerStruct struct {
			Field1 string
			Field2 int
		}
		type dstInnerStruct struct {
			Field1 string
			Field2 int
		}

		var src = struct {
			FieldStruct    srcInnerStruct
			FieldStructPtr *srcInnerStruct
		}{
			FieldStruct: srcInnerStruct{
				Field1: "it's struct",
				Field2: 1,
			},
			FieldStructPtr: &srcInnerStruct{
				Field1: "it's ptr struct",
				Field2: 2,
			},
		}
		var dst struct {
			FieldStruct    dstInnerStruct
			FieldStructPtr *dstInnerStruct
		}

		cpr := copy.NewCopier()
		assert.NoError(t, cpr.Copy(&dst, src))

		assert.Equal(t, src.FieldStruct.Field1, dst.FieldStruct.Field1)
		assert.Equal(t, src.FieldStruct.Field2, dst.FieldStruct.Field2)
		assert.Equal(t, src.FieldStructPtr.Field1, dst.FieldStructPtr.Field1)
		assert.Equal(t, src.FieldStructPtr.Field2, dst.FieldStructPtr.Field2)
	})

	t.Run("struct in slice copy", func(t *testing.T) {
		type srcInnerStruct struct {
			Field1 string
			Field2 int32
		}
		type dstInnerStruct struct {
			Field1 string
			Field2 int8
		}

		var src = struct {
			FieldStructs    []srcInnerStruct
			FieldStructPtrs []*srcInnerStruct
		}{
			FieldStructs: []srcInnerStruct{
				{
					Field1: "it's struct",
					Field2: 1,
				},
			},
			FieldStructPtrs: []*srcInnerStruct{
				{
					Field1: "it's ptr struct",
					Field2: 2,
				},
			},
		}
		var dst struct {
			FieldStructs    []dstInnerStruct
			FieldStructPtrs []*dstInnerStruct
		}

		cpr := copy.NewCopier()
		assert.NoError(t, cpr.Copy(&dst, src))
		if !assert.Len(t, dst.FieldStructs, 1) {
			t.FailNow()
		}
		for i := range src.FieldStructs {
			assert.Equal(t, src.FieldStructs[i].Field1, dst.FieldStructs[i].Field1)
			assert.Equal(t, int(src.FieldStructs[i].Field2), int(dst.FieldStructs[i].Field2))
		}
		if !assert.Len(t, dst.FieldStructPtrs, 1) {
			t.FailNow()
		}
		for i := range src.FieldStructPtrs {
			assert.Equal(t, src.FieldStructPtrs[i].Field1, dst.FieldStructPtrs[i].Field1)
			assert.Equal(t, int(src.FieldStructPtrs[i].Field2), int(dst.FieldStructPtrs[i].Field2))
		}
	})

	t.Run("struct in map copy", func(t *testing.T) {
		type srcInnerStruct struct {
			Field1 string
			Field2 int32
		}
		type dstInnerStruct struct {
			Field1 string
			Field2 int8
		}

		var src = struct {
			FieldStructs    map[int]srcInnerStruct
			FieldStructPtrs map[int]*srcInnerStruct
		}{
			FieldStructs: map[int]srcInnerStruct{
				1: {
					Field1: "it's struct",
					Field2: 1,
				},
			},
			FieldStructPtrs: map[int]*srcInnerStruct{
				1: {
					Field1: "it's ptr struct",
					Field2: 2,
				},
			},
		}
		var dst struct {
			FieldStructs    map[int]dstInnerStruct
			FieldStructPtrs map[int]*dstInnerStruct
		}

		cpr := copy.NewCopier()
		assert.NoError(t, cpr.Copy(&dst, src))
		if !assert.Len(t, dst.FieldStructs, 1) {
			t.FailNow()
		}
		for i := range src.FieldStructs {
			assert.Equal(t, src.FieldStructs[i].Field1, dst.FieldStructs[i].Field1)
			assert.Equal(t, int(src.FieldStructs[i].Field2), int(dst.FieldStructs[i].Field2))
		}
		if !assert.Len(t, dst.FieldStructPtrs, 1) {
			t.FailNow()
		}
		for i := range src.FieldStructPtrs {
			assert.Equal(t, src.FieldStructPtrs[i].Field1, dst.FieldStructPtrs[i].Field1)
			assert.Equal(t, int(src.FieldStructPtrs[i].Field2), int(dst.FieldStructPtrs[i].Field2))
		}
	})

	t.Run("struct in array copy", func(t *testing.T) {
		type srcInnerStruct struct {
			Field1 string
			Field2 int32
		}
		type dstInnerStruct struct {
			Field1 string
			Field2 int8
		}

		var src = struct {
			FieldStructs    map[int]srcInnerStruct
			FieldStructPtrs map[int]*srcInnerStruct
		}{
			FieldStructs: map[int]srcInnerStruct{
				1: {
					Field1: "it's struct",
					Field2: 1,
				},
			},
			FieldStructPtrs: map[int]*srcInnerStruct{
				1: {
					Field1: "it's ptr struct",
					Field2: 2,
				},
			},
		}
		var dst struct {
			FieldStructs    map[int]dstInnerStruct
			FieldStructPtrs map[int]*dstInnerStruct
		}

		cpr := copy.NewCopier()
		assert.NoError(t, cpr.Copy(&dst, src))
		if !assert.Len(t, dst.FieldStructs, 1) {
			t.FailNow()
		}
		for i := range src.FieldStructs {
			assert.Equal(t, src.FieldStructs[i].Field1, dst.FieldStructs[i].Field1)
			assert.Equal(t, int(src.FieldStructs[i].Field2), int(dst.FieldStructs[i].Field2))
		}
		if !assert.Len(t, dst.FieldStructPtrs, 1) {
			t.FailNow()
		}
		for i := range src.FieldStructPtrs {
			assert.Equal(t, src.FieldStructPtrs[i].Field1, dst.FieldStructPtrs[i].Field1)
			assert.Equal(t, int(src.FieldStructPtrs[i].Field2), int(dst.FieldStructPtrs[i].Field2))
		}
	})
}
