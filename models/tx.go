package models

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/xblockchainlabs/myqueue/utils"
)

type Tx struct {
	name   string
	once   sync.Once
	failed bool
	tx     *gorm.DB
}

func (t *Tx) Close() {
	t.once.Do(func() {
		if t.failed {
			utils.InfoLogf("Transaction (%s) is rolled back \n", t.name)
			t.tx.Rollback()
		} else {
			utils.InfoLogf("Transaction (%s) is commited \n", t.name)
			t.tx.Commit()
		}
	})
}

func (t *Tx) Fail(err error) {
	utils.WarningLogf("Transaction (%s) Error: %s \n", t.name, err)
	t.failed = true
}

func NewTx(name string, db *gorm.DB) (t *Tx, c *gorm.DB) {
	c = db.Begin()
	t = &Tx{
		name:   name,
		failed: false,
		tx:     c,
	}
	utils.InfoLogf("Transaction begins: %s \n", name)
	return
}
