package migration

import "testing"

func TestMigration(t *testing.T) {
	err := Migrate()
	if err != nil {
		t.Error(err)
	}
}

func TestMigrationOverMigration(t *testing.T) {
	err := Migrate()
	if err != nil {
		t.Error(err)
	}
}
