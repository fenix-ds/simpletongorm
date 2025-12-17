package simpletongorm

import (
	"fmt"

	"github.com/fenix-ds/simpletongorm/enuns"
	"github.com/fenix-ds/simpletongorm/models"
	"github.com/fenix-ds/simpletongorm/utils"
	"gorm.io/gorm"
)

func dbMigrate(db *gorm.DB, tables []interface{}) error {
	if err := db.AutoMigrate(tables...); err != nil {
		return err
	}

	return nil
}

func (sg *SimpletonGorm) dbConnectionActive() (*gorm.DB, error) {
	if sg.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return sg.db, nil
}

func (sg *SimpletonGorm) applyJoinsFields(fieldsQueryMain string, joins []models.SimpletonGormFindJoins) string {
	//CAMPOS QUE VEM DO JOIN
	var fieldsListView string
	for index, join := range joins {
		if join.Type == enuns.JT_LEFT {
			for index, fieldView := range join.TableRelatedFieldsView {
				if fieldView.FieldName == "*" {
					fieldsListView = fmt.Sprintf("%s.%s", join.TableRelatedName, fieldView.FieldName)
					break
				} else {
					if join.TableRelatedNameRename != nil {
						fieldsListView = fieldsListView + fmt.Sprintf("%s.%s", join.TableRelatedNameRename, fieldView.FieldName)
					} else {
						fieldsListView = fieldsListView + fmt.Sprintf("%s.%s", join.TableRelatedName, fieldView.FieldName)
					}

					if fieldView.FieldRename != nil {
						fieldsListView = fieldsListView + " AS " + *fieldView.FieldRename
					}

					if index+1 < len(join.TableRelatedFieldsView) {
						fieldsListView = fieldsListView + ", "
					}
				}
			}

			if index+1 < len(joins) {
				if len(fieldsListView) > 0 && len(joins[index+1].TableRelatedFieldsView) > 0 {
					fieldsListView = fieldsListView + ", "
				}
			}
		}
	}

	if len(fieldsListView) == 0 {
		return fieldsQueryMain
	} else {
		return fmt.Sprintf("%s, %s", fieldsQueryMain, fieldsListView)
	}
}

func (sg *SimpletonGorm) applyJoinsConditions(query *gorm.DB, joins []models.SimpletonGormFindJoins, isCount bool) error {
	for _, join := range joins {
		var joinSQL string

		switch join.Type {
		case enuns.JT_LEFT:
			joinSQL = fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s", join.TableRelatedName, join.TableRelatedName,
				join.TableRelatedField, join.TableMainName, join.TableMainField)

			if join.TableRelatedNameRename != nil {
				joinSQL = fmt.Sprintf("LEFT JOIN %s as %s ON %s.%s = %s.%s", join.TableRelatedName, join.TableRelatedNameRename,
					join.TableRelatedNameRename, join.TableRelatedField, join.TableMainName, join.TableMainField)
			}

			if join.TableRelatedConditionAdditional != nil {
				joinSQL = fmt.Sprintf("%s %s", joinSQL, join.TableRelatedConditionAdditional)
			}
		case enuns.JT_CROSS:
			joinSQL = fmt.Sprintf("CROSS JOIN %s", join.TableRelatedName)

			if join.TableRelatedNameRename != nil {
				joinSQL = fmt.Sprintf("CROSS JOIN %s as %s", join.TableRelatedName, join.TableRelatedNameRename)
			}
		default:
			return fmt.Errorf("join invalid")
		}

		query.Joins(joinSQL)
	}

	return nil
}

func (sg *SimpletonGorm) validateJoins(db *gorm.DB, joins []models.SimpletonGormFindJoins) error {
	for _, join := range joins {
		switch join.Type {
		case enuns.JT_LEFT:
			//CHECK TABLE MAIN
			if !db.Migrator().HasTable(join.TableMainName) {
				return fmt.Errorf("%s table to create join not found", join.TableMainName)
			}

			//CHECK FIELD OF THE TABLE MAIN
			if !db.Migrator().HasColumn(join.TableMainName, join.TableMainField) {
				return fmt.Errorf("%s field not found in %s", join.TableMainField, join.TableMainName)
			}

			//CHECK TABLE RELATED
			if !db.Migrator().HasTable(join.TableRelatedName) {
				return fmt.Errorf("%s table to create join not found", join.TableRelatedName)
			}

			//CHECK FIELD OF THE TABLE RELATED
			if !db.Migrator().HasColumn(join.TableRelatedName, join.TableRelatedField) {
				return fmt.Errorf("%s field not found in %s", join.TableMainField, join.TableRelatedName)
			}

			//CHECK THE NEW NAME OF THE RELATED TABLE
			if join.TableRelatedNameRename != nil {
				if !utils.IsSnakeCase(join.TableRelatedNameRename.(string)) {
					return fmt.Errorf("%s table name does not follow the required pattern Snake Case (variable_name)", join.TableRelatedNameRename.(string))
				}
			}

			//CHECK FIELDS VIEW OF THE TABEL RELATED
			if len(join.TableRelatedFieldsView) != 0 {
				for _, fieldView := range join.TableRelatedFieldsView {
					if fieldView.FieldName != "*" {
						if !db.Migrator().HasColumn(join.TableRelatedName, fieldView.FieldName) {
							return fmt.Errorf("%s field not found in %s", join.TableMainField, join.TableMainName)
						}

						if fieldView.FieldRename != nil {
							if !utils.IsSnakeCase(join.TableRelatedNameRename.(string)) {
								return fmt.Errorf("%s field name does not follow the required pattern Snake Case (variable_name)", *fieldView.FieldRename)
							}
						}
					}
				}
			}
		case enuns.JT_CROSS:
			if !db.Migrator().HasTable(join.TableRelatedName) {
				return fmt.Errorf("%s table to create join not found", join.TableRelatedName)
			}

			if join.TableRelatedNameRename != nil {
				if !utils.IsSnakeCase(join.TableRelatedNameRename.(string)) {
					return fmt.Errorf("%s table name does not follow the required pattern Snake Case (variable_name)", join.TableRelatedNameRename.(string))
				}
			}
		default:
			return fmt.Errorf("join invalid")
		}
	}

	return nil
}

func (sg *SimpletonGorm) applyFilters(filters []models.SimpletonGormFindFilters, query *gorm.DB) error {
	if query.Error != nil {
		return query.Error
	}

	for index, filter := range filters {
		var err error
		dataQuery, err := sg.generateSQLFilter(&filter)
		if err != nil {
			return err
		}

		if filter.Data != nil {
			if filter.Data, err = sg.modifyDataAccordingToRepositoryOperatorComparison(filter.OpComparison, filter.Data); err != nil {
				return err
			}
		}

		if index == 0 {
			query = query.Where(*dataQuery, filter.Data)
		} else {
			switch filters[index-1].OpLogic {
			case enuns.OPLC_OR:
				query = query.Or(*dataQuery, filter.Data)
			case enuns.OPLC_AND:
				query = query.Where(*dataQuery, filter.Data)
			}
		}
	}

	return nil
}

func (sg *SimpletonGorm) generateSQLFilter(filter *models.SimpletonGormFindFilters) (*string, error) {
	var query string
	var err error

	if err = filter.OpComparison.Validate(); err != nil {
		return nil, err
	}

	switch filter.OpComparison {
	case enuns.OPCN_LIKE_ANYWHERE, enuns.OPCN_LIKE_JUSTBEGINNING:
		query = fmt.Sprintf("%s.%s LIKE ?", filter.TableNameFind, filter.Field)
	case enuns.OPCN_BETWEEN:
		query, err = sg.setDataToQueryComparison_BETWEEN(filter)
	case enuns.OPCN_IN:
		query, err = sg.setDataToQueryComparison_IN(filter)
	case enuns.OPCN_IS:
		query = fmt.Sprintf("%s.%s IS ?", filter.TableNameFind, filter.Field)
	case enuns.OPCN_ISNULL:
		query = fmt.Sprintf("%s.%s ISNULL", filter.TableNameFind, filter.Field)
	case enuns.OPCN_LESS_EQUAL_ISNULL:
		query = fmt.Sprintf("%s.%s <= ? OR %s.%s IS NULL", filter.TableNameFind, filter.Field, filter.TableNameFind, filter.Field)
	default:
		query = fmt.Sprintf("%s.%s %s ?", filter.TableNameFind, filter.Field, string(filter.OpComparison))
	}

	if err != nil {
		return nil, err
	}

	return &query, nil
}

func (sg *SimpletonGorm) setDataToQueryComparison_BETWEEN(filter *models.SimpletonGormFindFilters) (string, error) {
	if dataList, ok := filter.Data.([]interface{}); ok {
		if len(dataList) != 2 {
			return "", fmt.Errorf("the BETWEEN selector must have 2 values")
		}

		result := fmt.Sprintf("%s.%s >= ? AND %s.%s <= ?", filter.TableNameFind, filter.Field, filter.TableNameFind, filter.Field)

		return result, nil

	} else {
		return "", fmt.Errorf("values sent are not compatible with the selected comparator")
	}
}

func (sg *SimpletonGorm) setDataToQueryComparison_IN(filter *models.SimpletonGormFindFilters) (string, error) {
	if dataList, ok := filter.Data.([]interface{}); ok {
		if len(dataList) < 2 {
			return "", fmt.Errorf("the IN selector must have at least 2 values")
		}

		var values string
		for i := 1; i < len(dataList); i++ {
			if i == 1 {
				values = "?"
			} else {
				values = values + (",?")
			}

		}

		result := fmt.Sprintf("%s.%s IN (%s)", filter.TableNameFind, filter.Field, values)

		return result, nil

	} else {
		return "", fmt.Errorf("values sent are not compatible with the selected comparator")
	}
}

func (sg *SimpletonGorm) modifyDataAccordingToRepositoryOperatorComparison(comparison enuns.OpComparison, data interface{}) (interface{}, error) {
	if err := comparison.Validate(); err != nil {
		return nil, err
	}

	var result interface{}
	if str, ok := data.(string); !ok {
		if len(str) > 0 {
			return nil, fmt.Errorf("value is not string")
		}

		result = data
	} else {
		switch comparison {
		case enuns.OPCN_LIKE_JUSTBEGINNING:
			result = str + "%"
		case enuns.OPCN_LIKE_ANYWHERE:
			result = "%" + str + "%"
		default:
			result = str
		}
	}

	return result, nil
}
