package buffer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBufferWrites(t *testing.T) {
	buf := NewPool().Get()

	tests := []struct {
		desc string
		f    func()
		want string
	}{
		{"AppendByte", func() { buf.AppendByte('v') }, "v"},
		{"AppendString", func() { buf.AppendString("foo") }, "foo"},
		{"AppendIntPositive", func() { buf.AppendInt(42) }, "42"},
		{"AppendIntNegative", func() { buf.AppendInt(-42) }, "-42"},
		{"AppendUint", func() { buf.AppendUint(42) }, "42"},
		{"AppendBool", func() { buf.AppendBool(true) }, "true"},
		{"AppendFloat64", func() { buf.AppendFloat(3.14, 64) }, "3.14"},
		// Intenationally introduce some floating-point error.
		{"AppendFloat32", func() { buf.AppendFloat(float64(float32(3.14)), 32) }, "3.14"},
		{"AppendWrite", func() { buf.Write([]byte("foo")) }, "foo"},
		{"AppendTime", func() { buf.AppendTime(time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC), time.RFC3339) }, "2000-01-02T03:04:05Z"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			buf.Reset()
			tt.f()
			assert.Equal(t, tt.want, buf.String(), "Unexpected buffer.String().")
			assert.Equal(t, tt.want, string(buf.Bytes()), "Unexpected string(buffer.Bytes()).")
			assert.Equal(t, len(tt.want), buf.Len(), "Unexpected buffer length.")
			// We're not writing more than a kibibyte in tests.
			assert.Equal(t, _size, buf.Cap(), "Expected buffer capacity to remain constant.")
		})
	}
}
