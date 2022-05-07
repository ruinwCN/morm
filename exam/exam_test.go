package exam

import "testing"

func Test_examModelInsert(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "insert ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			examModelInsert()
		})
	}
}
