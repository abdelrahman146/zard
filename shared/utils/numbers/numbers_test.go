package numbers

import "testing"

func TestStruct_GenerateRandomDigits(t *testing.T) {
	type args struct {
		digits int
	}
	tests := []struct {
		name     string
		args     args
		validate func(got int) bool
		want     string
	}{
		{
			name: "Test GenerateRandomDigits with 4 digits",
			args: args{
				digits: 4,
			},
			validate: func(got int) bool {
				return got >= 1000 && got <= 9999
			},
			want: "val >= 1000 && val <= 9999",
		},
		{
			name: "Test GenerateRandomDigits with 6 digits",
			args: args{
				digits: 6,
			},
			validate: func(got int) bool {
				return got >= 100000 && got <= 999999
			},
			want: "val >= 100000 && val <= 999999",
		},
		{
			name: "Test GenerateRandomDigits with 0 digits",
			args: args{
				digits: 0,
			},
			validate: func(got int) bool {
				return got == 0
			},
			want: "val == 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Struct{}
			if got, _ := n.GenerateRandomDigits(tt.args.digits); !tt.validate(got) {
				t.Errorf("Struct.GenerateRandomDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct_GenerateRandomInt(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name     string
		args     args
		validate func(got int) bool
		want     string
	}{
		{
			name: "Test Generate Random Int between 1 and 10",
			args: args{
				min: 1,
				max: 10,
			},
			validate: func(got int) bool {
				return got >= 1 && got <= 10
			},
			want: "val >= 1 && val <= 10",
		},
		{
			name: "Test Generate Random Int between 10 and 20",
			args: args{
				min: 10,
				max: 20,
			},
			validate: func(got int) bool {
				return got >= 10 && got <= 20
			},
			want: "val >= 10 && val <= 20",
		},
		{
			name: "Test Generate Random Int between 0 and 0",
			args: args{
				min: 0,
				max: 0,
			},
			validate: func(got int) bool {
				return got == 0
			},
			want: "val == 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Struct{}
			if got, _ := n.GenerateRandomInt(tt.args.min, tt.args.max); !tt.validate(got) {
				t.Errorf("Struct.GenerateRandomInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct_Round(t *testing.T) {
	type args struct {
		val       float64
		precision int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test Round 3.14159 to 2 decimal places",
			args: args{
				val:       3.14159,
				precision: 2,
			},
			want: 3.14,
		},
		{
			name: "Test Round 3.14159 to 3 decimal places",
			args: args{
				val:       3.14159,
				precision: 3,
			},
			want: 3.142,
		},
		{
			name: "Test Round 3.14159 to 0 decimal places",
			args: args{
				val:       3.14159,
				precision: 0,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Struct{}
			if got := n.Round(tt.args.val, tt.args.precision); got != tt.want {
				t.Errorf("Struct.Round() = %v, want %v", got, tt.want)
			}
		})
	}
}
