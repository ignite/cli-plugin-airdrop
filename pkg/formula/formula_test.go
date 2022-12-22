package formula

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestValue_Calculate(t *testing.T) {
	type args struct {
		amount math.Int
		staked math.Int
	}
	tests := []struct {
		name  string
		value Value
		args  args
		want  math.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.Calculate(tt.args.amount, tt.args.staked)
			require.Equal(t, tt.want, got)
		})
	}
}
