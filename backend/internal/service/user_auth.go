package service

import (
	"errors"
	"time"

	"github.com/nexus-id/backend/internal/models"
	"gorm.io/gorm"
)

var (
	MaxFailedAttempts      = 5
	DefaultLockoutDuration = 30 * time.Minute
)

// LockUserAccount locks a user account
func (s *UserService) LockUserAccount(userID int64, duration time.Duration) error {
	lockedUntil := time.Now().Add(duration)
	return s.userRepo.LockAccount(userID, lockedUntil)
}

// UnlockUserAccount unlocks a user account
func (s *UserService) UnlockUserAccount(userID int64) error {
	return s.userRepo.LockAccount(userID, nil)
}

// CheckAndLockUser checks if user should be locked and locks if necessary
func (s *UserService) CheckAndLockUser(user *models.User) error {
	// Already locked
	if user.IsLocked() {
		return ErrAccountLocked
	}

	// Check if threshold exceeded
	if user.FailedLoginAttempts >= MaxFailedAttempts {
		// Lock the account
		if err := s.LockUserAccount(user.ID, DefaultLockoutDuration); err != nil {
			return err
		}
		return ErrAccountLocked
	}

	return nil
}

// GetUserModelForAuth retrieves user model for authentication
func (s *UserService) GetUserModelForAuth(tenantID int64, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(tenantID, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
