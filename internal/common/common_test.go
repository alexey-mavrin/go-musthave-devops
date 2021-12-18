package common

import (
	"crypto/sha256"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testInt   int64             = 1
	hashInt   [sha256.Size]byte = sha256.Sum256([]byte("x:counter:1abcdef"))
	testFloat float64           = 1.1
	hashFloat [sha256.Size]byte = sha256.Sum256([]byte("x:gauge:1.100000abcdef"))
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
		want    *[sha256.Size]byte
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
		{
			name: "delta",
			fields: fields{
				ID:    "x",
				MType: "counter",
				Delta: &testInt,
			},
			args: args{key: "abcdef"},
			want: &hashInt,
		},
		{
			name: "gauge",
			fields: fields{
				ID:    "x",
				MType: "gauge",
				Value: &testFloat,
			},
			args: args{key: "abcdef"},
			want: &hashFloat,
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
			wantRes: `e679eeabe6696315863ef8518c3859f860a4ba52914902078c133d1c92d20a51`,
		},
		{
			name: "gauge",
			fields: fields{
				ID:    "x",
				MType: "gauge",
				Value: &testFloat,
			},
			args:    args{key: "abcdef"},
			wantRes: `0463a27af714a7631fddc3bf34a75ee0b0628b03f3a0f0cf9d7eb825e6b9af7b`,
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
