package migration1603706814

import "github.com/jinzhu/gorm"

const up = `
ALTER TABLE jobs
ADD COLUMN name text UNIQUE NOT NULL;
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}
