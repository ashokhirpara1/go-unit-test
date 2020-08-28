package restapi

import (
	"learning/unit-testing/controllers"
	"net/http"
)

func configureAPI(api *operations.ClientAPI) http.Handler {

	api.ConnectionsUsersConnectionsGroupsByUserIDAndGroupIDDeleteHandler = connections.UsersConnectionsGroupsByUserIDAndGroupIDDeleteHandlerFunc(controllers.UsersConnectionsGroupsByUserIDAndGroupIDDeleteController)

	api.ConnectionsUsersConnectionsGroupsByUserIDAndGroupIDGetHandler = connections.UsersConnectionsGroupsByUserIDAndGroupIDGetHandlerFunc(controllers.UsersConnectionsGroupsByUserIDAndGroupIDGetController)

	api.ConnectionsUsersConnectionsGroupsByUserIDAndGroupIDPatchHandler = connections.UsersConnectionsGroupsByUserIDAndGroupIDPatchHandlerFunc(controllers.UsersConnectionsGroupsByUserIDAndGroupIDPatchController)

	api.ConnectionsUsersConnectionsGroupsByUserIDGetHandler = connections.UsersConnectionsGroupsByUserIDGetHandlerFunc(controllers.UsersConnectionsGroupsByUserIDGetController)

	// Create a User's Connection Group.
	api.ConnectionsUsersConnectionsGroupsByUserIDPostHandler = connections.UsersConnectionsGroupsByUserIDPostHandlerFunc(controllers.CreateConnectionsGroupsByUserIDController)

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
