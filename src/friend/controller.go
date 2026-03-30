package friend

import (
	"net/http"
	"ws/src/auth"
	"ws/src/notify"
	"ws/src/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Controller struct {
	Repo     *Repository
	UserRepo *user.Repository
}

func NewController(r *Repository, ur *user.Repository) *Controller {
	return &Controller{Repo: r, UserRepo: ur}
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

	req, err := ctrl.Repo.GetRequestByID(requestObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found",
		})
		return
	}

	if err := ctrl.Repo.AcceptRequest(requestObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can not Accept friend request"})
		return
	}

	notify.SendToUser(req.FromUserID.Hex(), "Lời mời kết bạn đã được xác nhận !")
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

	req, err := ctrl.Repo.GetRequestByID(requestObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Request not found",
		})
		return
	}

	if err := ctrl.Repo.RefuseRequest(requestObjID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot refuse request",
		})
		return
	}

	// Gửi notify cho người nhận
	notify.SendToUser(req.FromUserID.Hex(), "Lời mời kết bạn đã bị từ chối!")
	c.JSON(http.StatusOK, gin.H{
		"message": "Friend Request Rejected!",
	})
}

func (ctrl *Controller) ListMyFriend(c *gin.Context) {
	userIDStr := c.MustGet(auth.UserIDKey).(string)
	userID, _ := bson.ObjectIDFromHex(userIDStr)

	friendIDs, _ := ctrl.Repo.ListFriends(userID)

	if len(friendIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"friends": []gin.H{},
		})
		return
	}

	users := ctrl.UserRepo.FindManyByID(friendIDs)

	result := make([]gin.H, len(users))

	for _, u := range users {
		result = append(result, gin.H{
			"id":       u.ID.Hex(),
			"username": u.Username,
			"email":    u.Email,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"friends": result,
	})
}
