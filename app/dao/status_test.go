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
	cleanupDB()

	status := &object.Status{
		AccountId: 1,
		Content:   "Test Content",
	}
	err := statusRepo.Create(ctx, status)
	assert.NoError(t, err)

	createdStatus, err := statusRepo.Retrieve(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, createdStatus)
}

func TestStatusRetrieve(t *testing.T) {
	ctx := context.Background()
	cleanupDB()

	status := &object.Status{
		AccountId: 1,
		Content:   "Test Content for Find",
	}
	err := statusRepo.Create(ctx, status)
	assert.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		foundStatus, err := statusRepo.Retrieve(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, foundStatus)
		assert.Equal(t, status.Content, foundStatus.Content)
	})

	t.Run("NotFound", func(t *testing.T) {
		foundStatus, err := statusRepo.Retrieve(ctx, 2)
		assert.NoError(t, err)
		assert.Nil(t, foundStatus)
	})
}

func TestStatusDelete(t *testing.T) {
	ctx := context.Background()
	cleanupDB()

	status := &object.Status{
		AccountId: 1,
		Content:   "Test Content for Delete",
	}
	err := statusRepo.Create(ctx, status)
	assert.NoError(t, err)

	t.Run("Delete", func(t *testing.T) {
		err = statusRepo.Delete(ctx, 1)
		assert.NoError(t, err)
		deletedStatus, err := statusRepo.Retrieve(ctx, 1)
		assert.NoError(t, err)
		assert.Nil(t, deletedStatus)
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		err = statusRepo.Delete(ctx, 2)
		assert.NoError(t, err)
	})
}

func TestPublicTimeline(t *testing.T) {
	ctx := context.Background()
	cleanupDB()

	// テストのためのデータを作成
	for i := 1; i <= 10; i++ {
		status := &object.Status{
			AccountId: uint64(i),
			Content:   "Test Content " + strings.Repeat("#", i),
		}
		err := statusRepo.Create(ctx, status)
		assert.NoError(t, err)
	}

	t.Run("All", func(t *testing.T) {
		allStatuses, err := statusRepo.PublicTimeline(ctx, nil, nil, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 10, len(allStatuses))
	})

	t.Run("Limit", func(t *testing.T) {
		limit := uint64(5)
		limitedStatuses, err := statusRepo.PublicTimeline(ctx, nil, nil, nil, &limit)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(limitedStatuses))
	})

	t.Run("SinceID", func(t *testing.T) {
		sinceID := uint64(5)
		sinceIDStatuses, err := statusRepo.PublicTimeline(ctx, nil, nil, &sinceID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 6, len(sinceIDStatuses))
	})

	t.Run("MaxID", func(t *testing.T) {
		maxID := uint64(5)
		maxIDStatuses, err := statusRepo.PublicTimeline(ctx, nil, &maxID, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(maxIDStatuses))
	})

	t.Run("SinceIDAndMaxID", func(t *testing.T) {
		sinceID := uint64(5)
		maxID := uint64(8)
		sinceIDAndMaxIDStatuses, err := statusRepo.PublicTimeline(ctx, nil, &maxID, &sinceID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(sinceIDAndMaxIDStatuses))
	})
}
