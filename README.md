# Simpleton Gorm

Simpleton Gorm is a lightweight Go package that simplifies common GORM database operations.
It supports SQLite in-memory, SQLite file, MariaDB, and PostgreSQL, with built-in migration, save, find, delete and raw SQL utilities.

## Requirements

- Go version: `1.25.4`

## Dependencies

- `gorm.io/gorm` v1.31.1
- `gorm.io/driver/sqlite` v1.6.0
- `gorm.io/driver/mysql` v1.6.0
- `gorm.io/driver/postgres` v1.6.0

Indirect dependencies used by the module:

- `filippo.io/edwards25519` v1.1.0
- `github.com/go-sql-driver/mysql` v1.9.3
- `github.com/jackc/pgpassfile` v1.0.0
- `github.com/jackc/pgservicefile` v0.0.0-20240606120523-5a60cdf6a761
- `github.com/jackc/pgx/v5` v5.7.6
- `github.com/jackc/puddle/v2` v2.2.2
- `github.com/mattn/go-sqlite3` v1.14.32
- `github.com/jinzhu/inflection` v1.0.0
- `github.com/jinzhu/now` v1.1.5
- `golang.org/x/crypto` v0.46.0
- `golang.org/x/sync` v0.19.0
- `golang.org/x/text` v0.32.0

## Installation

```bash
go get github.com/fenix-ds/simpletongorm
```

## Import

```go
import "github.com/fenix-ds/simpletongorm"
import "github.com/fenix-ds/simpletongorm/enuns"
import "github.com/fenix-ds/simpletongorm/models"
```

## Quick start

```go
package main

import (
    "fmt"

    "github.com/fenix-ds/simpletongorm"
    "github.com/fenix-ds/simpletongorm/enuns"
    "github.com/fenix-ds/simpletongorm/models"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name  string
    Email string
}

func main() {
    sg, err := simpletongorm.NewSimpletonGorm(&models.SimpletonGormParam{
        Database:      enuns.DB_SQLITEINMEMORY,
        MigrateTables: []interface{}{User{}},
    })
    if err != nil {
        panic(err)
    }

    user := User{Name: "Alice", Email: "alice@example.com"}
    if err := sg.Save(&models.SimpletonGormSave{TableName: "users", Data: &user}); err != nil {
        panic(err)
    }

    result, err := sg.Find(&models.SimpletonGormFind{TableName: "users"})
    if err != nil {
        panic(err)
    }

    fmt.Printf("found %d rows\n", *result.Count)
    fmt.Println(result.Data)
}
```

## Supported databases

- `enuns.DB_SQLITEINMEMORY`
- `enuns.DB_SQLITEFILE`
- `enuns.DB_MARIADB`
- `enuns.DB_POSTGRESQL`

## `NewSimpletonGorm` parameters

`models.SimpletonGormParam` fields:

- `Database` - database driver type from `enuns.Database`
- `FilePathOrDns` - SQLite file path or connection DSN for MariaDB / PostgreSQL
- `MigrateTables` - list of model structs to auto migrate
- `SeeLog` - optional pointer to enable GORM SQL logging

## Save data

```go
err := sg.Save(&models.SimpletonGormSave{
    TableName: "users",
    Data:      &user,
})
```

If `Data` contains a struct with `ID`, `Save` updates existing records and preserves `CreatedAt` when the record already exists.

## Find records

```go
result, err := sg.Find(&models.SimpletonGormFind{
    TableName: "users",
    Filters: []models.SimpletonGormFindFilters{
        {
            Field:        "name",
            Data:         "Alice",
            OpComparison: enuns.OPCN_EQUAL,
            OpLogic:      enuns.OPLC_EMPYT,
        },
    },
    Options: &models.SimpletonGormFindOptions{
        Limit:  10,
        Offset: 0,
        Orders: []models.SimpletonGormFindOptionsOrders{{
            Table:          "users",
            Field:          "id",
            OrderDirection: enuns.RFOOT_DESC,
        }},
    },
})
```

### Joins

You can include one or more join definitions using `models.SimpletonGormFindJoins`.

## Delete records

```go
err := sg.Delete(&models.SimpletonGormDelete{
    Type:      enuns.DT_SOFT,
    TableName: "users",
    FieldName: "id",
    FieldValue: user.ID,
    Model:     &User{},
})
```

Available delete types:

- `enuns.DT_SOFT`
- `enuns.DT_PERMANENT`

## Raw SQL

Execute raw SQL commands:

```go
err := sg.SQLExec(&models.SimpletonGormSQL{
    SQL:          "UPDATE users SET name = ? WHERE id = ?",
    FieldsValues: []interface{}{"Bob", user.ID},
})
```

Query raw SQL results:

```go
result, err := sg.SQLFind(&models.SimpletonGormSQL{
    SQL:          "SELECT * FROM users WHERE id = ?",
    FieldsValues: []interface{}{user.ID},
})
```

## Notes

- When using `DB_SQLITEINMEMORY`, `MigrateTables` must be provided.
- For file-based SQLite and SQL servers, `FilePathOrDns` is required.
- `Find` returns `models.SimpletonGormResult` with `Data` as `[]map[string]any` and `Count` as `*uint64`.
