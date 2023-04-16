package embedded

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestEmbedStorage_Delete(t *testing.T) {
	type fields struct {
		items           map[string]item
		cleanupInterval time.Duration
		stopCleaning    chan bool
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
			name: "not existing",
			args: args{
				key: "not-existing",
			},
			wantErr: false,
		},
		{
			name: "existing",
			args: args{
				key: "existing",
			},
			fields: fields{
				items: map[string]item{
					"existing": {
						Value: []byte("test"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EmbeddedStorage{
				items:           tt.fields.items,
				cleanupInterval: tt.fields.cleanupInterval,
				stopCleaning:    tt.fields.stopCleaning,
			}

			err := s.Delete(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.NoError(t, err)
		})
	}
}

func TestEmbedStorage_Get(t *testing.T) {
	type fields struct {
		items           map[string]item
		cleanupInterval time.Duration
		stopCleaning    chan bool
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "not found",
			args: args{
				key: "not-existing-key",
			},
			want: nil,
		},
		{
			name: "existing",
			args: args{
				key: "existing-key",
			},
			fields: fields{
				items: map[string]item{
					"existing-key": {
						Value: []byte("test"),
					},
				},
			},
			want: []byte("test"),
		},
		{
			name: "expired",
			args: args{
				key: "expired-key",
			},
			fields: fields{
				items: map[string]item{
					"expired-key": {
						Value:      []byte("test"),
						Expiration: time.Now().Add(-time.Second * 10).UnixNano(),
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EmbeddedStorage{
				items:           tt.fields.items,
				cleanupInterval: tt.fields.cleanupInterval,
				stopCleaning:    tt.fields.stopCleaning,
			}

			got, err := s.Get(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEmbedStorage_Set(t *testing.T) {
	type fields struct {
		items           map[string]item
		cleanupInterval time.Duration
		mx              sync.RWMutex
		stopCleaning    chan bool
	}
	type args struct {
		key        string
		value      []byte
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add item",
			fields: fields{
				items: make(map[string]item),
			},
			args: args{
				key:        "key",
				value:      []byte("test"),
				expiration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EmbeddedStorage{
				items:           tt.fields.items,
				cleanupInterval: tt.fields.cleanupInterval,
				stopCleaning:    tt.fields.stopCleaning,
			}
			assert.NoError(t, s.Set(tt.args.key, tt.args.value, tt.args.expiration))
		})
	}
}

func TestEmbedStorage_cleaner(t *testing.T) {
	type fields struct {
		items           map[string]item
		cleanupInterval time.Duration
		stopCleaning    chan bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "cleanup",
			fields: fields{
				cleanupInterval: 20 * time.Millisecond,
				stopCleaning:    make(chan bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EmbeddedStorage{
				items:           tt.fields.items,
				cleanupInterval: tt.fields.cleanupInterval,
				stopCleaning:    tt.fields.stopCleaning,
			}
			go s.cleaner()
			time.Sleep(100 * time.Millisecond)
			finalizer(s)
		})
	}
}

func TestEmbedStorage_deleteExpired(t *testing.T) {
	type fields struct {
		items           map[string]item
		cleanupInterval time.Duration
		mx              sync.RWMutex
		stopCleaning    chan bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "delete expired",
			fields: fields{
				items: map[string]item{
					"expired": {
						Value:      []byte("test-1"),
						Expiration: time.Now().Add(-time.Second * 10).UnixNano(),
					},
					"not-expired": {
						Value: []byte("test-2"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EmbeddedStorage{
				items:           tt.fields.items,
				cleanupInterval: tt.fields.cleanupInterval,
				stopCleaning:    tt.fields.stopCleaning,
			}
			s.deleteExpired()
			assert.Equal(t, map[string]item{
				"not-expired": {
					Value: []byte("test-2"),
				},
			}, s.items)
		})
	}
}

func TestNewEmbedStorage(t *testing.T) {
	type args struct {
		cleanupInterval time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creation",
			args: args{
				cleanupInterval: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewEmbeddedStorage(tt.args.cleanupInterval))
		})
	}
}

func Test_item_IsExpired(t *testing.T) {
	type fields struct {
		Value      []byte
		Expiration int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "not expired",
			fields: fields{
				Value:      []byte("test"),
				Expiration: time.Now().Add(time.Second * 10).UnixNano(),
			},
			want: false,
		},
		{
			name: "expired",
			fields: fields{
				Value:      []byte("test"),
				Expiration: time.Now().Add(-time.Second * 10).UnixNano(),
			},
			want: true,
		},
		{
			name: "eternal",
			fields: fields{
				Value:      []byte("test"),
				Expiration: 0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &item{
				Value:      tt.fields.Value,
				Expiration: tt.fields.Expiration,
			}

			assert.Equal(t, tt.want, i.IsExpired())
		})
	}
}
