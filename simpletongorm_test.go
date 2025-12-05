package simpletongorm_test

import (
	"testing"

	"github.com/fenix-ds/simpletongorm"
	"github.com/fenix-ds/simpletongorm/enuns"
	"github.com/fenix-ds/simpletongorm/models"
	"gorm.io/gorm"
)

func TestNewSimpletonGorm_Sucess(t *testing.T) {
	type Test struct{ gorm.Model }

	if _, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
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

	if sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&models.SimpletonGormSave{
		TableName: "tests", Data: &data,
	}); err != nil {
		t.Error(err)
	} else if data.ID == 0 {
		t.Error("data not save")
	} else {
		data.Name = "teste"

		if err = sg.Save(&models.SimpletonGormSave{
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

	if sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&models.SimpletonGormSave{
		TableName: "tests", Data: &Test{},
	}); err != nil {
		t.Error(err)
	} else if result, err := sg.Find(&models.SimpletonGormFind{
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
	type Test struct{ gorm.Model }
	data := Test{}

	if sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&models.SimpletonGormSave{
		TableName: "tests", Data: &data,
	}); err != nil {
		t.Error(err)
	} else if result, err := sg.Find(&models.SimpletonGormFind{
		TableName: "tests",
		Filters: []models.SimpletonGormFindFilters{
			{Field: "id", Data: data.ID, OpComparison: enuns.OPCN_EQUAL, OpLogic: enuns.OPLC_EMPYT},
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

	if sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}, TestItem{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&models.SimpletonGormSave{
		TableName: "tests", Data: &test,
	}); err != nil {
		t.Error(err)
	} else {
		testItem.TestId = test.ID

		if err = sg.Save(&models.SimpletonGormSave{
			TableName: "test_items", Data: &testItem,
		}); err != nil {
			t.Error(err)
		}

		if result, err := sg.Find(&models.SimpletonGormFind{
			TableName: "test_items",
			FieldsView: []map[string]any{
				{"id": "new_id"},
			},
			Joins: []models.SimpletonGormFindJoins{
				{
					Type: enuns.JT_LEFT, TableMainName: "test_items", TableMainField: "test_id", TableRelatedName: "tests", TableRelatedField: "id",
					TableRelatedFieldsView: []models.SimpletonGormFindJoinsFieldsView{
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

func TestSimpletonGorm_Find_WithOrders_Sucess(t *testing.T) {
	type Test struct {
		gorm.Model
	}

	test := []Test{{}, {}, {}, {}}

	if sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else {
		for i := 0; i < len(test); i++ {
			if err = sg.Save(&models.SimpletonGormSave{
				TableName: "tests", Data: &test[i],
			}); err != nil {
				t.Error(err)
				break
			}
		}

		if result, err := sg.Find(&models.SimpletonGormFind{
			TableName: "tests",
			Options: &models.SimpletonGormFindOptions{
				Limit: 10, Offset: 3,
				Orders: []models.SimpletonGormFindOptionsOrders{
					{Table: "tests", Field: "id", OrderDirection: enuns.RFOOT_DESC},
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

	if sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
		Database:      enuns.DB_SQLITEINMEMORY,
		MigrateTables: []interface{}{Test{}},
	}); err != nil {
		t.Error(err)
	} else if err = sg.Save(&models.SimpletonGormSave{
		TableName: "tests", Data: &data,
	}); err != nil {
		t.Error(err)
	} else if data.ID == 0 {
		t.Error("data not save")
	} else if err = sg.Delete(&models.SimpletonGormDelete{
		Type:      enuns.DT_SOFT,
		TableName: "tests", FieldName: "id", FieldValue: data.ID, Model: &data,
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}
