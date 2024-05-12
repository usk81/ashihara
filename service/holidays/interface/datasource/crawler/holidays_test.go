package crawler

import "testing"

func Test_convertHolidayName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "休日（祝日扱い）",
			args: args{
				s: "休日（祝日扱い）",
			},
			want: "休日",
		},
		{
			name: "体育の日（スポーツの日）",
			args: args{
				s: "体育の日（スポーツの日）",
			},
			want: "体育の日",
		},
		{
			name: "スポーツの日",
			args: args{
				s: "スポーツの日",
			},
			want: "スポーツの日",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertHolidayName(tt.args.s); got != tt.want {
				t.Errorf("convertHolidayName() = %v, want %v", got, tt.want)
			}
		})
	}
}
