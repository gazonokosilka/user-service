package dto

import (
	"errors"
	"strings"
	"time"
)

type CreateCustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Timezone  string `json:"timezone"`
	Birthday  string `json:"birthday"`
	UserID    string `json:"user_id"`
}

func (c *CreateCustomerRequest) Validate() error {
	if strings.TrimSpace(c.FirstName) == "" {
		return errors.New("first_name is required")
	}
	if len(c.FirstName) > 100 {
		return errors.New("first_name too long, max 100 characters")
	}
	if strings.TrimSpace(c.LastName) == "" {
		return errors.New("last_name is required")
	}
	if len(c.LastName) > 100 {
		return errors.New("last_name too long, max 100 characters")
	}
	if strings.TrimSpace(c.Gender) == "" {
		return errors.New("gender is required")
	}
	gender := strings.ToLower(strings.TrimSpace(c.Gender))
	if gender != "male" && gender != "female" {
		return errors.New("gender must be male or female")
	}
	if c.Timezone == "" {
		c.Timezone = "UTC"
	}
	if strings.TrimSpace(c.Birthday) == "" {
		return errors.New("birthday is required")
	}
	birthday, err := time.Parse("2006-01-02", strings.TrimSpace(c.Birthday))
	if err != nil {
		return errors.New("birthday must be in format YYYY-MM-DD (e.g., 1990-05-15)")
	}

	if birthday.After(time.Now()) {
		return errors.New("birthday cannot be in the future")
	}

	minDate := time.Now().AddDate(-150, 0, 0)
	if birthday.Before(minDate) {
		return errors.New("birthday is too far in the past")
	}
	if c.UserID == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func (c *CreateCustomerRequest) ParseBirthday() (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(c.Birthday))
}

type UpdateCustomerRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Gender    *string `json:"gender,omitempty"`
	Timezone  *string `json:"timezone,omitempty"`
	Birthday  *string `json:"birthday,omitempty"`
}

func (r *UpdateCustomerRequest) Validate() error {
	if r.FirstName != nil {
		if strings.TrimSpace(*r.FirstName) == "" {
			return errors.New("first_name cannot be empty")
		}
		if len(*r.FirstName) > 100 {
			return errors.New("first_name too long (max 100 characters)")
		}
	}
	if r.LastName != nil {
		if strings.TrimSpace(*r.LastName) == "" {
			return errors.New("last_name cannot be empty")
		}
		if len(*r.LastName) > 100 {
			return errors.New("last_name too long (max 100 characters)")
		}
	}

	if r.Gender != nil {
		if strings.TrimSpace(*r.Gender) == "" {
			return errors.New("gender cannot be empty")
		}
		gender := strings.ToLower(strings.TrimSpace(*r.Gender))
		if gender != "male" && gender != "female" {
			return errors.New("gender must be male or female")
		}
	}

	if r.Birthday != nil {
		birthday, err := time.Parse("2006-01-02", strings.TrimSpace(*r.Birthday))
		if err != nil {
			return errors.New("birthday must be in format YYYY-MM-DD")
		}
		if birthday.After(time.Now()) {
			return errors.New("birthday cannot be in the future")
		}
		minDate := time.Now().AddDate(-150, 0, 0)
		if birthday.Before(minDate) {
			return errors.New("birthday is too far in the past")
		}
	}

	return nil
}
func (r *UpdateCustomerRequest) ParseBirthday() (*time.Time, error) {
	if r.Birthday == nil {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", strings.TrimSpace(*r.Birthday))
	if err != nil {
		return nil, err
	}

	return &t, nil
}
