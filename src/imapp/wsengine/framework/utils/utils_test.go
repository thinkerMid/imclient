package utils

import (
	"reflect"
	"testing"
)

func TestIntToBigEndianBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{name: "test", args: args{n: 192222391}, want: []byte{11, 117, 20, 183}},
		{name: "test", args: args{n: 183}, want: []byte{183}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToBigEndianBytes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntToBigEndianBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigEndianBytesToInt(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "test", args: args{buf: []byte{11, 117, 20, 183}}, want: 192222391},
		{name: "test", args: args{buf: []byte{183}}, want: 183},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BigEndianBytesToInt(tt.args.buf); got != tt.want {
				t.Errorf("BigEndianBytesToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToLittleEndianBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{name: "test", args: args{n: 192222391}, want: []byte{183, 20, 117, 11}},
		{name: "test", args: args{n: 183}, want: []byte{183}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToLittleEndianBytes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntToLittleEndianBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestIntToBytes(t *testing.T) {
//	type args struct {
//		n int64
//	}
//	tests := []struct {
//		name string
//		args args
//		want []byte
//	}{
//		{name: "test", args: args{n: 192222391}, want: []byte{0, 0, 0, 0, 11, 117, 20, 183}},
//		{name: "test", args: args{n: 183}, want: []byte{0, 0, 0, 0, 0, 0, 0, 183}},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := IntToBytes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("IntToBytes() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestIntToBytes1(t *testing.T) {
//	type args struct {
//		n int
//		b byte
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []byte
//		wantErr bool
//	}{
//		{name: "test", args: args{n: 192222391, b: 4}, want: []byte{11, 117, 20, 183}},
//		{name: "test", args: args{n: 183, b: 1}, want: []byte{183}},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := IntToBytes1(tt.args.n, tt.args.b)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("IntToBytes1() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("IntToBytes1() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestIntToBytes2(t *testing.T) {
//	type args struct {
//		n int64
//	}
//	tests := []struct {
//		name string
//		args args
//		want []byte
//	}{
//		{name: "test", args: args{n: 192222391}, want: []byte{183, 20, 117, 11}},
//		{name: "test", args: args{n: 183}, want: []byte{183}},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := IntToBytes2(tt.args.n); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("IntToBytes2() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestBytesToInt(t *testing.T) {
//	type args struct {
//		b []byte
//	}
//	tests := []struct {
//		name string
//		args args
//		want int
//	}{
//		{name: "test", args: args{b: []byte{11, 117, 20, 183}}, want: 192222391},
//		{name: "test", args: args{b: []byte{183}}, want: 183},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := BytesToInt(tt.args.b); got != tt.want {
//				t.Errorf("BytesToInt() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestRandInt64(t *testing.T) {
	t.Run("TestRandInt64", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			t.Log(RandInt64(1, 3))
		}
	})
}

func BenchmarkGenerateEd25519CredentialGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _ = GenerateEd25519Credential()
	}
}