package myqueue

import (
	"github.com/xblockchainlabs/myqueue/models"
	"github.com/xblockchainlabs/myqueue/utils"
)

// Migrate the schema of database if needed
func connectMySQL(user string, pass string, host string, port string, name string) error {
	var job models.Job
	var sched models.Schedule
	utils.SetConnection(
		user,
		pass,
		host,
		port,
		name,
	)
	db, err := utils.GetDB()
	if err != nil {
		return err
	}
	db.AutoMigrate(&job)
	db.AutoMigrate(&sched).AddForeignKey("job_id", "jobs(id)", "CASCADE", "SET NULL")
	return nil
}


func GetQueue(user string, pass string, host string, port string, name string) (*ConsumerGroup, error) {
	err := connectMySQL(user, pass, host, port, name)
	if err != nil {
		return nil, err
	}
	cg := NewCG()
	return cg, nil
}
