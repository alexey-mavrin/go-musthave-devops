package common

import "testing"

func TestFirstSet(t *testing.T) {
	// fistStr := "first"
	secondStr := "second"
	thirdStr := "third"

	type args struct {
		s []*string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "three elements, first is nil",
			args: args{
				s: []*string{
					nil,
					&secondStr,
					&thirdStr,
				},
			},
			want: &secondStr,
		},
		{
			name: "all elemets are nil",
			args: args{
				s: []*string{
					nil,
					nil,
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstSet(tt.args.s...); got != tt.want {
				t.Errorf("FirstSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
