package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"yatter-backend-go/app/domain/object"
)

func TestAccountCreate(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		accounts  []object.Account
		wantError bool
	}{
		{
			name: "Success",
			accounts: []object.Account{{
				Username: "testuser", PasswordHash: "hashedpassword",
			}},
			wantError: false,
		},
		{
			name: "Duplicate username",
			accounts: []object.Account{{
				Username:     "testuser",
				PasswordHash: "hashedpassword",
			}, {
				Username:     "testuser",
				PasswordHash: "hashedpassword",
			}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupDB()

			for i, account := range tt.accounts {
				err := accountRepo.Create(ctx, &account)
				if i == len(tt.accounts)-1 {
					if tt.wantError {
						assert.Error(t, err)
					}
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestAccountRetrieve(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		accounts []object.Account
		wantErr  bool
	}{
		{
			name: "Success",
			accounts: []object.Account{{
				Username: "testuser", PasswordHash: "hashedpassword",
			}},
			wantErr: false,
		},
		{
			name:    "Not found",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupDB()

			for _, account := range tt.accounts {
				err := accountRepo.Create(ctx, &account)
				assert.NoError(t, err)
			}

			account, err := accountRepo.Retrieve(ctx, "testuser")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
			}
		})
	}
}
