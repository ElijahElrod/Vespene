package signal

import (
	"reflect"
	"testing"
)

func TestNewDonchian(t *testing.T) {

	type args struct{}

	tests := []struct {
		name    string
		args    args
		want    *Donchian
		wantErr bool
	}{
		{
			name: "Test New Donchian",
			args: args{},
			want: &Donchian{upper: 0,
				lower:      0,
				mid:        0,
				lowPeriod:  40,
				highPeriod: 50,
				channel:    DonchianChannel,
				lows:       make([]float64, 0, 40),
				highs:      make([]float64, 0, 50)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDonchian(50, 40)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDonchian() got = %v, want %v", got, tt.want)
			}
		})
	}
}
