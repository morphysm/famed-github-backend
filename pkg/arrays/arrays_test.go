package arrays

import (
	"reflect"
	"testing"
)

type T any

func TestRemove(t *testing.T) {
	t.Parallel()

	type args struct {
		slice []T
		s     int
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		{
			name: "string",
			args: args{
				slice: []T{"morphysm", 1337, "deadbeef", "66616D6564", "famed"},
				s:     2,
			},
			want: []T{"morphysm", 1337, "66616D6564", "famed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Remove(tt.args.slice, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
