package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// ChangePassword
func (s *Server) ChangePassword(c *gin.Context) {
	id := c.GetInt("userID")

	load := struct {
		NewPassword string `json:"newPassword"`
		OldPassword string `json:"oldPassword"`
	}{}
	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := s.DB.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if !checkPassword(user.Pass, load.OldPassword) {
		c.JSON(http.StatusForbidden, gin.H{"err": "Wrong password"})
		return
	}

	hash, err := hashPassword(load.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if err := s.DB.SetPassword(user.ID, hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

///////////////////////////////////////////////////////////////////////////
// UpdateProfilePicture

func (s *Server) UpdateProfilePicture(c *gin.Context) {
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

	pictureURL, err := s.DB.GetProfilePictureURL(uint(userID))
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
		if err = s.DB.SetProfilePicture(uint(userID), pictureURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"newUrl": pictureURL})
}

///////////////////////////////////////////////////////////////////////////
// UpdateProfilePicture

func (s *Server) DeleteProfilePicture(c *gin.Context) {
	userID := c.GetInt("userID")

	url, err := s.DB.GetProfilePictureURL(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user has no image to delete"})
		return
	}
	if err = s.DB.SetProfilePicture(uint(userID), ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err = s.Storage.DeleteProfilePicture(url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
