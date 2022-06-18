package handlers

import (
	"fmt"
	"net/http"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserGroups(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	groups, err := s.DB.GetUserGroups(userUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if len(groups) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, groups)

}

func (s *Server) CreateGroup(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
	}

	var group models.Group
	err = c.ShouldBindJSON(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if group.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad name"})
		return
	}
	if group.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad description"})
		return
	}

	group, err = s.DB.CreateGroup(userUID, group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.actionChan <- &communication.Action{Group: group.ID, User: userUID, Action: "CREATE_GROUP"}

	c.JSON(http.StatusCreated, group)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	userID := c.GetString("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}

	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid group ID"})
		return
	}

	member, err := s.DB.GetUserGroupMember(uid, groupUID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}
	if !member.Creator {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}

	group, err := s.DB.DeleteGroup(groupUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	s.actionChan <- &communication.Action{Group: group.ID, Action: "DELETE_GROUP"}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}

func (s *Server) SetGroupProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid ID"})
		return
	}

	imageFileHeader, err := c.FormFile("avatarFile")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"err": fmt.Sprintf("Max request body size is %v bytes\n", s.maxBodyBytes),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if imageFileHeader == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no file provided"})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")
	if !isAllowedImageType(mimeType) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "image extention not allowed"})
		return
	}

	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid group ID"})
		return
	}

	member, err := s.DB.GetUserGroupMember(userUID, groupUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if !member.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
		return
	}

	pictureURL, err := s.DB.GetGroupProfilePicture(groupUID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		return
	}
	wasEmpty := false
	if pictureURL == "" {
		pictureURL = uuid.New().String()
		wasEmpty = true
	}

	file, err := imageFileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "bad image"})
		return
	}

	if err = s.Storage.UpdateProfilePicture(file, pictureURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if wasEmpty {
		if err = s.DB.SetGroupProfilePicture(groupUID, pictureURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"newUrl": pictureURL})
}

func (s *Server) DeleteGroupProfilePicture(c *gin.Context) {
	userID := c.GetString("userID")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid ID"})
		return
	}
	groupID := c.Param("groupID")
	groupUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid grouo ID"})
		return
	}
	member, err := s.DB.GetUserGroupMember(userUID, groupUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if !member.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
	}
	url, err := s.DB.GetGroupProfilePicture(groupUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group has no image to delete"})
		return
	}
	if err = s.DB.SetGroupProfilePicture(groupUID, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err = s.Storage.DeleteProfilePicture(url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
