package models

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
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
			fmt.Printf("Transaction (%s) is rolled back \n", t.name)
			t.tx.Rollback()
		} else {
			fmt.Printf("Transaction (%s) is commited \n", t.name)
			t.tx.Commit()
		}
	})
}

func (t *Tx) Fail(err error) {
	fmt.Printf("Transaction (%s) Error: %s \n", t.name, err)
	t.failed = true
}

func NewTx(name string, db *gorm.DB) (t *Tx, c *gorm.DB) {
	c = db.Begin()
	t = &Tx{
		name:   name,
		failed: false,
		tx:     c,
	}
	fmt.Printf("Transaction begins: %s \n", name)
	return
}
