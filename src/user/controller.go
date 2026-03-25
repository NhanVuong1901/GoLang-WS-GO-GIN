package user

import (
	"net/http"
	"ws/src/common"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Controller struct {
	Repo *Repository
}

func NewController(repo *Repository) *Controller {
	return &Controller{Repo: repo}
}

func (ctrl Controller) Register(ctx *gin.Context) {
	var input User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	hashed, _ := common.HashPassword(input.Password)
	input.Password = hashed

	// Call Repo
	if err := ctrl.Repo.Create(&input); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username or email existed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Register success!"})
}

func (ctrl Controller) Update(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	userObjID, _ := bson.ObjectIDFromHex(userID)

	var input struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	update := bson.M{}

	if input.Username != nil {
		update["username"] = *input.Username
	}

	if input.Password != nil {
		hashed, _ := common.HashPassword(*input.Password)
		update["password"] = hashed
	}

	if len(update) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nothing to update"})
		return
	}

	if err := ctrl.Repo.UpdateByID(userObjID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}
