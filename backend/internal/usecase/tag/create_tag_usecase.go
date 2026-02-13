package tag

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
)

type CreateTagUseCase interface {
	Execute(ctx context.Context, name string) (*model.Tag, error)
}

type createTagInteractor struct {
	tagRepo repository.TagRepository
}

func NewCreateTagUseCase(tagRepo repository.TagRepository) CreateTagUseCase {
	return &createTagInteractor{tagRepo: tagRepo}
}

func (uc *createTagInteractor) Execute(ctx context.Context, name string) (*model.Tag, error) {
	if err := validation.ValidateTag(name); err != nil {
		return nil, err
	}

	// Check if tag already exists
	existingTag, err := uc.tagRepo.GetByName(ctx, name)
	if err == nil && existingTag != nil {
		return existingTag, nil // Return existing tag
	}

	tag := &model.Tag{Name: name}
	if err := uc.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}
