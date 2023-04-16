package memcache

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestMemcacheStorage_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name              string
		getMemcacheClient func() Memcacher
		args              args
		wantErr           bool
	}{
		{
			name: "with error",
			getMemcacheClient: func() Memcacher {
				ctrl := gomock.NewController(t)
				mockedClient := NewMockMemcacher(ctrl)
				mockedClient.EXPECT().
					Delete("testkey").
					Return(errors.New("something went wrong")).
					Times(1)
				return mockedClient
			},
			args: args{
				key: "testkey",
			},
			wantErr: true,
		},
		{
			name: "without error",
			getMemcacheClient: func() Memcacher {
				ctrl := gomock.NewController(t)
				mockedClient := NewMockMemcacher(ctrl)
				mockedClient.EXPECT().
					Delete("testkey").
					Return(nil).
					Times(1)
				return mockedClient
			},
			args: args{
				key: "testkey",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemcacheStorage{
				memcacheClient: tt.getMemcacheClient(),
			}

			if tt.wantErr {
				assert.Error(t, s.Delete(tt.args.key))
			} else {
				assert.NoError(t, s.Delete(tt.args.key))
			}
		})
	}
}

func TestMemcacheStorage_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name              string
		getMemcacheClient func() Memcacher
		args              args
		want              []byte
		wantErr           bool
	}{
		{
			name: "with error",
			getMemcacheClient: func() Memcacher {
				ctrl := gomock.NewController(t)
				mockedClient := NewMockMemcacher(ctrl)
				mockedClient.EXPECT().
					Get("testkey").
					Return(nil, errors.New("something went wrong")).
					Times(1)
				return mockedClient
			},
			args: args{
				key: "testkey",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "without error",
			getMemcacheClient: func() Memcacher {
				ctrl := gomock.NewController(t)
				mockedClient := NewMockMemcacher(ctrl)
				mockedClient.EXPECT().
					Get("testkey").
					Return([]byte("data"), nil).
					Times(1)
				return mockedClient
			},
			args: args{
				key: "testkey",
			},
			want:    []byte("data"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemcacheStorage{
				memcacheClient: tt.getMemcacheClient(),
			}
			got, err := s.Get(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMemcacheStorage_Set(t *testing.T) {
	type args struct {
		key   string
		value []byte
		ttl   time.Duration
	}
	tests := []struct {
		name              string
		getMemcacheClient func() Memcacher
		args              args
		wantErr           bool
	}{
		{
			name: "with error",
			getMemcacheClient: func() Memcacher {
				ctrl := gomock.NewController(t)
				mockedClient := NewMockMemcacher(ctrl)
				mockedClient.EXPECT().
					Set("testkey", []byte("data"), int64((time.Second * 5).Seconds())).
					Return(errors.New("something went wrong")).
					Times(1)
				return mockedClient
			},
			args: args{
				key:   "testkey",
				value: []byte("data"),
				ttl:   time.Second * 5,
			},
			wantErr: true,
		},
		{
			name: "without error",
			getMemcacheClient: func() Memcacher {
				ctrl := gomock.NewController(t)
				mockedClient := NewMockMemcacher(ctrl)
				mockedClient.EXPECT().
					Set("testkey", []byte("data"), int64((time.Second * 5).Seconds())).
					Return(nil).
					Times(1)
				return mockedClient
			},
			args: args{
				key:   "testkey",
				value: []byte("data"),
				ttl:   time.Second * 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemcacheStorage{
				memcacheClient: tt.getMemcacheClient(),
			}

			if tt.wantErr {
				assert.Error(t, s.Set(tt.args.key, tt.args.value, tt.args.ttl))
			} else {
				assert.NoError(t, s.Set(tt.args.key, tt.args.value, tt.args.ttl))
			}
		})
	}
}

func TestNewMemcacheStorage(t *testing.T) {
	type args struct {
		memcacheClient Memcacher
	}

	ctrl := gomock.NewController(t)
	mockedClient := NewMockMemcacher(ctrl)
	tests := []struct {
		name string
		args args
		want *MemcacheStorage
	}{
		{
			name: "creation",
			args: args{
				memcacheClient: mockedClient,
			},
			want: &MemcacheStorage{
				memcacheClient: mockedClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemcacheStorage(tt.args.memcacheClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemcacheStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
