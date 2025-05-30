package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Location struct {
	ent.Schema
}

func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("locality_name"),
		field.String("city"),
		field.String("state"),
		field.String("country").Default("India"),
		field.String("pincode"),
		field.String("area_type"), // Sector, Phase, Block, etc.
		field.Float("latitude").Optional(),
		field.Float("longitude").Optional(),
		field.String("google_map_link").Optional(),
		field.Text("location_description").Optional(),
		field.JSON("nearby_landmarks", NearbyLandmarks{}),
		field.JSON("connectivity", LocationConnectivity{}),
		field.Bool("is_active").Default(true),
		field.String("slug"), // URL-friendly version of locality name
	}
}

func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}

// nearby landmarks
type NearbyLandmarks struct {
	Landmarks []struct {
		Name     string `json:"name"`
		Type     string `json:"type"` // hospital, school, mall, metro, etc.
		Distance string `json:"distance"`
		Icon     string `json:"icon"`
	} `json:"landmarks"`
}

// location connectivity
type LocationConnectivity struct {
	Metro []struct {
		StationName string `json:"station_name"`
		Distance    string `json:"distance"`
		Line        string `json:"line"`
	} `json:"metro"`
	Airport struct {
		Name     string `json:"name"`
		Distance string `json:"distance"`
	} `json:"airport"`
	Railway []struct {
		StationName string `json:"station_name"`
		Distance    string `json:"distance"`
	} `json:"railway"`
	Highways []struct {
		HighwayName string `json:"highway_name"`
		Distance    string `json:"distance"`
	} `json:"highways"`
}
