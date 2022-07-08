package main

import (
	"strings"
	"testing"
)

func TestArguments_Description(t *testing.T) {
	type fields struct {
		Server *Server
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "normal",
			fields: fields{},
			want:   "security tool that manages the vulnerability",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := Arguments{
				Server: tt.fields.Server,
			}
			if got := ar.Description(); !strings.Contains(got, tt.want) {
				t.Errorf("Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArguments_Version(t *testing.T) {
	type fields struct {
		Server *Server
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "normal",
			fields: fields{},
			want:   "CompilerVersion:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := Arguments{
				Server: tt.fields.Server,
			}
			if got := ar.Version(); !strings.Contains(got, tt.want) {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
