package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iqbal2604/dear-talk-api.git/internal/handler"
)

func registerRoomRoutes(protected *gin.RouterGroup, roomHandler *handler.RoomHandler) {
	rooms := protected.Group("/rooms")
	{
		rooms.POST("", roomHandler.CreateRoom)
		rooms.GET("", roomHandler.GetRooms)
		rooms.GET(":id", roomHandler.GetRoomByID)
		rooms.PUT(":id", roomHandler.UpdateRoom)
		rooms.DELETE(":id", roomHandler.DeleteRoom)

		// Members
		rooms.POST(":id/members", roomHandler.AddMember)
		rooms.DELETE(":id/members/:userId", roomHandler.RemoveMember)
	}
}
