package routers

import (
	"gin-blog/middleware/jwt"
	"gin-blog/pkg/export"
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/upload"
	"gin-blog/routers/api"
	v1 "gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	//
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	// 获取图片
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	// auth 认证
	r.POST("/auth", api.GetAuth)
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 获取标签
		apiv1.GET("/tags", v1.GetTags)
		// 添加标签
		apiv1.POST("/tags", v1.AddTag)
		// 修改标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		// 删除标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 更新指定文章
		apiv1.PUT("/article/:id", v1.EditArticle)
		// 删除文章
		apiv1.DELETE("/article/:id", v1.DeleteArticle)

		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
		// 导出标签
		r.POST("/tags/export", v1.ExportTag)
		// 导入
		r.POST("/tags/import", v1.ImportTag)
	}
	return r
}
