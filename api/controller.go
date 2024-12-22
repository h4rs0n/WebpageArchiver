package api

import (
	"WebpageArchiver/assets"
	"WebpageArchiver/common"
	"WebpageArchiver/search"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

func SearchByKeyword(c *gin.Context) {
	queryString := c.Query("q")
	queryPage := c.Query("p")
	pageNum := int64(0)
	if queryPage == "" {
		pageNum = 1
	} else {
		var err error
		pageNum, err = strconv.ParseInt(queryPage, 10, 64)
		if err != nil {
			c.JSON(403, gin.H{
				"Status":  "0",
				"Message": "参数 p 格式错误",
			})
			return
		}
	}
	if queryString == "" {
		c.JSON(403, gin.H{
			"Status":  "0",
			"Message": "缺少关键参数 q",
		})
		return
	}
	queryResult, pageAndHits := search.QueryByKeyword(queryString, pageNum)

	if queryResult == "Error" {
		c.JSON(500, gin.H{
			"Status":  "0",
			"Message": "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":     "1",
		"Message":    "",
		"Result":     queryResult,
		"TotalHits":  pageAndHits["TotalHits"],
		"TotalPages": pageAndHits["TotalPages"],
	})
}

func AddDocByURL(c *gin.Context) {
	// todo
}

func AddHTMLFile(c *gin.Context) {
	htmlFile, err := c.FormFile("file")
	// 上传到临时目录
	filePath := common.ARCHIVEFILELOACTION + "/Temporary/" + htmlFile.Filename
	if err != nil {
		c.JSON(500, gin.H{
			"Status":  "0",
			"Message": "上传文件失败",
		})
		return
	}

	if err := c.SaveUploadedFile(htmlFile, filePath); err != nil {
		c.JSON(500, gin.H{
			"Status":  "0",
			"Message": "上传文件失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "1",
		"Message": "文件上传成功",
	})
}

type File struct {
	Uid      string   `json:"uid"`
	File     struct{} `json:"file"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Percent  int      `json:"percent"`
	Response struct {
		Message string `json:"Message"`
		Status  string `json:"Status"`
	} `json:"response"`
}

type AddDocRequest struct {
	Domain string `json:"domain"`
	Files  []File `json:"files"`
}

func AddDocByHTMLFile(c *gin.Context) {
	var req AddDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(403, gin.H{
			"Status":  "0",
			"Message": "请求参数错误",
		})
		return
	}
	if req.Domain == "" {
		c.JSON(403, gin.H{
			"Status":  "0",
			"Message": "请求参数错误",
		})
		return
	}
	if len(req.Files) != 1 {
		c.JSON(403, gin.H{
			"Status":  "0",
			"Message": "仅支持单个文件上传",
		})
		return
	}
	if req.Files[0].Name == "" {
		c.JSON(403, gin.H{
			"Status":  "0",
			"Message": "请求参数错误",
		})
		return
	}

	if err := search.AddDocFile(req.Files[0].Name, req.Domain); err != nil {
		c.JSON(500, gin.H{
			"Status":  "0",
			"Message": "上传文件失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "1",
		"Message": "文件上传成功",
	})
	return
}

var Templates embed.FS

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func WebStarter(debugMode bool) {
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	if debugMode {
		router.Use(CORSMiddleware())
	}
	router.Static("/static", "./static/web/")
	router.Static("/archive", common.ARCHIVEFILELOACTION)
	router.StaticFS("/assets", http.FS(assets.LoadFile()))

	router.GET("/api/search", SearchByKeyword)
	router.POST("/api/uploadHtmlFile", AddHTMLFile)
	router.POST("/api/upload", AddDocByHTMLFile)

	tmpl := template.Must(template.New("").ParseFS(assets.WebFiles, "web/*.html"))
	router.SetHTMLTemplate(tmpl)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	err := router.Run("0.0.0.0:7845")
	if err != nil {
		fmt.Print("Maybe the port is already in use. Please check it.")
		return
	}
}
