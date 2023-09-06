package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"yatter-backend-go/app/domain/object"
)

func TestStatusCreate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		statuses []object.Status
	}{
		{
			name: "Success",
			statuses: []object.Status{{
				AccountId: 1, Content: "Test Content",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupDB()

			for _, status := range tt.statuses {
				err := statusRepo.Create(ctx, &status)
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusRetrieve(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		statuses []object.Status
		wantErr  bool
	}{
		{
			name: "Success",
			statuses: []object.Status{{
				AccountId: 1, Content: "Test Content",
			}},
			wantErr: false,
		},
		{
			name:     "NotFound",
			statuses: []object.Status{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupDB()

			for _, status := range tt.statuses {
				err := statusRepo.Create(ctx, &status)
				assert.NoError(t, err)
			}

			status, err := statusRepo.Retrieve(ctx, 1)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, status)
			}
		})
	}
}

func TestStatusDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		statuses []object.Status
		wantErr  bool
	}{
		{
			name: "Success",
			statuses: []object.Status{{
				AccountId: 1, Content: "Test Content",
			}},
			wantErr: false,
		},
		{
			name:     "NotFound",
			statuses: []object.Status{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupDB()

			for _, status := range tt.statuses {
				err := statusRepo.Create(ctx, &status)
				assert.NoError(t, err)
			}

			err := statusRepo.Delete(ctx, 1)
			if tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
				deletedStatus, err := statusRepo.Retrieve(ctx, 1)
				assert.Error(t, err)
				assert.Nil(t, deletedStatus)
			}
		})
	}
}

func TestPublicTimeline(t *testing.T) {
	ctx := context.Background()
	cleanupDB()

	for i := 1; i <= 10; i++ {
		status := &object.Status{
			AccountId: uint64(i),
			Content:   "Test Content " + strings.Repeat("#", i),
		}
		err := statusRepo.Create(ctx, status)
		assert.NoError(t, err)
	}

	tests := []struct {
		name      string
		expectLen int
		wantArgs  func() (only_media, max_id, since_id, limit *uint64)
	}{
		{
			name:      "All",
			expectLen: 10,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, nil, nil, nil
			},
		},
		{
			name:      "Limit",
			expectLen: 5,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, nil, nil, newUint64(5)
			},
		},
		{
			name:      "SinceID",
			expectLen: 6,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, nil, newUint64(5), nil
			},
		},
		{
			name:      "MaxID",
			expectLen: 5,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, newUint64(5), nil, nil
			},
		},
		{
			name:      "SinceIDAndMaxID",
			expectLen: 4,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, newUint64(8), newUint64(5), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			only_media, max_id, since_id, limit := tt.wantArgs()
			allStatuses, err := statusRepo.PublicTimeline(ctx, only_media, max_id, since_id, limit)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectLen, len(allStatuses))
		})
	}
}

func TestHomeTimeline(t *testing.T) {
	ctx := context.Background()
	cleanupDB()
	defer cleanupDB()
	followingID := uint64(1)
	insertRelationshipDB(t, ctx, []object.Relationship{
		{
			FollowingId: followingID,
			FollowerId:  2,
		},
	})
	for i := 1; i <= 10; i++ {
		status := &object.Status{
			AccountId: 2,
			Content:   "Test Content " + strings.Repeat("#", i),
		}
		err := statusRepo.Create(ctx, status)
		assert.NoError(t, err)
	}

	tests := []struct {
		name      string
		accountId uint64
		expectLen int
		wantArgs  func() (only_media, max_id, since_id, limit *uint64)
	}{
		{
			name:      "All",
			accountId: followingID,
			expectLen: 10,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, nil, nil, nil
			},
		},
		{
			name:      "Limit",
			accountId: followingID,
			expectLen: 5,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, nil, nil, newUint64(5)
			},
		},
		{
			name:      "SinceID",
			accountId: followingID,
			expectLen: 6,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, nil, newUint64(5), nil
			},
		},
		{
			name:      "MaxID",
			accountId: followingID,
			expectLen: 5,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, newUint64(5), nil, nil
			},
		},
		{
			name:      "SinceIDAndMaxID",
			accountId: followingID,
			expectLen: 4,
			wantArgs: func() (only_media, max_id, since_id, limit *uint64) {
				return nil, newUint64(8), newUint64(5), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			only_media, max_id, since_id, limit := tt.wantArgs()
			allStatuses, err := statusRepo.HomeTimeline(ctx, tt.accountId, only_media, max_id, since_id, limit)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectLen, len(allStatuses))
		})
	}
}

func newUint64(i uint64) *uint64 {
	return &i
}
