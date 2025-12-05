package simpletongorm

import (
	"errors"
	"strings"
)

type RepositoryOperatorComparison string

const (
	RYOC_EQUAL              RepositoryOperatorComparison = "="
	RYOC_EQUAL_NOT          RepositoryOperatorComparison = "<>"
	RYOC_GREATER            RepositoryOperatorComparison = ">"
	RYOC_GREATER_EQUAL      RepositoryOperatorComparison = ">="
	RYOC_LESS               RepositoryOperatorComparison = "<"
	RYOC_LESS_EQUAL         RepositoryOperatorComparison = "<="
	RYOC_LIKE               RepositoryOperatorComparison = "LIKE"
	RYOC_LIKE_ANYWHERE      RepositoryOperatorComparison = "%LIKE%"
	RYOC_LIKE_JUSTBEGINNING RepositoryOperatorComparison = "LIKE%"
	RYOC_BETWEEN            RepositoryOperatorComparison = "BETWEEN"
	RYOC_IN                 RepositoryOperatorComparison = "IN"
	RYOC_IS                 RepositoryOperatorComparison = "IS"
	RYOC_LESS_EQUAL_ISNULL  RepositoryOperatorComparison = "<=,ISNULL"
)

func (ryoc *RepositoryOperatorComparison) Validate() error {
	switch *ryoc {
	case
		RYOC_EQUAL,
		RYOC_EQUAL_NOT,
		RYOC_GREATER,
		RYOC_GREATER_EQUAL,
		RYOC_LESS,
		RYOC_LESS_EQUAL,
		RYOC_LIKE,
		RYOC_LIKE_ANYWHERE,
		RYOC_LIKE_JUSTBEGINNING,
		RYOC_BETWEEN,
		RYOC_IN,
		RYOC_IS,
		RYOC_LESS_EQUAL_ISNULL:
		return nil
	default:
		return errors.New("invalid comparison operator")
	}
}

type RepositoryOperatorLogic int

const (
	RYOL_EMPYT RepositoryOperatorLogic = iota
	RYOL_AND
	RYOL_OR
)

type RepositoryFindOptionsOrderType int

const (
	RFOOT_ASC RepositoryFindOptionsOrderType = iota
	RFOOT_DESC
)

func StrToRepositoryFindOptionsOrderType(value string) (*RepositoryFindOptionsOrderType, error) {
	var result RepositoryFindOptionsOrderType
	switch strings.ToLower(value) {
	case "asc":
		result = RFOOT_ASC
	case "desc":
		result = RFOOT_DESC
	default:
		return nil, errors.New("value sent is not valid RepositoryFindOptionsOrderType")
	}

	return &result, nil
}
