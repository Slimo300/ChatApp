package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) GetUserGroups(c *gin.Context) {
	userID := c.GetInt("userID")

	groups, err := s.DB.GetUserGroups(uint(userID))
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
	userID := c.GetInt("userID")

	var group models.Group
	err := c.ShouldBindJSON(&group)
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

	group, err = s.DB.CreateGroup(uint(userID), group.Name, group.Desc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	s.actionChan <- &communication.Action{Group: int(group.ID), User: int(userID), Action: "CREATE_GROUP"}

	c.JSON(http.StatusCreated, group)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	userID := c.GetInt("userID")

	groupID := c.Param("groupID")
	groupIDint, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group not specified"})
		return
	}

	member, err := s.DB.GetUserGroupMember(uint(userID), uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}
	if !member.Creator {
		c.JSON(http.StatusForbidden, gin.H{"err": "couldn't delete group"})
		return
	}

	group, err := s.DB.DeleteGroup(uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	s.actionChan <- &communication.Action{Group: int(group.ID), Action: "DELETE_GROUP"}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}

func (s *Server) SetGroupProfilePicture(c *gin.Context) {
	userID := c.GetInt("userID")

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
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "group not specified"})
		return
	}
	groupIDint, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	member, err := s.DB.GetUserGroupMember(uint(userID), uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if !member.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
		return
	}

	pictureURL, err := s.DB.GetGroupProfilePicture(uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"err": err.Error()})
		return
	}
	if pictureURL == "" {
		pictureURL = uuid.New().String()
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

	if err = s.DB.SetGroupProfilePicture(uint(groupIDint), pictureURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newUrl": pictureURL})
}

func (s *Server) DeleteGroupProfilePicture(c *gin.Context) {
	userID := c.GetInt("userID")
	groupID := c.Param("groupID")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "groupID not specified"})
		return
	}
	groupIDint, err := strconv.Atoi(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	member, err := s.DB.GetUserGroupMember(uint(userID), uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if !member.Setting {
		c.JSON(http.StatusForbidden, gin.H{"err": "no rights to set"})
	}
	url, err := s.DB.GetGroupProfilePicture(uint(groupIDint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user has no image to delete"})
		return
	}
	if err = s.DB.SetGroupProfilePicture(uint(groupIDint), ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err = s.Storage.DeleteProfilePicture(url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
