package tag

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type CreateTagUseCase struct {
	tagRepo repository.TagRepository
}

func NewCreateTagUseCase(tagRepo repository.TagRepository) *CreateTagUseCase {
	return &CreateTagUseCase{tagRepo: tagRepo}
}

func (uc *CreateTagUseCase) Execute(ctx context.Context, name string) (*entity.Tag, error) {
	// Check if tag already exists
	existingTag, err := uc.tagRepo.GetByName(ctx, name)
	if err == nil && existingTag != nil {
		return existingTag, nil // Return existing tag
	}

	tag := &entity.Tag{Name: name}
	if err := uc.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}
