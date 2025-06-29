// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
)

// 統一註冊所有路由
func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api")

	AuthRoutes(api)
	TodoRoutes(api)
	MemberRoutes(api)
	// 其他模組路由也可以在這邊加
}
