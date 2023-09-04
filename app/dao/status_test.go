package dao_test

import (
	"context"
	"strings"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/stretchr/testify/assert"
)

func setupStatusDAO(t *testing.T) (repository.Status, func()) {
	dao, err := setupDAO()
	if err != nil {
		t.Fatal(err)
	}
	dao.InitAll()
	statusRepo := dao.Status()

	cleanup := func() {
		dao.InitAll()
	}

	return statusRepo, cleanup
}

func TestStatusCreate(t *testing.T) {
	repo, cleanup := setupStatusDAO(t)
	defer cleanup()

	ctx := context.Background()
	status := &object.Status{
		AccountId: 1,
		Content:   "Test Content",
	}
	err := repo.Create(ctx, status)
	assert.NoError(t, err)

	createdStatus, err := repo.Retrieve(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, createdStatus)
}

func TestStatusRetrieve(t *testing.T) {
	repo, cleanup := setupStatusDAO(t)
	defer cleanup()

	ctx := context.Background()

	status := &object.Status{
		AccountId: 1,
		Content:   "Test Content for Find",
	}
	err := repo.Create(ctx, status)
	assert.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		foundStatus, err := repo.Retrieve(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, foundStatus)
		assert.Equal(t, status.Content, foundStatus.Content)
	})

	t.Run("NotFound", func(t *testing.T) {
		foundStatus, err := repo.Retrieve(ctx, 2)
		assert.NoError(t, err)
		assert.Nil(t, foundStatus)
	})
}

func TestStatusDelete(t *testing.T) {
	repo, cleanup := setupStatusDAO(t)
	defer cleanup()

	ctx := context.Background()

	status := &object.Status{
		AccountId: 1,
		Content:   "Test Content for Delete",
	}
	err := repo.Create(ctx, status)
	assert.NoError(t, err)

	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(ctx, 1)
		assert.NoError(t, err)
		deletedStatus, err := repo.Retrieve(ctx, 1)
		assert.NoError(t, err)
		assert.Nil(t, deletedStatus)
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		err = repo.Delete(ctx, 2)
		assert.NoError(t, err)
	})
}

func TestPublicTimeline(t *testing.T) {
	repo, cleanup := setupStatusDAO(t)
	defer cleanup()

	ctx := context.Background()

	// テストのためのデータを作成
	for i := 1; i <= 10; i++ {
		status := &object.Status{
			AccountId: uint64(i),
			Content:   "Test Content " + strings.Repeat("#", i),
		}
		err := repo.Create(ctx, status)
		assert.NoError(t, err)
	}

	t.Run("All", func(t *testing.T) {
		allStatuses, err := repo.PublicTimeline(ctx, nil, nil, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 10, len(allStatuses))
	})

	t.Run("Limit", func(t *testing.T) {
		limit := uint64(5)
		limitedStatuses, err := repo.PublicTimeline(ctx, nil, nil, nil, &limit)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(limitedStatuses))
	})

	t.Run("SinceID", func(t *testing.T) {
		sinceID := uint64(5)
		sinceIDStatuses, err := repo.PublicTimeline(ctx, nil, nil, &sinceID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 6, len(sinceIDStatuses))
	})

	t.Run("MaxID", func(t *testing.T) {
		maxID := uint64(5)
		maxIDStatuses, err := repo.PublicTimeline(ctx, nil, &maxID, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(maxIDStatuses))
	})

	t.Run("SinceIDAndMaxID", func(t *testing.T) {
		sinceID := uint64(5)
		maxID := uint64(8)
		sinceIDAndMaxIDStatuses, err := repo.PublicTimeline(ctx, nil, &maxID, &sinceID, nil)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(sinceIDAndMaxIDStatuses))
	})
}
