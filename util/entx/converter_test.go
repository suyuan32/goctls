package entx

import "testing"

func TestIsIDField(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 1",
			args: args{
				field: "Id",
			},
			want: false,
		},
		{
			name: "test 2",
			args: args{
				field: "Ids",
			},
			want: true,
		},
		{
			name: "test 3",
			args: args{
				field: "Idx",
			},
			want: true,
		},
		{
			name: "test 4",
			args: args{
				field: "UserId",
			},
			want: false,
		},
		{
			name: "test 5",
			args: args{
				field: "IdxId",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotIDField(tt.args.field); got != tt.want {
				t.Errorf("IsIDField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertIdFieldToUpper(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				target: "Id",
			},
			want: "ID",
		},
		{
			name: "test 2",
			args: args{
				target: "UserId",
			},
			want: "UserID",
		},
		{
			name: "test 3",
			args: args{
				target: "Identity",
			},
			want: "Identity",
		},
		{
			name: "test 4",
			args: args{
				target: "Idx",
			},
			want: "Idx",
		},
		{
			name: "test 5",
			args: args{
				target: "IdxId",
			},
			want: "IdxID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertIdFieldToUpper(tt.args.target); got != tt.want {
				t.Errorf("ConvertIdFieldToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertSpecificNounToUpper(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				str: "IdentityUuid",
			},
			want: "IdentityUUID",
		},
		{
			name: "test 2",
			args: args{
				str: "IdxApi",
			},
			want: "IdxAPI",
		},
		{
			name: "test 3",
			args: args{
				str: "IdxApi",
			},
			want: "IdxAPI",
		},
		{
			name: "test 4",
			args: args{
				str: "TestUrl",
			},
			want: "TestURL",
		},
		{
			name: "test 5",
			args: args{
				str: "IdxId",
			},
			want: "IdxID",
		},
		{
			name: "test 6",
			args: args{
				str: "IdxIds",
			},
			want: "IdxIds",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertSpecificNounToUpper(tt.args.str); got != tt.want {
				t.Errorf("ConvertSpecificNounToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}
