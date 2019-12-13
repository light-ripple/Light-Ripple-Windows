package common

import "testing"

func TestPaginate(t *testing.T) {
	type args struct {
		page     string
		limit    string
		maxLimit int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1",
			args{
				"10",
				"",
				100,
			},
			" LIMIT 450,50 ",
		},
		{
			"2",
			args{
				"-5",
				"-15",
				100,
			},
			" LIMIT 0,50 ",
		},
		{
			"3",
			args{
				"2",
				"150",
				100,
			},
			" LIMIT 100,100 ",
		},
	}
	for _, tt := range tests {
		if got := Paginate(tt.args.page, tt.args.limit, tt.args.maxLimit); got != tt.want {
			t.Errorf("%q. Paginate() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
