package gtool

import (
	"context"
	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
)

func WithOrmContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if db == nil {
		return nil
	}

	return apmgorm.WithContext(ctx, db)
}
