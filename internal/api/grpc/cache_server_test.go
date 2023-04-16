package grpc

import (
	"context"
	"errors"
	v1 "github.com/dimuska139/cacher/internal/api/grpc/gen/cacher/cache/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheServer_Delete(t *testing.T) {
	type fields struct {
		logger  Logger
		storage Storage
	}
	type args struct {
		ctx     context.Context
		request *v1.DeleteRequest
	}

	tests := []struct {
		name      string
		getFields func(storage *MockStorage, logger *MockLogger) fields
		args      args
		want      *v1.DeleteResponse
		wantErr   bool
	}{
		{
			name: "without error",
			getFields: func(mockedStorage *MockStorage, _ *MockLogger) fields {
				mockedStorage.EXPECT().
					Delete("key").
					Return(nil).
					Times(1)
				return fields{
					storage: mockedStorage,
					logger:  nil,
				}
			},
			args: args{
				ctx: context.Background(),
				request: &v1.DeleteRequest{
					Key: "key",
				},
			},
			want:    &v1.DeleteResponse{},
			wantErr: false,
		},
		{
			name: "with error",
			getFields: func(mockedStorage *MockStorage, mockedLogger *MockLogger) fields {
				err := errors.New("error")

				mockedLogger.EXPECT().
					Error("Can't delete data from storage", "err", err, "key", "key").
					Times(1)

				mockedStorage.EXPECT().
					Delete("key").
					Return(err).
					Times(1)
				return fields{
					storage: mockedStorage,
					logger:  mockedLogger,
				}
			},
			args: args{
				ctx: context.Background(),
				request: &v1.DeleteRequest{
					Key: "key",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockedStorage := NewMockStorage(ctrl)
			mockedLogger := NewMockLogger(ctrl)

			mockedFields := tt.getFields(mockedStorage, mockedLogger)

			s := &CacheServer{
				logger:  mockedFields.logger,
				storage: mockedFields.storage,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCacheServer_Get(t *testing.T) {
	type fields struct {
		logger  Logger
		storage Storage
	}
	type args struct {
		ctx     context.Context
		request *v1.GetRequest
	}

	tests := []struct {
		name      string
		getFields func(storage *MockStorage, logger *MockLogger) fields
		args      args
		want      *v1.GetResponse
		wantErr   bool
	}{
		{
			name: "without error",
			getFields: func(mockedStorage *MockStorage, _ *MockLogger) fields {
				mockedStorage.EXPECT().
					Get("key").
					Return([]byte("data"), nil).
					Times(1)
				return fields{
					storage: mockedStorage,
					logger:  nil,
				}
			},
			args: args{
				ctx: context.Background(),
				request: &v1.GetRequest{
					Key: "key",
				},
			},
			want: &v1.GetResponse{
				Value: []byte("data"),
			},
			wantErr: false,
		},
		{
			name: "with error",
			getFields: func(mockedStorage *MockStorage, mockedLogger *MockLogger) fields {
				err := errors.New("error")

				mockedLogger.EXPECT().
					Error("Can't get data from storage", "err", err, "key", "key").
					Times(1)

				mockedStorage.EXPECT().
					Get("key").
					Return(nil, err).
					Times(1)
				return fields{
					storage: mockedStorage,
					logger:  mockedLogger,
				}
			},
			args: args{
				ctx: context.Background(),
				request: &v1.GetRequest{
					Key: "key",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockedStorage := NewMockStorage(ctrl)
			mockedLogger := NewMockLogger(ctrl)

			mockedFields := tt.getFields(mockedStorage, mockedLogger)

			s := &CacheServer{
				logger:  mockedFields.logger,
				storage: mockedFields.storage,
			}
			got, err := s.Get(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCacheServer_Set(t *testing.T) {
	type fields struct {
		logger  Logger
		storage Storage
	}
	type args struct {
		ctx     context.Context
		request *v1.SetRequest
	}

	tests := []struct {
		name      string
		getFields func(storage *MockStorage, logger *MockLogger) fields
		args      args
		want      *v1.SetResponse
		wantErr   bool
	}{
		{
			name: "without error",
			getFields: func(mockedStorage *MockStorage, _ *MockLogger) fields {
				mockedStorage.EXPECT().
					Set("key", []byte("test"), time.Second*time.Duration(10)).
					Return(nil).
					Times(1)
				return fields{
					storage: mockedStorage,
					logger:  nil,
				}
			},
			args: args{
				ctx: context.Background(),
				request: &v1.SetRequest{
					Key:   "key",
					Value: []byte("test"),
					Ttl:   uint64(10),
				},
			},
			want:    &v1.SetResponse{},
			wantErr: false,
		},
		{
			name: "with error",
			getFields: func(mockedStorage *MockStorage, mockedLogger *MockLogger) fields {
				err := errors.New("error")

				mockedLogger.EXPECT().
					Error("Can't get save data to storage", "err", err).
					Times(1)

				mockedStorage.EXPECT().
					Set("key", []byte("test"), time.Second*time.Duration(10)).
					Return(err).
					Times(1)
				return fields{
					storage: mockedStorage,
					logger:  mockedLogger,
				}
			},
			args: args{
				ctx: context.Background(),
				request: &v1.SetRequest{
					Key:   "key",
					Value: []byte("test"),
					Ttl:   uint64(10),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockedStorage := NewMockStorage(ctrl)
			mockedLogger := NewMockLogger(ctrl)

			mockedFields := tt.getFields(mockedStorage, mockedLogger)

			s := &CacheServer{
				logger:  mockedFields.logger,
				storage: mockedFields.storage,
			}
			got, err := s.Set(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewCacheServer(t *testing.T) {
	type args struct {
		logger  Logger
		storage Storage
	}

	ctrl := gomock.NewController(t)
	mockedLogger := NewMockLogger(ctrl)
	mockedStorage := NewMockStorage(ctrl)

	tests := []struct {
		name string
		args args
		want *CacheServer
	}{
		{
			name: "creation",
			args: args{
				logger:  mockedLogger,
				storage: mockedStorage,
			},
			want: &CacheServer{
				logger:  mockedLogger,
				storage: mockedStorage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewCacheServer(tt.args.logger, tt.args.storage))
		})
	}
}
