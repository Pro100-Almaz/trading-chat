package utils

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.Encode(obj)
}

func MigrateDB(db *sqlx.DB) {
	// Create users table if not exists
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		google_id VARCHAR(255) DEFAULT '',
		avatar_emoji INTEGER DEFAULT 0,
		name VARCHAR(255) DEFAULT '',
		password VARCHAR(255) DEFAULT '',
		email VARCHAR(255) NOT NULL UNIQUE,
		phone VARCHAR(255) DEFAULT '',
		is_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP
		);
	`)

	// Create verification_codes table
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS verification_codes (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	// Create index on verification_codes for faster lookups
	db.MustExec(`CREATE INDEX IF NOT EXISTS idx_verification_codes_user_id ON verification_codes(user_id)`)

	// Migration: rename profile_picture to avatar_emoji if old column exists
	var columnExists bool
	err := db.Get(&columnExists, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'users' AND column_name = 'profile_picture'
		)
	`)
	if err != nil {
		log.Error("Error checking for profile_picture column: ", err)
		return
	}

	if columnExists {
		log.Info("Migrating: dropping profile_picture column and adding avatar_emoji")
		db.MustExec(`ALTER TABLE users DROP COLUMN IF EXISTS profile_picture`)
		db.MustExec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_emoji INTEGER DEFAULT 0`)
	}

	// Migration: add is_verified column if not exists
	var isVerifiedExists bool
	err = db.Get(&isVerifiedExists, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'users' AND column_name = 'is_verified'
		)
	`)
	if err != nil {
		log.Error("Error checking for is_verified column: ", err)
		return
	}

	if !isVerifiedExists {
		log.Info("Migrating: adding is_verified column")
		db.MustExec(`ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT FALSE`)
	}
}

func SetCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
}
