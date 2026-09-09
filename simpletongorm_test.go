package simpletongorm_test

import (
	"testing"
	"time"

	"github.com/fenix-ds/simpletongorm"
	sgenums "github.com/fenix-ds/simpletongorm/enums"
	sgmodels "github.com/fenix-ds/simpletongorm/models"
	"gorm.io/gorm"
)

func TestNewSimpletonGorm_Sucess(t *testing.T) {
	type Test struct{ gorm.Model }

	if _, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	}
}

func TestSimpletonGorm_Save_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
		Name string
	}

	data := Test{}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&sgmodels.SimpletonGormSave{
		TableName: "tests", Data: &data,
	}); err != nil {
		t.Error(err)
	} else if data.ID == 0 {
		t.Error("data not save")
	} else {
		data.Name = "teste"

		if err = sg.Save(&sgmodels.SimpletonGormSave{
			TableName: "tests", Data: &data,
		}); err != nil {
			t.Error(err)
		} else if len(data.Name) == 0 {
			t.Error("data not save")
		}
	}
}

func TestSimpletonGorm_Find_Single_Sucess(t *testing.T) {
	type Test struct{ gorm.Model }

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&sgmodels.SimpletonGormSave{
		TableName: "tests", Data: &Test{},
	}); err != nil {
		t.Error(err)
	} else if result, err := sg.Find(&sgmodels.SimpletonGormFind{
		TableName: "tests",
	}); err != nil {
		t.Error(err)
	} else if *result.Count == 0 {
		t.Error("data not found")
	} else {
		t.Log(result)
	}
}

func TestSimpletonGorm_Find_WithFilters_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
		Date *time.Time
	}
	data := Test{}

	b := true

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
		SeeLog:        &b,
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&sgmodels.SimpletonGormSave{
		TableName: "tests", Data: &data,
	}); err != nil {
		t.Error(err)
	} else if result, err := sg.Find(&sgmodels.SimpletonGormFind{
		TableName: "tests",
		Filters: []sgmodels.SimpletonGormFindFilters{
			{Field: "date", OpComparison: sgenums.OPCN_ISNOTNULL, OpLogic: sgenums.OPLC_AND},
			{Field: "id", Data: data.ID, OpComparison: sgenums.OPCN_EQUAL, OpLogic: sgenums.OPLC_EMPYT},
		},
	}); err != nil {
		t.Error(err)
	} else if *result.Count == 0 {
		t.Error("data not found")
	} else {
		t.Log(result)
	}
}

func TestSimpletonGorm_Find_WithJoins_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
	}
	type TestItem struct {
		gorm.Model
		TestId uint
	}

	test := Test{}
	testItem := TestItem{}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}, TestItem{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&sgmodels.SimpletonGormSave{
		TableName: "tests", Data: &test,
	}); err != nil {
		t.Error(err)
	} else {
		testItem.TestId = test.ID

		if err = sg.Save(&sgmodels.SimpletonGormSave{
			TableName: "test_items", Data: &testItem,
		}); err != nil {
			t.Error(err)
		}

		if result, err := sg.Find(&sgmodels.SimpletonGormFind{
			TableName: "test_items",
			FieldsView: []map[string]any{
				{"id": "new_id"},
			},
			Joins: []sgmodels.SimpletonGormFindJoins{
				{
					Type: sgenums.JT_LEFT, TableMainName: "test_items", TableMainField: "test_id", TableRelatedName: "tests", TableRelatedField: "id",
					TableRelatedFieldsView: []sgmodels.SimpletonGormFindJoinsFieldsView{
						{FieldName: "*"}},
				},
			},
		}); err != nil {
			t.Error(err)
		} else if *result.Count == 0 {
			t.Error("data not found")
		} else {
			t.Log(result)
		}
	}
}

func TestSimpletonGorm_Find_WithFiltersAndJoins_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
	}
	type TestItem struct {
		gorm.Model
		TestId uint
		Name   string
	}

	test := Test{}
	testItem := TestItem{
		Name: "Test 2",
	}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}, TestItem{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&sgmodels.SimpletonGormSave{
		TableName: "tests", Data: &test,
	}); err != nil {
		t.Error(err)
	} else {
		testItem.TestId = test.ID

		if err = sg.Save(&sgmodels.SimpletonGormSave{
			TableName: "test_items", Data: &testItem,
		}); err != nil {
			t.Error(err)
		}

		tbl := "test_items"

		if result, err := sg.Find(&sgmodels.SimpletonGormFind{
			TableName: tbl,
			FieldsView: []map[string]any{
				{"id": "new_id"},
			},
			Joins: []sgmodels.SimpletonGormFindJoins{
				{
					Type: sgenums.JT_LEFT, TableMainName: "test_items", TableMainField: "test_id", TableRelatedName: "tests", TableRelatedField: "id",
					TableRelatedFieldsView: []sgmodels.SimpletonGormFindJoinsFieldsView{
						{FieldName: "*"}},
				},
			},
			Filters: []sgmodels.SimpletonGormFindFilters{
				{
					TableNameFind: &tbl,
					Field:         "name", Data: "Test 2", OpComparison: sgenums.OPCN_EQUAL, OpLogic: sgenums.OPLC_EMPYT,
				},
			},
		}); err != nil {
			t.Error(err)
		} else if *result.Count == 0 {
			t.Error("data not found")
		} else {
			t.Log(result)
		}
	}
}

func TestSimpletonGorm_Find_WithOrders_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
	}

	test := []Test{{}, {}, {}, {}}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else {
		for i := 0; i < len(test); i++ {
			if err = sg.Save(&sgmodels.SimpletonGormSave{
				TableName: "tests", Data: &test[i],
			}); err != nil {
				t.Error(err)
				break
			}
		}

		if result, err := sg.Find(&sgmodels.SimpletonGormFind{
			TableName: "tests",
			Options: &sgmodels.SimpletonGormFindOptions{
				Limit: 10, Offset: 3,
				Orders: []sgmodels.SimpletonGormFindOptionsOrders{
					{Table: "tests", Field: "id", OrderDirection: sgenums.RFOOT_DESC},
				},
			},
		}); err != nil {
			t.Error(err)
		} else if *result.Count == 0 {
			t.Error("data not found")
		} else {
			t.Log(result)
		}
	}
}

func TestSimpletonGorm_Delete_Sucess(t *testing.T) {
	type Test struct{ gorm.Model }

	data := Test{}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&sgmodels.SimpletonGormSave{
		TableName: "tests", Data: &data,
	}); err != nil {
		t.Error(err)
	} else if data.ID == 0 {
		t.Error("data not save")
	} else if err = sg.Delete(&sgmodels.SimpletonGormDelete{
		Type:      sgenums.DT_SOFT,
		TableName: "tests", FieldName: "id", FieldValue: data.ID, Model: &data,
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestSimpletonGorm_SQLFind_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
	}

	test := []Test{{}, {}, {}, {}}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else {
		for i := 0; i < len(test); i++ {
			if err = sg.Save(&sgmodels.SimpletonGormSave{
				TableName: "tests", Data: &test[i],
			}); err != nil {
				t.Error(err)
				break
			}
		}

		if result, err := sg.SQLFind(&sgmodels.SimpletonGormSQL{
			SQL:          "Select * from tests",
			FieldsValues: nil,
		}); err != nil {
			t.Error(err)
		} else if *result.Count == 0 {
			t.Error("data not found")
		} else {
			t.Log(result)
		}
	}
}

func TestSimpletonGorm_SQLExec_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
	}

	test := []Test{{}, {}, {}, {}}

	if sg, err := simpletongorm.NewSimpletonGorm(&sgmodels.SimpletonGormParam{
		Database:      sgenums.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else {
		for i := 0; i < len(test); i++ {
			if err = sg.Save(&sgmodels.SimpletonGormSave{
				TableName: "tests", Data: &test[i],
			}); err != nil {
				t.Error(err)
				break
			}
		}

		var countBeforeCommand uint64
		var countAfterCommand uint64

		if result, err := sg.SQLFind(&sgmodels.SimpletonGormSQL{
			SQL:          "Select * from tests",
			FieldsValues: nil,
		}); err != nil {
			t.Error(err)
		} else if *result.Count == 0 {
			t.Error("data not found")
		} else {
			countBeforeCommand = *result.Count
		}

		if err := sg.SQLExec(&sgmodels.SimpletonGormSQL{
			SQL:          "Delete from tests Where id = ?",
			FieldsValues: []interface{}{1},
		}); err != nil {
			t.Error(err)
		}

		if result, err := sg.SQLFind(&sgmodels.SimpletonGormSQL{
			SQL:          "Select * from tests",
			FieldsValues: nil,
		}); err != nil {
			t.Error(err)
		} else if *result.Count == 0 {
			t.Error("data not found")
		} else {
			countAfterCommand = *result.Count
		}

		if countAfterCommand >= countBeforeCommand {
			t.Error("The data did not change after the command was executed.")
		}
	}
}
