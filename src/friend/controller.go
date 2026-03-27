package friend

import (
	"net/http"
	"ws/src/auth"
	"ws/src/notify"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Controller struct {
	Repo *Repository
}

func NewController(r *Repository) *Controller {
	return &Controller{Repo: r}
}

func (ctrl *Controller) SendRequest(c *gin.Context) {
	var input struct {
		ToUserID string `json:"to_user_id"`
	}

	c.BindJSON(&input)

	fromID := c.MustGet(auth.UserIDKey).(string)
	fromObjID, _ := bson.ObjectIDFromHex(fromID)
	toObjiD, _ := bson.ObjectIDFromHex(input.ToUserID)

	if err := ctrl.Repo.SendRequest(fromObjID, toObjiD); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not send friend request"})
		return
	}

	// Gửi notify cho người nhận
	notify.SendToUser(toObjiD.Hex(), "Bạn có một lời mời kết bạn !")

	c.JSON(http.StatusOK, gin.H{"message": "Friend Request Sent !"})

}

func (ctrl *Controller) AcceptRequest(c *gin.Context) {
	var input struct {
		RequestID string `json:"request_id"`
	}

	c.BindJSON(&input)

	requestObjID, _ := bson.ObjectIDFromHex(input.RequestID)

	if err := ctrl.Repo.AcceptRequest(requestObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not Accept friend request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend Request Accepted !"})
}

func (ctrl *Controller) ListFriends(c *gin.Context) {
	userID := c.MustGet(auth.UserIDKey).(string)
	userObjID, _ := bson.ObjectIDFromHex(userID)

	ids, err := ctrl.Repo.ListFriends(userObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot get friends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"friends": ids,
	})
}

func (ctrl *Controller) RefuseRequest(c *gin.Context) {
	var input struct {
		RequestID string `json:"request_id"`
	}

	c.BindJSON(&input)

	requestObjID, _ := bson.ObjectIDFromHex(input.RequestID)

	if err := ctrl.Repo.RefuseRequest(requestObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot refuse request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Friend Request Rejected!",
	})
}
