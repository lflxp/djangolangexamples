package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/lflxp/djangolangexamples/demo/service"
)

func RegisterDemo(router *gin.Engine) {
	demoGroup := router.Group("/apis/v1")
	{
		demoGroup.GET("/demo", service.GetAllDemo)
		demoGroup.GET("/:id/demo", service.GetDemoById)
		demoGroup.POST("/demo", service.CreateDemo)
		demoGroup.DELETE("/:id/demo", service.DeleteDemo)
		demoGroup.PUT("/:id/demo", service.PutDemo)
	}
}

func RegisterVpn(router *gin.Engine) {
	demoGroup := router.Group("/apis/v1")
	{
		demoGroup.GET("/vpn", service.GetAllVpn)
		demoGroup.GET("/:id/vpn", service.GetVpnById)
		demoGroup.POST("/vpn", service.CreateVpn)
		demoGroup.DELETE("/:id/vpn", service.DeleteVpn)
		demoGroup.PUT("/:id/vpn", service.PutVpn)
	}
}
