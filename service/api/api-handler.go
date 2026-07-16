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
	rt.router.PUT("/messages/:messageId/status", rt.authWrap(rt.updateMessageStatus))
	rt.router.POST("/messages/:messageId/forward", rt.authWrap(rt.forwardMessage))

	// Chats
	rt.router.POST("/chats", rt.authWrap(rt.createChat))
	rt.router.GET("/chats", rt.authWrap(rt.getMyConversations))
	rt.router.GET("/chats/:chatId", rt.authWrap(rt.getConversation))

	// Comments
	rt.router.POST("/messages/:messageId/comments", rt.authWrap(rt.createComment))
	rt.router.DELETE("/messages/:messageId/comments/:commentId", rt.authWrap(rt.deleteComment))

	// Groups
	rt.router.POST("/groups", rt.authWrap(rt.createGroup))
	rt.router.PUT("/groups/:groupId", rt.authWrap(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.authWrap(rt.setGroupPhoto))
	rt.router.POST("/groups/:groupId/members", rt.authWrap(rt.addToGroup))
	rt.router.DELETE("/groups/:groupId/members/:memberId", rt.authWrap(rt.leaveGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
