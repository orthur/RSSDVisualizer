package main

func initializeRoutes() {

	router.GET("/", showIndexPage)
	router.POST("/upload", uploadFile)
	router.Static("/static/", "static/")
	fileRoutes := router.Group("/file")
	{
		fileRoutes.GET("/view/:file_id", getFile)
		fileRoutes.GET("/points/:file_id", getPoints)
		fileRoutes.GET("/delete/:file_id", showFileDeletionPage)
		fileRoutes.POST("/delete/:file_id", deleteFile)
	}
}
