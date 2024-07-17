// routes/register.go
package routes

import (
	"log"
	"net/http"
	"strconv"
	"example.com/chat/models"
	"github.com/gin-gonic/gin"
)


func getEvents(context *gin.Context) {
	events, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch events"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}
	log.Printf("User data: %+v\n", user)

	if err := user.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}


func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := models.Authenticate(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	// For simplicity, we'll just set a cookie here
	c.SetCookie("user_id", strconv.Itoa(user.ID), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/chat")
}