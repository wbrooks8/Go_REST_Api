package routes

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/REST-api/db"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func setupRegisterTestDB(t *testing.T, schema string) {
	t.Helper()

	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	if _, err := testDB.Exec(schema); err != nil {
		testDB.Close()
		t.Fatalf("create test schema: %v", err)
	}

	db.DB = testDB
	t.Cleanup(func() {
		testDB.Close()
		db.DB = nil
	})
}

func registerTestContext(t *testing.T, eventID string, userID int64) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest(http.MethodPost, "/events/"+eventID+"/register", nil)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: eventID}}
	context.Set("userId", userID)

	return context, recorder
}

func TestRegisterForEventReturnsBadRequestForInvalidEventID(t *testing.T) {
	context, recorder := registerTestContext(t, "not-an-id", 1)

	registerForEvent(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
	if recorder.Body.String() != `{"message":"Could not parse event id."}` {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}

func TestRegisterForEventReturnsInternalServerErrorWhenEventCannotBeFetched(t *testing.T) {
	setupRegisterTestDB(t, "")
	context, recorder := registerTestContext(t, "1", 1)

	registerForEvent(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
	if recorder.Body.String() != `{"message":"Could not fetch event."}` {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}

func TestRegisterForEventCreatesRegistration(t *testing.T) {
	setupRegisterTestDB(t, `
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
	_, err := db.DB.Exec(
		"INSERT INTO events (name, description, location, dateTime, userID) VALUES (?, ?, ?, ?, ?)",
		"Go meetup", "Learn Go", "Online", time.Now(), 1,
	)
	if err != nil {
		t.Fatalf("insert test event: %v", err)
	}

	context, recorder := registerTestContext(t, "1", 42)
	registerForEvent(context)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}
	if recorder.Body.String() != `{"message":"Registered!"}` {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}

	var registrationCount int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM regestrations WHERE event_id = 1 AND user_id = 42").Scan(&registrationCount); err != nil {
		t.Fatalf("query registration: %v", err)
	}
	if registrationCount != 1 {
		t.Fatalf("expected one registration, got %d", registrationCount)
	}
}

func TestRegisterForEventReturnsInternalServerErrorWhenRegistrationFails(t *testing.T) {
	setupRegisterTestDB(t, `
		CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime DATETIME NOT NULL,
			userID INTEGER
		);
	`)
	_, err := db.DB.Exec(
		"INSERT INTO events (name, description, location, dateTime, userID) VALUES (?, ?, ?, ?, ?)",
		"Go meetup", "Learn Go", "Online", time.Now(), 1,
	)
	if err != nil {
		t.Fatalf("insert test event: %v", err)
	}

	context, recorder := registerTestContext(t, "1", 42)
	registerForEvent(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
	if recorder.Body.String() != `{"message":"Could not register user for event."}` {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}

func TestCancelRegistrationReturnsBadRequestForInvalidEventID(t *testing.T) {
	context, recorder := registerTestContext(t, "not-an-id", 1)

	cancelRegistration(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
	if recorder.Body.String() != `{"message":"Could not parse event id."}` {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}
