package simpletongorm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fenix-ds/simpletongorm/enuns"
	"github.com/fenix-ds/simpletongorm/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SimpletonGorm struct {
	database *enuns.Database
	db       *gorm.DB
}

func NewSimpletonGorm(param *models.SimpletonGormParam) (sg *SimpletonGorm, err error) {
	if err := param.CheckData(); err != nil {
		return nil, err
	}

	gormConfig := gorm.Config{}
	if param.SeeLog != nil && *param.SeeLog {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var db *gorm.DB
	switch param.Database {
	case enuns.DB_SQLITEINMEMORY:
		db, err = gorm.Open(sqlite.Open(":memory:"), &gormConfig)
	case enuns.DB_SQLITEFILE:
		db, err = gorm.Open(sqlite.Open(param.FilePathOrDns), &gormConfig)
	case enuns.DB_MARIADB:
		db, err = gorm.Open(mysql.Open(param.FilePathOrDns), &gormConfig)
	case enuns.DB_POSTGRESQL:
		db, err = gorm.Open(postgres.Open(param.FilePathOrDns), &gormConfig)
	default:
		return nil, fmt.Errorf("database type invalid")
	}

	if err != nil {
		return nil, err
	}

	if param.Database != enuns.DB_SQLITEINMEMORY {
		if param.MigrateTables != nil {
			if err = dbMigrate(db, param.MigrateTables); err != nil {
				return nil, err
			}
		}
	} else {
		if param.MigrateTables == nil {
			return nil, fmt.Errorf("there are no models to migrate to the database")
		} else if err = dbMigrate(db, param.MigrateTables); err != nil {
			return nil, err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &SimpletonGorm{
		database: &param.Database,
		db:       db,
	}, nil
}

func (sg *SimpletonGorm) Save(param *models.SimpletonGormSave) error {
	db, err := sg.dbConnectionActive()
	if err != nil {
		return err
	}

	if param.CheckData {
		if !db.Migrator().HasTable(param.TableName) {
			return fmt.Errorf("%s table not found", param.TableName)
		}

		if param.Data == nil {
			return fmt.Errorf("no data found to save")
		}
	}

	val := reflect.ValueOf(param.Data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.Uint() > 0 {
		original := reflect.New(val.Type()).Interface()
		if err := db.Table(param.TableName).First(original, idField.Uint()).Error; err == nil {
			createdField := val.FieldByName("CreatedAt")
			origCreated := reflect.ValueOf(original).Elem().FieldByName("CreatedAt")
			if createdField.IsValid() && origCreated.IsValid() {
				createdField.Set(origCreated)
			}
		}
	}

	if result := db.Table(param.TableName).Save(param.Data); result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("no records were saved")
	}

	if *sg.database != enuns.DB_SQLITEINMEMORY {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return nil
}

func (sg *SimpletonGorm) Find(param *models.SimpletonGormFind) (result *models.SimpletonGormResult, err error) {
	// DATABASE CONNECTION
	db, err := sg.dbConnectionActive()
	if err != nil {
		return nil, err
	}

	if param.CheckData {
		if !db.Migrator().HasTable(param.TableName) {
			return nil, fmt.Errorf("%s table not found", param.TableName)
		}
	}

	//INIT QUERYS
	var fieldsQueryMain string
	queryMain := db.Table(param.TableName)
	queryCount := db.Table(param.TableName)
	queryCount = queryCount.Select("*")

	//FIELDS VIEW
	if len(param.FieldsView) == 0 {
		fieldsQueryMain = fmt.Sprintf("%s.*", param.TableName)
	} else {
		for index, fieldname := range param.FieldsView {
			for fieldNameOrigin, fieldNameAs := range fieldname {
				if param.CheckData {
					if !db.Migrator().HasColumn(param.TableName, fieldNameOrigin) {
						return nil, fmt.Errorf("%s field not found in %s", fieldNameOrigin, param.TableName)
					}
				}

				fieldResult := fmt.Sprintf("%s.%s", param.TableName, fieldNameOrigin)
				if fieldNameAs != nil {
					fieldResult = fmt.Sprintf("%s AS %s", fieldResult, fieldNameAs)
				}

				if index == 0 {
					fieldsQueryMain += fieldResult
				} else {
					fieldsQueryMain += "," + fieldResult
				}
			}
		}
	}

	//JOINS
	if len(param.Joins) > 0 {
		for _, j := range param.Joins {
			if j.Type != enuns.JT_LEFT && j.Type != enuns.JT_CROSS {
				return nil, fmt.Errorf("join type invalid")
			}
		}

		if param.CheckData {
			if err := sg.validateJoins(db, param.Joins); err != nil {
				return nil, err
			}
		}

		queryMain.Select(sg.applyJoinsFields(fieldsQueryMain, param.Joins))

		if err := sg.applyJoinsConditions(queryMain, param.Joins, false); err != nil {
			return nil, err
		}

		if err := sg.applyJoinsConditions(queryCount, param.Joins, true); err != nil {
			return nil, err
		}
	} else {
		queryMain = queryMain.Select(fieldsQueryMain)
	}

	//FILTERS
	if len(param.Filters) > 0 {
		for index := range param.Filters {
			if param.Filters[index].TableNameFind == nil {
				param.Filters[index].TableNameFind = param.TableName
			}
		}

		if err := sg.applyFilters(param.Filters, queryMain); err != nil {
			return nil, err
		}

		if err := sg.applyFilters(param.Filters, queryCount); err != nil {
			return nil, err
		}
	}

	queryMain = queryMain.Where(fmt.Sprintf("%s.deleted_at IS NULL", param.TableName))
	queryCount = queryCount.Where(fmt.Sprintf("%s.deleted_at IS NULL", param.TableName))

	//OPTIONS
	if options := param.Options; options != nil {
		if options.Limit > 0 {
			queryMain = queryMain.Limit(int(options.Limit))
		}

		if options.Offset > 0 {
			queryMain = queryMain.Offset(int(options.Offset))
		}

		if len(options.Orders) > 0 {
			var orderData string
			for index, order := range options.Orders {
				if order.Table != nil {
					if param.CheckData {
						if !db.Migrator().HasColumn(order.Table, order.Field) {
							return nil, fmt.Errorf("%s field not found in %s", order.Field, order.Table)
						}
					}

					orderData += fmt.Sprintf("%s.%s", order.Table, order.Field)

				} else {
					orderData += order.Field
				}

				if order.OrderDirection == enuns.RFOOT_ASC {
					orderData += " ASC"
				} else {
					orderData += " DESC"
				}

				if index+1 < len(options.Orders) {
					orderData += ", "
				}
			}

			queryMain = queryMain.Order(orderData)
		}
	}

	//BUSCAR RESULTADOS
	if err = queryMain.Error; err != nil {
		return nil, queryMain.Error
	}

	if err = queryCount.Error; err != nil {
		return nil, queryCount.Error
	}

	var count int64
	if err := queryCount.Count(&count).Error; err != nil {
		return nil, err
	}

	countResult := uint64(count)

	var data []map[string]any
	if err := queryMain.Find(&data).Error; err != nil {
		return nil, err
	}

	if *sg.database != enuns.DB_SQLITEINMEMORY {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return &models.SimpletonGormResult{Data: data, Count: &countResult}, nil
}

func (sg *SimpletonGorm) Delete(param *models.SimpletonGormDelete) error {
	db, err := sg.dbConnectionActive()
	if err != nil {
		return err
	}

	if param.Type != enuns.DT_SOFT && param.Type != enuns.DT_PERMANENT {
		return fmt.Errorf("invalid deletion type")
	}

	if param.CheckData {
		if !db.Migrator().HasTable(param.TableName) {
			return fmt.Errorf("%s table not found", param.TableName)
		}

		if !db.Migrator().HasColumn(param.TableName, param.FieldName) {
			return fmt.Errorf("%s field not found in %s", param.FieldName, param.TableName)
		}
	}

	if param.FieldValue == nil {
		return fmt.Errorf("no data found to delete")
	} else if param.Type == enuns.DT_SOFT && param.Model == nil {
		return fmt.Errorf("no data found to delete")
	}

	var result *gorm.DB
	switch param.Type {
	case enuns.DT_SOFT:
		result = db.Table(param.TableName).Where(param.FieldName+"=?", param.FieldValue).Delete(param.Model)
	case enuns.DT_PERMANENT:
		result = db.Table(param.TableName).Where(param.FieldName+"=?", param.FieldValue).Delete(nil)
	}

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no records were deleted")
	}

	if *sg.database != enuns.DB_SQLITEINMEMORY {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return nil
}
