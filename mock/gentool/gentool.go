package main

import (
	"fmt"
	_const "gitee.com/liujit/shop/mock/const"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

func main() {
	db, _ := gorm.Open(mysql.Open(_const.Database.Source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	g := gen.NewGenerator(gen.Config{
		OutPath:           "../server/lib/data/query",
		OutFile:           "",
		ModelPkgPath:      "models",
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	m := make(map[string]func(columnType gorm.ColumnType) (dataType string))
	m["tinyint"] = func(columnType gorm.ColumnType) (dataType string) {
		c, ok := columnType.ColumnType()
		if ok {
			if c == "tinyint(1)" || c == "tinyint(1) unsigned" {
				return "*bool"
			}
		}
		return "int32"
	}
	g.WithDataTypeMap(m)

	models, err := genModels(g, db, nil)
	if err != nil {
		log.Fatalln("get tables info fail:", err)
	}

	g.ApplyBasic(models...)

	g.Execute()
}

func genModels(g *gen.Generator, db *gorm.DB, tables []string) (models []interface{}, err error) {
	if len(tables) == 0 {
		// Execute tasks for all tables in the database
		tables, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}

	// Execute some data table tasks
	models = make([]interface{}, len(tables))
	for i, tableName := range tables {
		models[i] = g.GenerateModel(tableName)
	}
	return models, nil
}
