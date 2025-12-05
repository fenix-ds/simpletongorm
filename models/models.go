package models

import (
	"fmt"

	"github.com/fenix-ds/simpletongorm/enuns"
	"github.com/fenix-ds/simpletongorm/utils"
)

type SimpletonGormParam struct {
	Database      enuns.Database
	FilePathOrDns string
	MigrateTables []interface{}
	SeeLog        *bool
}

func (param SimpletonGormParam) CheckData() error {
	if err := param.Database.Validate(); err != nil {
		return err
	}

	if param.Database != enuns.DB_SQLITEINMEMORY && len(param.FilePathOrDns) == 0 {
		return fmt.Errorf("file path or dns not found")
	}

	if param.Database == enuns.DB_SQLITEINMEMORY && param.MigrateTables == nil {
		return fmt.Errorf("to use SQLite in memory, there must be tables to perform the migration")
	}

	if param.MigrateTables != nil {
		if len(param.MigrateTables) == 0 {
			return fmt.Errorf("no items found in the list")
		} else {
			for _, table := range param.MigrateTables {
				if !utils.IsStruct(table) {
					return fmt.Errorf("there are items that are not struct")
				}
			}
		}
	}

	return nil
}

type SimpletonGormSave struct {
	TableName string
	Data      any
	CheckData bool
}

type SimpletonGormFind struct {
	TableName  string
	FieldsView []map[string]any
	Joins      []SimpletonGormFindJoins
	Filters    []SimpletonGormFindFilters
	Options    *SimpletonGormFindOptions
	CheckData  bool
}

type SimpletonGormFindJoins struct {
	Type                            enuns.JoinType
	TableMainName                   string
	TableMainField                  string
	TableRelatedName                string
	TableRelatedNameRename          any
	TableRelatedConditionAdditional any
	TableRelatedField               string
	TableRelatedFieldsView          []SimpletonGormFindJoinsFieldsView
}

type SimpletonGormFindJoinsFieldsView struct {
	FieldName   string
	FieldRename *string
}

type SimpletonGormFindFilters struct {
	TableNameFind any
	Field         string
	Data          any
	OpComparison  enuns.OpComparison
	OpLogic       enuns.OpLogic
}

type SimpletonGormFindOptions struct {
	Limit  uint
	Offset uint
	Orders []SimpletonGormFindOptionsOrders
}

type SimpletonGormFindOptionsOrders struct {
	OrderDirection enuns.OptionsOrderDirection
	Table          any
	Field          string
}

type SimpletonGormResult struct {
	Data  []map[string]any
	Count *uint64
}

type SimpletonGormDelete struct {
	Type       enuns.DeleteType
	TableName  string
	FieldName  string
	FieldValue any
	Model      any
	CheckData  bool
}
