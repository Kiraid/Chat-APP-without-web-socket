package routes

import (
	"net/http"
	"strconv"
	"time"

	"example.com/chat/models"
	"github.com/gin-gonic/gin"
)


func ChatPage(c *gin.Context) {
	// Add authentication check here
	userID, err := c.Cookie("user_id")
	if err != nil || userID == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	c.HTML(http.StatusOK, "chatroom.html", nil)
}

func CreateChannel(c *gin.Context) {
	var request struct{
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message":"Invalid Request"})
		return
	}
	channel := models.Channel{
		Name: request.Name,
	}
	if err := channel.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success":false, "message":err.Error()})
		return
}
	userIDstr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success":false, "message":"Failed to get user ID"})
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success":false, "message":"Invalid user ID"})
		return
	}
	if err := models.AddUsertoChannel(userID, channel.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Success":false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func joinChannel(c *gin.Context) {
    var req struct {
        Name string `json:"name"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
        return
    }

    channel, err := models.GetChannelByName(req.Name)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Channel not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "channel": channel})
}




func GetMessages(c *gin.Context) {
	channelName := c.Param("channel")
	channelID, err := models.GetChannelIDByName(channelName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}
	messages, err := models.GetMessagebyChannel(channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Couldnt get messages"})
		return
	}
	var messagesWithUsers []struct {
        Content string `json:"content"`
        User    string `json:"user"`
    }
	for _, message := range messages {
        user, err := models.GetUserByID(message.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user details"})
            return
        }
		messagesWithUsers = append(messagesWithUsers, struct {
            Content string `json:"content"`
            User    string `json:"user"`
        }{
            Content: message.Content,
            User:    user.Username,
        })
    }

	c.JSON(http.StatusOK, gin.H{"messages": messagesWithUsers})}

func PutMessage(c *gin.Context){	
	channelName := c.Param("channel")
	var request struct{
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	channelID, err := models.GetChannelIDByName(channelName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	userIDstr, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(userIDstr)
	message := models.Message {
		ChannelID: channelID,
		UserID: userID,
		Content: request.Content,
		Timestamp: time.Now(),
	}
	if err := models.SaveMessage(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}	