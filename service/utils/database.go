package utils

import (
	"fmt"
	"service/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DatabaseManager      = gorm.DB{}
	DatabaseManagerSlave = gorm.DB{}
	masterDsn            = GetEnv("DB_INFO_MASTER")
	slaveDsn             = GetEnv("DB_INFO_SLAVE")
)

func InitDB() bool {
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	success := false

	for range ticker.C {
		success = excute()
		PrintObj("InitDB retry:"+ToString(count), "")

		if count == 60 || success {
			return success
		}

		count++
	}

	return true
}

func excute() bool {
	// connect master
	masterDb, err := gorm.Open(postgres.Open(masterDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		PrintObj(err.Error(), "init db err")
		return false
	}

	DatabaseManager = *masterDb
	PrintObj("db master connected")

	// connect slave
	slaveDb, err := gorm.Open(postgres.Open(slaveDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		PrintObj(err.Error(), "init db err")
		return false
	}

	DatabaseManagerSlave = *slaveDb
	PrintObj("db slave connected")

	// init uuid function
	init := masterDb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if init.Error != nil {
		PrintObj(init.Error.Error(), "gorm DB init uuid function fail")
		return false
	}

	// init table or migrate
	migrateErr := masterDb.AutoMigrate(
		&model.User{},
		&model.ArContent{},
		&model.Program{},
		&model.PermissionGroup{},
		&model.Log{},
		&model.StatisticArContentUserDay{},
		&model.StatisticArContentUserSum{},
		&model.MailState{},
		&model.AccountSettings{},
	)

	if migrateErr != nil {
		PrintObj(migrateErr, "gorm migrateErr")
		return false
	}

	TestDbMasterSlave()

	return true
}

func DisConnect() bool {
	db, _ := DatabaseManager.DB()
	err := db.Close()
	fmt.Println(err)

	return err == nil
}

func TestDbMasterSlave() {
	testId := uuid.New()
	checkValue := uuid.New()

	res := DatabaseManager.Create(&model.Log{
		Id:    testId,
		Key:   "test db master slave",
		Value: checkValue.String(),
	})

	if res.Error != nil {
		panic(res.Error.Error())
	}

	time.Sleep(1 * time.Second)
	// PrintObj("sleep over")

	findResult := model.Log{}
	DatabaseManagerSlave.Where(model.Log{
		Id: testId,
	}).Find(&findResult)

	if res.Error != nil {
		panic(res.Error.Error())
	}

	if findResult.Value != checkValue.String() {
		PrintObj("", "test db master slave")
		panic("test fail")
	} else {
		PrintObj("test success", "test db master slave")
	}
}
