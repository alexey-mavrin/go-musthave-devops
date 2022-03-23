package main

import (
	"testing"
	"time"

	"github.com/alexey-mavrin/go-musthave-devops/internal/agent"
	"github.com/stretchr/testify/assert"
)

func Test_setAgentArgs(t *testing.T) {
	type want struct {
		ServerAddr     string
		PollInterval   time.Duration
		ReportInterval time.Duration
		Key            string
		useJSON        bool
		useBatch       bool
	}
	tests := []struct {
		name    string
		wantErr bool
		want    want
	}{
		{
			name:    "setAgentArgs makes default config",
			wantErr: false,
			want: want{
				ServerAddr:     "http://localhost:8080",
				PollInterval:   2 * time.Second,
				ReportInterval: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setAgentArgs(); (err != nil) != tt.wantErr {
				t.Errorf("setAgentArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t,
				tt.want.ServerAddr,
				agent.Config.ServerAddr)
			assert.Equal(t,
				tt.want.PollInterval,
				agent.Config.PollInterval)
			assert.Equal(t,
				tt.want.ReportInterval,
				agent.Config.ReportInterval)
			assert.Equal(t,
				tt.want.Key,
				agent.Config.Key)
		})
	}
}
