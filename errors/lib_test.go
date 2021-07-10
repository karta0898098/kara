package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestExceptionBuild(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "ResourceNotFound",
			err:  ErrResourceNotFound.BuildWithError(gorm.ErrRecordNotFound),
		},
		{
			name: "ResourceNotFoundWithError",
			err: ErrResourceNotFound.WithDetails(
				DetailData{
					"domain": "article",
					"reason": "id not found",
				}).BuildWithError(gorm.ErrRecordNotFound),
		},
		{
			name: "ResourceNotFoundWithMsg",
			err: ErrResourceNotFound.WithDetails(
				DetailData{
					"domain": "article",
					"reason": "id not found",
				}).Build("occur error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("%+v \n", tt.err)
		})
	}
}

func TestExceptionToRestfulView(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    error
		view   *exceptionView
	}{
		{
			name:   "ResourceNotFound",
			status: 404,
			err:    ErrResourceNotFound.BuildWithError(gorm.ErrRecordNotFound),
			view: &exceptionView{
				Code:    404002,
				Message: "The specified resource does not exist.",
				Details: nil,
			},
		},
		{
			name:   "ResourceNotFoundWithDetails",
			status: 404,
			err: ErrResourceNotFound.WithDetails(
				DetailData{
					"domain": "article",
					"reason": "id not found",
				},
			).BuildWithError(gorm.ErrRecordNotFound),
			view: &exceptionView{
				Code:    404002,
				Message: "The specified resource does not exist.",
				Details: DetailData{
					"domain": "article",
					"reason": "id not found",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TryConvert(tt.err)
			status, payload := err.ToViewModel()
			assert.Equal(t, tt.status, status)
			assert.Equal(t, tt.view, payload)
		})
	}
}

func TestExceptionIs(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "ResourceNotFound",
			err:  ErrResourceNotFound.BuildWithError(gorm.ErrRecordNotFound),
		},
		{
			name: "ResourceNotFoundWithError",
			err: ErrResourceNotFound.WithDetails(
				DetailData{
					"domain": "article",
					"reason": "id not found",
				}).BuildWithError(gorm.ErrRecordNotFound),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, true, ErrResourceNotFound.Is(tt.err))
		})
	}
}
