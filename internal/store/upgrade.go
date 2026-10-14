package store

import (
	"database/sql"
	"fmt"
)

// upgradeStats adds the click bucket key and deletion cascade to older databases.
// Dangling history is archived rather than discarded or attached to future links.
func upgradeStats(db *sql.DB) error {
	var keys, cascade int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('Stats') WHERE (name='ID' AND pk=1) OR (name='Created' AND pk=2)`).Scan(&keys); err != nil {
		return err
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_foreign_key_list('Stats') WHERE "table"='Links' AND "from"='ID' AND "to"='ID' AND on_delete='CASCADE'`).Scan(&cascade); err != nil {
		return err
	}
	if keys == 2 && cascade == 1 {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var invalid int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM Stats s JOIN Links l ON l.ID=s.ID WHERE typeof(s.Created)!='integer' OR s.Clicks<0 OR (s.Clicks IS NOT NULL AND typeof(s.Clicks)!='integer')`).Scan(&invalid); err != nil {
		return err
	}
	if invalid > 0 {
		return fmt.Errorf("cannot upgrade click history: %d invalid rows", invalid)
	}
	for _, statement := range []string{
		`CREATE TABLE StatsNext (ID TEXT NOT NULL REFERENCES Links(ID) ON DELETE CASCADE,Created INTEGER NOT NULL,Clicks INTEGER NOT NULL CHECK(Clicks>=0),PRIMARY KEY(ID,Created))`,
		`CREATE TABLE IF NOT EXISTS StatsOrphans_v1 (ID TEXT,Created INTEGER,Clicks INTEGER)`,
		`INSERT INTO StatsOrphans_v1 SELECT s.ID,s.Created,s.Clicks FROM Stats s WHERE NOT EXISTS(SELECT 1 FROM Links l WHERE l.ID=s.ID)`,
		`INSERT INTO StatsNext SELECT s.ID,s.Created,SUM(COALESCE(s.Clicks,0)) FROM Stats s JOIN Links l ON l.ID=s.ID GROUP BY s.ID,s.Created`,
		`DROP TABLE Stats`,
		`ALTER TABLE StatsNext RENAME TO Stats`,
	} {
		if _, err := tx.Exec(statement); err != nil {
			return fmt.Errorf("upgrade click history: %w", err)
		}
	}
	return tx.Commit()
}
