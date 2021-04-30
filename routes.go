package main

// initializeRoutes defines the available routes for
// the app and sets the route's handler functions
func initializeRoutes() {
	router.GET("/", handleIndexPage)
}
