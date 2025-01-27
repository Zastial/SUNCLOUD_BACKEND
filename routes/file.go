package routes

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"SUNCLOUD_BACKEND/model"
)

// FileHandler structure pour contenir les dépendances
type FileHandler struct {
	DB *model.Database
}

func (h *FileHandler) CreateFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun fichier fourni"})
		return
	}
	defer file.Close()

	filePath, err := SaveFile(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileEntry := model.File{
		Name: header.Filename,
		Path: filePath,
		Size: header.Size,
		Type: header.Header.Get("Content-Type"),
	}
	if err := h.DB.Create(&fileEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, fileEntry)
}

// GetFile récupère un fichier par son ID
func (h *FileHandler) GetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var file model.File
	if err := h.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fichier non trouvé"})
		return
	}

	c.JSON(http.StatusOK, file)
}

// GetAllFiles récupère tous les fichiers
func (h *FileHandler) GetAllFiles(c *gin.Context) {
	var files []model.File
	if err := h.DB.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var file model.File
	if err := h.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fichier non trouvé"})
		return
	}

	if err := os.Remove(file.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du fichier"})
		return
	}

	if err := h.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fichier supprimé avec succès"})
}
