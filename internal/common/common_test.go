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
		want    *[]byte
		name    string
		fields  fields
		args    args
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
			wantRes: `dce08d478091ced02919c91029b1aaf0d55d44c348f1099684d75b9033032114`,
		},
		{
			name: "gauge",
			fields: fields{
				ID:    "x",
				MType: "gauge",
				Value: &testFloat,
			},
			args:    args{key: "abcdef"},
			wantRes: `af54a8b01e14ae89834b040ddc5177eec57b19a1c67a920e8ac6956f2346d579`,
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

func TestMetrics_CheckHash(t *testing.T) {
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
		wantErr bool
	}{
		{
			name: "check hash produced",
			fields: fields{
				ID:    "x",
				MType: "counter",
				Delta: &testInt,
			},
			args: args{key: "12345"},
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
			m.StoreHash(tt.args.key)
			if err := m.CheckHash(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Metrics.CheckHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
