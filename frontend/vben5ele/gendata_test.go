// Copyright (C) 2023  Ryan SU (https://github.com/suyuan32)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package vben5ele

import (
	"reflect"
	"testing"
)

func TestConvertTagToRules(t *testing.T) {
	type args struct {
		tagString string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "testString",
			args:    args{"omitempty,min=1,max=50"},
			want:    []string{"min: 1", "max: 50"},
			wantErr: false,
		},
		{
			name:    "testLen",
			args:    args{"omitempty,len=50"},
			want:    []string{"len: 50"},
			wantErr: false,
		},
		{
			name:    "testNum",
			args:    args{"omitempty,gte=1,lte=50"},
			want:    []string{"min: 1", "max: 50"},
			wantErr: false,
		},
		{
			name:    "testFloat",
			args:    args{"omitempty,gte=1.1,lte=50.0"},
			want:    []string{"min: 1.1", "max: 50.0"},
			wantErr: false,
		},
		{
			name:    "testFloat2",
			args:    args{"omitempty,gt=1.11,lt=50.01"},
			want:    []string{"min: 1.12", "max: 50.00"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTagToRules(tt.args.tagString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTagToRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTagToRules() got = %v, want %v", got, tt.want)
			}
		})
	}
}
