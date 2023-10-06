package database

import (
	"context"
	"strconv"

	"github.com/JitenPalaparthi/atipaday/models"

	"gorm.io/gorm"
)

type Tip struct {
	DB any
}

func (c *Tip) Create(ctx context.Context, Tip *models.Tip) (*models.Tip, error) {
	// Need to check this implementation. What if second time it is called
	if err := c.DB.(*gorm.DB).AutoMigrate(&models.Tip{}); err != nil {
		return nil, err
	}
	tx := c.DB.(*gorm.DB).Create(Tip)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return Tip, nil
}

func (c *Tip) UpdateBy(ctx context.Context, id string, data map[string]any) (*models.Tip, error) {
	Tip := new(models.Tip)
	_id, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	Tip.ID = uint(_id)
	tx := c.DB.(*gorm.DB).Model(Tip).Updates(data)
	if tx.Error != nil {
		return nil, tx.Error
	}
	Tip, err = c.GetBy(ctx, id)
	if err != nil {
		return nil, err
	}
	return Tip, nil
}
func (c *Tip) GetBy(ctx context.Context, id string) (*models.Tip, error) {
	Tip := new(models.Tip)
	tx := c.DB.(*gorm.DB).First(Tip, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return Tip, nil
}

func (c *Tip) DeleteBy(ctx context.Context, id string) (any, error) {
	tx := c.DB.(*gorm.DB).Delete(&models.Tip{}, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx.RowsAffected, nil
}

func (c *Tip) GetAllByOffSet(ctx context.Context, offset, limit int) ([]models.Tip, error) {
	companies := make([]models.Tip, 0)
	tx := c.DB.(*gorm.DB).Limit(limit).Offset(offset).Find(&companies)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return companies, nil
}

func (c *Tip) Search(ctx context.Context, offset, limit int, search string) ([]models.Tip, error) {
	companies := make([]models.Tip, 0)
	//tx := c.DB.(*gorm.DB).Limit(limit).Offset(offset).Find(&companies)
	tx := c.DB.(*gorm.DB).Where("tags @@ to_tsquery(?)", search).Limit(limit).Offset(offset).Find(&companies)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return companies, nil
}
