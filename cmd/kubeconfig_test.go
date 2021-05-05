package cmd

import (
	"reflect"
	"testing"
)

func Test_splitCommand(t *testing.T) {
	type args struct {
		commandline string
	}
	tests := []struct {
		name     string
		args     args
		wantCmd  string
		wantArgs []string
		wantErr  bool
	}{
		{
			name: "command with args",
			args: args{
				commandline: "/cmd -args1 -args2",
			},
			wantCmd:  "/cmd",
			wantArgs: []string{"-args1", "-args2"},
		},
		{
			name: "command only",
			args: args{
				commandline: "/cmd",
			},
			wantCmd:  "/cmd",
			wantArgs: []string{},
		},
		{
			name: "empty",
			args: args{
				commandline: "",
			},
			wantCmd:  "",
			wantArgs: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := splitCommand(tt.args.commandline)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantCmd {
				t.Errorf("splitCommand() got = %v, want %v", got, tt.wantCmd)
			}
			if !reflect.DeepEqual(got1, tt.wantArgs) {
				t.Errorf("splitCommand() got1 = %+v, want %+v", got1, tt.wantArgs)
			}
		})
	}
}
