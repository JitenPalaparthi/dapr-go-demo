package interfaces

import (
	"context"

	"github.com/JitenPalaparthi/atipaday/models"
)

type ITip interface {
	Create(ctx context.Context, tip *models.Tip) (*models.Tip, error)
	UpdateBy(ctx context.Context, id string, data map[string]any) (*models.Tip, error)
	GetBy(ctx context.Context, id string) (*models.Tip, error)
	GetAllByOffSet(ctx context.Context, offset, limit int) ([]models.Tip, error)
	DeleteBy(ctx context.Context, id string) (any, error)
	Search(ctx context.Context, offset, limit int, search string) ([]models.Tip, error)
}
