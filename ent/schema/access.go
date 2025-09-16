package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Access holds the schema definition for the Access entity.
type Access struct {
	ent.Schema
}

type Location string

const (
	ESPACIO_URBANO Location = "102"
	CALAMA         Location = "104"
	PACIFICO       Location = "105"
	ARAUCO         Location = "106"
	IQUIQUE        Location = "107"
	ANGAMOS        Location = "108"
)

// Annotations of the User.
func (Access) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "access"},
	}
}

// Fields of the Access.
func (Access) Fields() []ent.Field {
	return []ent.Field{
		field.String("run"),
		field.Enum("location").
			Values(
				string(ESPACIO_URBANO),
				string(CALAMA),
				string(PACIFICO),
				string(ARAUCO),
				string(IQUIQUE),
				string(ANGAMOS),
			),
		field.Time("entry_at"),
		field.Time("exit_at").Optional(),
	}
}

// Edges of the Access.
func (Access) Edges() []ent.Edge {
	return nil
}

func (Access) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("run", "location", "entry_at").Unique(),
	}
}
