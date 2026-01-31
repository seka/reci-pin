package tag

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type DeleteTagUseCase struct {
	tagRepo repository.TagRepository
}

func NewDeleteTagUseCase(tagRepo repository.TagRepository) *DeleteTagUseCase {
	return &DeleteTagUseCase{tagRepo: tagRepo}
}

func (uc *DeleteTagUseCase) Execute(ctx context.Context, id int64) error {
	return uc.tagRepo.Delete(ctx, id)
}
