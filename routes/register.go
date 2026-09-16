package routes

import (
	"net/http"
	"strconv"

	"example.com/REST-api/models"
	"github.com/gin-gonic/gin"
)

// Handler notes:
//  1. Middleware runs before this function and places authenticated data in the
//     Gin context.
//  2. The handler validates input, calls the model layer, and chooses an HTTP
//     response. Keeping database work in models keeps routes focused on HTTP.
//  3. Return after writing an error response so the handler does not try to send
//     a second response.
//
// registerForEvent handles POST requests that register the authenticated user
// for the event identified by the route parameter, such as /events/:id/register.
func registerForEvent(context *gin.Context) {
	// The auth middleware stores the user's ID in the Gin context before this
	// handler runs. GetInt64 returns the stored value for use by the model layer.
	userId := context.GetInt64("userId")

	// Route parameters are strings, so convert the event ID to int64 before
	// passing it to the database/model functions.
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		// A 400 response tells the client that its request contained an invalid
		// event ID. The return prevents the handler from continuing with bad data.
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	// Fetch the event first so Register can associate the user with a real event.
	event, err := models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	// Register performs the model/database operation using the authenticated
	// user's ID and the event that was fetched above.
	err = event.Register(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}

	// 201 Created indicates that the registration was successfully created.
	context.JSON(http.StatusCreated, gin.H{"message": "Registered!"})
}

// TODO: Implement the DELETE endpoint that removes a user's event registration.
// Unlike POST registration, DELETE should normally return 200 (with a message)
// or 204 (with no response body) after the row is removed.
func cancelRegistration(context *gin.Context) {
	// The auth middleware stores the user's ID in the Gin context before this
	// handler runs. GetInt64 returns the stored value for use by the model layer.
	userId := context.GetInt64("userId")

	// Route parameters are strings, so convert the event ID to int64 before
	// passing it to the database/model functions.
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration!"})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{"message": "Registration canceled!"})

}
