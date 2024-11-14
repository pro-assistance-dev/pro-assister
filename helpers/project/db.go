package project

import (
	"context"
	"log"

	"github.com/uptrace/bun"
)

func UpdateSchemasDB(db *bun.DB, schemas Schemas) {
	initSchemaTables(db)

	schemasForInsert := make([]*Schema, 0)
	for _, item := range schemas {
		schemasForInsert = append(schemasForInsert, item)
	}
	insertSchemas(db, schemasForInsert)

	fields := make(SchemaFields, 0)
	for _, item := range schemas {
		fields = append(fields, item.Fields...)
	}
	insertFields(db, fields)
}

func initSchemaTables(db *bun.DB) {
	_, err := db.NewCreateTable().Model((*Schema)(nil)).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.NewCreateIndex().
		Model((*Schema)(nil)).
		IfNotExists().
		Index("name_camel_idx").
		Unique().
		Column("name_camel").
		Exec(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.NewCreateTable().Model((*SchemaField)(nil)).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.NewCreateIndex().
		IfNotExists().
		Model((*SchemaField)(nil)).
		Index("field_unique_paid").
		Unique().
		Column("name_camel", "schema_id").
		Exec(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}

func insertSchemas(db *bun.DB, schemas []*Schema) {
	_, err := db.NewInsert().Model(&schemas).On("CONFLICT (name_camel) do update").
		Set("name_table = EXCLUDED.name_table").
		Set("name_pascal = EXCLUDED.name_pascal").
		Set("name_plural = EXCLUDED.name_plural").
		Set("name_rus = EXCLUDED.name_rus").
		Set("sort_column = EXCLUDED.sort_column").
		Set("label = EXCLUDED.label").
		Set("value = EXCLUDED.value").
		Exec(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}

func insertFields(db *bun.DB, fields SchemaFields) {
	_, err := db.NewInsert().Model(&fields).On("CONFLICT (name_camel, schema_id) do update").
		Set("name_col = EXCLUDED.name_col").
		Set("name_pascal = EXCLUDED.name_pascal").
		Set("name_rus = EXCLUDED.name_rus").
		Set("type = EXCLUDED.type").
		Exec(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
