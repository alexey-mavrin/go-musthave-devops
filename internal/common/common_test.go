package common

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testInt   int64   = 1
	testFloat float64 = 1.1
)

func TestMetrics_ComputeHash(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
		Hash  string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]byte
		wantErr bool
	}{
		{
			name:    "empty key",
			args:    args{key: ""},
			wantErr: true,
		},
		{
			name: "empty id",
			fields: fields{
				ID: "",
			},
			args:    args{key: "abcdef"},
			wantErr: true,
		},
		{
			name: "nil value",
			fields: fields{
				ID:    "x",
				MType: "gauge",
				Value: nil,
			},
			args:    args{key: "abcdef"},
			wantErr: true,
		},
		{
			name: "nil delta",
			fields: fields{
				ID:    "x",
				MType: "counter",
				Delta: nil,
			},
			args:    args{key: "abcdef"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
				Hash:  tt.fields.Hash,
			}
			got, err := m.ComputeHash(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metrics.ComputeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Metrics.ComputeHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_StoreHash(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
		Hash  string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes string
		wantErr bool
	}{
		{
			name: "delta",
			fields: fields{
				ID:    "x",
				MType: "counter",
				Delta: &testInt,
			},
			args:    args{key: "abcdef"},
			wantRes: `4dfa5caf0f7bce10f304ded32e9a680341c87bbd8de913966c3f31cf85cd47cb`,
		},
		{
			name: "gauge",
			fields: fields{
				ID:    "x",
				MType: "gauge",
				Value: &testFloat,
			},
			args:    args{key: "abcdef"},
			wantRes: `8d81fbfecf9f7efe52d8bd783e005cb90e4f139c8f5d28e2e6b095d18c2645e2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
				Hash:  tt.fields.Hash,
			}
			if err := m.StoreHash(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Metrics.StoreHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantRes, m.Hash)
		})
	}
}
