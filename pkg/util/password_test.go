package util

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "password",
			args: args{password: "123456"},
			want: "$2a$08$3RuMTEk5ZaCslUa6yWBO6ecx00I7zB84Mpny2N.r.o0SVorI0iQjC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashPassword(tt.args.password); got != tt.want {
				fmt.Println(got)

				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
