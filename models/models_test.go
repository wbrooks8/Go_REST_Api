package models

import (
	"database/sql"
	"testing"
	"time"

	"example.com/REST-api/db"
	_ "github.com/mattn/go-sqlite3"
)

func setupModelTestDB(t *testing.T) {
	t.Helper()

	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	_, err = testDB.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
		CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime DATETIME NOT NULL,
			userID INTEGER
		);
		CREATE TABLE regestrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER,
			user_id INTEGER
		);
	`)
	if err != nil {
		t.Fatalf("create test schema: %v", err)
	}

	db.DB = testDB
	t.Cleanup(func() {
		testDB.Close()
		db.DB = nil
	})
}

func TestUserSaveAndValidateCredentials(t *testing.T) {
	setupModelTestDB(t)

	user := User{Email: "user@example.com", Password: "secret"}
	if err := user.Save(); err != nil {
		t.Fatalf("save user: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected saved user to receive an ID")
	}

	credentials := User{Email: "user@example.com", Password: "secret"}
	if err := credentials.ValidateCredentials(); err != nil {
		t.Fatalf("validate credentials: %v", err)
	}
	if credentials.ID != user.ID {
		t.Fatalf("expected user ID %d, got %d", user.ID, credentials.ID)
	}

	invalidCredentials := User{Email: "user@example.com", Password: "wrong"}
	if err := invalidCredentials.ValidateCredentials(); err == nil {
		t.Fatal("expected invalid credentials to fail")
	}
}

func TestEventCRUDAndRegistration(t *testing.T) {
	setupModelTestDB(t)

	event := Event{
		Name:        "Go meetup",
		Description: "Learn Go",
		Location:    "Online",
		DateTime:    time.Now(),
		UserID:      1,
	}
	if err := event.Save(); err != nil {
		t.Fatalf("save event: %v", err)
	}
	if event.ID == 0 {
		t.Fatal("expected saved event to receive an ID")
	}

	loaded, err := GetEventById(event.ID)
	if err != nil {
		t.Fatalf("load event: %v", err)
	}
	if loaded.Name != event.Name {
		t.Fatalf("expected event name %q, got %q", event.Name, loaded.Name)
	}

	event.Name = "Advanced Go meetup"
	if err := event.Update(); err != nil {
		t.Fatalf("update event: %v", err)
	}
	loaded, err = GetEventById(event.ID)
	if err != nil {
		t.Fatalf("reload event: %v", err)
	}
	if loaded.Name != event.Name {
		t.Fatalf("expected updated event name %q, got %q", event.Name, loaded.Name)
	}

	if err := event.Register(7); err != nil {
		t.Fatalf("register event: %v", err)
	}
	var count int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM regestrations WHERE event_id = ? AND user_id = ?", event.ID, 7).Scan(&count); err != nil {
		t.Fatalf("count registration: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one registration, got %d", count)
	}

	if err := event.CancelRegistration(7); err != nil {
		t.Fatalf("cancel registration: %v", err)
	}
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM regestrations WHERE event_id = ? AND user_id = ?", event.ID, 7).Scan(&count); err != nil {
		t.Fatalf("count cancellation: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no registration, got %d", count)
	}

	if err := event.Delete(); err != nil {
		t.Fatalf("delete event: %v", err)
	}
	_, err = GetEventById(event.ID)
	if err == nil {
		t.Fatal("expected deleted event lookup to fail")
	}
}
