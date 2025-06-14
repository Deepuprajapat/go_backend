package migration

import (
	"context"
	"fmt"
	"project-schema/ent"

	_ "github.com/go-sql-driver/mysql"
)

// NewNewDBConnection creates a new connection using ent client
func NewNewDBConnection() (*ent.Client, error) {
	// Connect using ent
	client, err := ent.Open("mysql", "root:password@tcp(localhost:3306)/mydb?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("error connecting to new database: %v", err)
	}

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}
