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
	rt.router.PUT("/users/:userId", rt.authWrap(rt.SetMyUserName))
	rt.router.PUT("/users/:userId/photo", rt.authWrap(rt.setMyPhoto))

	// Messages
	rt.router.POST("/messages", rt.authWrap(rt.sendMessage))
	rt.router.PUT("/messages/:messageId", rt.authWrap(rt.updateMessage))
	rt.router.DELETE("/messages/:messageId", rt.authWrap(rt.deleteMessage))
	rt.router.POST("/messages/:messageId/forward", rt.authWrap(rt.forwardMessage))

	// Chats
	rt.router.POST("/chats", rt.authWrap(rt.createChat))

	// Comments
	rt.router.POST("/messages/:messageId/comments", rt.authWrap(rt.createComment))
	rt.router.DELETE("/messages/:messageId/comments", rt.authWrap(rt.deleteComment))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
