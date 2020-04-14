package api

import (
	"testing"
)

func TestGetFiles(t *testing.T) {
	type args struct {
		path      string
		auth      string
		rootDrive string
	}
	tests := []struct {
		name string
		args args
		want []*File
	}{
		{
			name: "Test",
			args: args{
				path:      "Music/",
				auth:      "",
				rootDrive: "XXXX",
			},
			want: nil,
		},
		{
			name: "Test",
			args: args{
				path:      "/", // root
				auth:      "",
				rootDrive: "XXXX",
			},
			want: nil,
		},
		{
			name: "Test",
			args: args{
				path:      "Music", // is WRONG!
				auth:      "",
				rootDrive: "XXXX",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := ShowDir(tt.args.path, tt.args.auth, tt.args.rootDrive); got == nil {
				t.Errorf("ShowDir() = %v, want %v", got, tt.want)
			} else {
				for _, f := range got {
					t.Log(f.Id, f.Name, f.Ext, "IsFolder:", f.IsFolder, "IsPlayable:", f.IsPlayable)
				}
			}
		})
	}
}
