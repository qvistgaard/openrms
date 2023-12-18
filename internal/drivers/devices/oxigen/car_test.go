package oxigen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Conversion(t *testing.T) {
	toByte := percentageToByte(100)
	u := toByte >> 1

	fmt.Printf("%X\n", toByte)
	fmt.Printf("%X\n", u)
	fmt.Printf("%d\n", u)
}

func Test_percentageToByte(t *testing.T) {
	type args struct {
		percent uint8
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "100 percent",
			want: 255,
			args: args{
				percent: 100,
			},
		},
		{
			name: "0 percent",
			want: 0,
			args: args{
				percent: 0,
			},
		},
		{
			name: "50 percent",
			want: 127,
			args: args{
				percent: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, percentageToByte(tt.args.percent), "percentageToByte(%v)", tt.args.percent)
		})
	}
}
