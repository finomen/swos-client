package swos_client

import (
	"net"
	"reflect"
	"testing"
)

func Test_structToMikrotik(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "linkChange",
			args: args{
				value: linkChange{
					Nm:   []string{"506f727431", "506f727432", "506f727433", "506f727434", "506f727435", "534650"},
					En:   0x3f,
					An:   0x3f,
					Spdc: []int{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					Dpxc: 0x3f,
					Fct:  0x3f,
					Prio: []int{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					Poe:  []int{0x01, 0x01, 0x01, 0x02, 0x03, 0x00},
				},
			},
			want: "{nm:['506f727431','506f727432','506f727433','506f727434','506f727435','534650'],en:0x3f,an:0x3f,spdc:[0x00,0x00,0x00,0x00,0x00,0x00],dpxc:0x3f,fct:0x3f,prio:[0x00,0x00,0x00,0x00,0x00,0x00],poe:[0x01,0x01,0x01,0x02,0x03,0x00]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := anyToMikrotik(tt.args.value); got != tt.want {
				t.Errorf("structToMikrotik() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayToBitMask(t *testing.T) {
	type args struct {
		in []bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "All ones",
			args: args{
				in: []bool{true, true, true, true, true, true},
			},
			want: 0x3f,
		},
		{
			name: "All zeroes",
			args: args{
				in: []bool{false, false, false, false, false, false},
			},
			want: 0x0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arrayToBitMask(tt.args.in); got != tt.want {
				t.Errorf("arrayToBitMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_macFromMikrotik(t *testing.T) {
	type args struct {
		in string
	}

	addr, _ := net.ParseMAC("f4:1e:57:5c:86:7a")
	tests := []struct {
		name    string
		args    args
		want    net.HardwareAddr
		wantErr bool
	}{
		{
			name: "Valid MAC",
			args: args{
				in: "f41e575c867a",
			},
			want:    addr,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := macFromMikrotik(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("macFromMikrotik() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("macFromMikrotik() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ipFromMikrotik(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		{
			name: "Valid IP",
			args: args{
				in: "0x0158a8c0",
			},
			want:    net.IP{192, 168, 88, 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipFromMikrotik(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ipFromMikrotik() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ipFromMikrotik() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ipToMikrotik(t *testing.T) {
	type args struct {
		in net.IP
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Valid IP",
			args: args{
				in: net.IP{192, 168, 88, 1},
			},
			want: 0x0158a8c0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipToMikrotik(tt.args.in); got != tt.want {
				t.Errorf("ipToMikrotik() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringToMikrotik(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid string",
			args: args{
				in: "IoT Switch3",
			},
			want: "496f542053776974636833",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToMikrotik(tt.args.in); got != tt.want {
				t.Errorf("stringToMikrotik() = %v, want %v", got, tt.want)
			}
		})
	}
}
