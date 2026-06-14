package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Users
	rt.router.PUT("/users/:userId", rt.wrap(rt.SetMyUserName))
	rt.router.PUT("/users/:userId/photo", rt.wrap(rt.setMyPhoto))

	// Messages
	rt.router.POST("/messages", rt.wrap(rt.sendMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
