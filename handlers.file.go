package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Show home page with list of files
func showIndexPage(c *gin.Context) {
	files := getAllFiles()
	render(c, gin.H{"payload": files}, "index.html")
}

// Delete file by file_id
func deleteFile(c *gin.Context) {
	if fileID, err := strconv.Atoi(c.Param("file_id")); err == nil {
		if file, err := deleteFileByID(fileID); err == nil {
			render(c, gin.H{"payload": file}, "delete_successful.html")
		} else {
			render(c, gin.H{"payload": map[string]string{
				"status": strconv.FormatInt(http.StatusNotFound, 10),
				"reason": "file could not be deleted"}}, "error.html")
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		render(c, gin.H{"payload": map[string]string{
			"status": strconv.FormatInt(http.StatusNotFound, 10),
			"reason": "file not found"}}, "error.html")
		c.AbortWithError(http.StatusNotFound, err)
	}
}

// Show page for deletion
func showFileDeletionPage(c *gin.Context) {
	if fileID, err := strconv.Atoi(c.Param("file_id")); err == nil {
		if file, err := getFileByID(fileID); err == nil {
			render(c, gin.H{"payload": file}, "delete.html")
		} else {
			render(c, gin.H{"payload": map[string]string{
				"status": strconv.FormatInt(http.StatusNotFound, 10),
				"reason": "file not found"}}, "error.html")
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		render(c, gin.H{"payload": map[string]string{
			"status": strconv.FormatInt(http.StatusNotFound, 10),
			"reason": "file not found"}}, "error.html")
		c.AbortWithError(http.StatusNotFound, err)
	}
}

// Fetch points of the file by file_id
func getPoints(c *gin.Context) {
	if fileID, err := strconv.Atoi(c.Param("file_id")); err == nil {
		if file, err := getFileByID(fileID); err == nil {
			if len(file.Points) != 0 {
				render(c, gin.H{"payload": file.Points}, "")
			} else {
				if points, err := parseFile(file); err == nil {
					file.Points = points
					render(c, gin.H{"payload": file.Points}, "")
				} else {
					render(c, gin.H{"payload": map[string]string{
						"status": strconv.FormatInt(http.StatusBadRequest, 10),
						"reason": "error getting points"}}, "error.html")
					c.AbortWithError(http.StatusBadRequest, err)
				}
			}
		} else {
			render(c, gin.H{"payload": map[string]string{
				"status": strconv.FormatInt(http.StatusBadRequest, 10),
				"reason": "error getting points"}}, "error.html")
			c.AbortWithError(http.StatusBadRequest, err)
		}
	}
}

// Fetch file by file_id
func getFile(c *gin.Context) {
	if fileID, err := strconv.Atoi(c.Param("file_id")); err == nil {
		if file, err := getFileByID(fileID); err == nil {
			if len(file.Points) != 0 {
				render(c, gin.H{"payload": file}, "file.html")
			} else if points, err := parseFile(file); err == nil {
				file.Points = points
				render(c, gin.H{"payload": file}, "file.html")
			} else {
				render(c, gin.H{"payload": map[string]string{
					"status": strconv.FormatInt(http.StatusBadRequest, 10),
					"reason": "error processing file"}}, "error.html")
				c.AbortWithError(http.StatusBadRequest, err)
			}

		} else {
			render(c, gin.H{"payload": map[string]string{
				"status": strconv.FormatInt(http.StatusNotFound, 10),
				"reason": "file not found"}}, "error.html")
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		render(c, gin.H{"payload": map[string]string{
			"status": strconv.FormatInt(http.StatusNotFound, 10),
			"reason": "file not found"}}, "error.html")
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// Upload file
func uploadFile(c *gin.Context) {
	newFile, handler, err := c.Request.FormFile("file")
	if err != nil {
		render(c, gin.H{"payload": map[string]string{
			"status": strconv.FormatInt(http.StatusBadRequest, 10),
			"reason": "bad request"}}, "error.html")
		c.AbortWithError(http.StatusBadRequest, err)
	}
	defer newFile.Close()
	var isFound bool
	for _, i := range fileList {
		if i.Name == handler.Filename {
			isFound = true
			break
		}
	}
	if isFound {
		render(c, gin.H{"payload": map[string]string{
			"status": strconv.FormatInt(http.StatusBadRequest, 10),
			"reason": "file " + handler.Filename + " already exist"}}, "error.html")
		c.AbortWithError(http.StatusBadRequest, err)
	} else if _, err := uploadNewFile(newFile, handler); err == nil {
		render(c, gin.H{
			"payload": map[string]string{"Name": handler.Filename}},
			"upload_successful.html")
	} else {
		render(c, gin.H{
			"payload": map[string]string{
				"status": strconv.FormatInt(http.StatusBadRequest, 10),
				"reason": "file could not be uploaded"}}, "error.html")
		c.AbortWithError(http.StatusBadRequest, err)
	}
}
