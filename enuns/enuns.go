package enuns

import (
	"errors"
	"fmt"
	"strings"
)

type Database uint

const (
	DB_SQLITEINMEMORY Database = iota + 1
	DB_SQLITEFILE
	DB_MARIADB
	DB_POSTGRESQL
)

func (db Database) Validate() error {
	if db == 0 || db > 4 {
		return fmt.Errorf("database invalid")
	}

	return nil
}

func (db Database) ToPoint() *Database {
	return &db
}

type JoinType string

const (
	JT_LEFT  JoinType = "LEFT"
	JT_CROSS JoinType = "CROSS"
)

func (jt JoinType) ToPoint() *JoinType {
	return &jt
}

type OpComparison string

const (
	OPCN_EQUAL              OpComparison = "="
	OPCN_EQUAL_NOT          OpComparison = "<>"
	OPCN_GREATER            OpComparison = ">"
	OPCN_GREATER_EQUAL      OpComparison = ">="
	OPCN_LESS               OpComparison = "<"
	OPCN_LESS_EQUAL         OpComparison = "<="
	OPCN_LIKE               OpComparison = "LIKE"
	OPCN_LIKE_ANYWHERE      OpComparison = "%LIKE%"
	OPCN_LIKE_JUSTBEGINNING OpComparison = "LIKE%"
	OPCN_BETWEEN            OpComparison = "BETWEEN"
	OPCN_IN                 OpComparison = "IN"
	OPCN_IS                 OpComparison = "IS"
	OPCN_ISNULL             OpComparison = "ISNULL"
	OPCN_LESS_EQUAL_ISNULL  OpComparison = "<=,ISNULL"
)

func (ryoc *OpComparison) Validate() error {
	switch *ryoc {
	case
		OPCN_EQUAL,
		OPCN_EQUAL_NOT,
		OPCN_GREATER,
		OPCN_GREATER_EQUAL,
		OPCN_LESS,
		OPCN_LESS_EQUAL,
		OPCN_LIKE,
		OPCN_LIKE_ANYWHERE,
		OPCN_LIKE_JUSTBEGINNING,
		OPCN_BETWEEN,
		OPCN_IN,
		OPCN_IS,
		OPCN_ISNULL,
		OPCN_LESS_EQUAL_ISNULL:
		return nil
	default:
		return errors.New("invalid comparison operator")
	}
}

type OpLogic int

const (
	OPLC_EMPYT OpLogic = iota
	OPLC_AND
	OPLC_OR
)

type OptionsOrderDirection uint

const (
	RFOOT_ASC OptionsOrderDirection = iota
	RFOOT_DESC
)

func StrToOptionsOrderType(value string) (*OptionsOrderDirection, error) {
	var result OptionsOrderDirection
	switch strings.ToLower(value) {
	case "asc":
		result = RFOOT_ASC
	case "desc":
		result = RFOOT_DESC
	default:
		return nil, errors.New("value sent is not valid OptionsOrderType")
	}

	return &result, nil
}

type DeleteType uint

const (
	DT_SOFT DeleteType = iota
	DT_PERMANENT
)
