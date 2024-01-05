package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/lflxp/djangolang"
	"github.com/lflxp/djangolang/middlewares"
	ctls "github.com/lflxp/djangolang/tls"
	"github.com/lflxp/djangolang/utils"
	"github.com/skratchdot/open-golang/open"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/lflxp/djangolangexamples/demo"
	_ "github.com/lflxp/djangolangexamples/docs"
)

var (
	// 是否开启https访问
	isHttps bool = true
	// 设置服务host
	Host string = "0.0.0.0"
	// 设置服务端口
	Port string = "8000"
	// 是否开启自动打开浏览器
	OpenBrowser bool = true
)

func init() {
	// 设置404自动跳转
	middlewares.NoRoutePath = "/swagger/index.html"
}

func main() {
	r := gin.Default()
	// 设置跨域设置
	r.Use(middlewares.Cors())
	// 设置恢复策略
	r.Use(gin.Recovery())
	// 设置日志
	r.Use(gin.Logger())
	// 设置健康检查接口
	r.GET("/health", middlewares.RegisterHealthMiddleware)
	// 设置swagger
	middlewares.RegisterSwaggerMiddleware(r)
	// 注册demo和vpn接口
	demo.RegisterDemo(r)
	demo.RegisterVpn(r)
	// 设置自动跳转
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/admin/index")
	})

	// 注册admin接口
	djangolang.RegisterControllerAdmin(r)
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// 配置gin优雅启动
	var server *http.Server
	if isHttps {
		err := ctls.Refresh()
		if err != nil {
			panic(err)
		}

		pool := x509.NewCertPool()
		caCeretPath := "ca.crt"

		caCrt, err := os.ReadFile(caCeretPath)
		if err != nil {
			panic(err)
		}

		pool.AppendCertsFromPEM(caCrt)

		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", Host, Port),
			Handler: r,
			TLSConfig: &tls.Config{
				ClientCAs:  pool,
				ClientAuth: tls.RequestClientCert,
			},
		}
	} else {
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", Host, Port),
			Handler: r,
		}

	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Warn("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	var openUrl string
	for index, ip := range utils.GetIps() {
		if isHttps {
			log.Infof("Listening and serving HTTPS on https://%s:%s", ip, Port)
		} else {
			log.Infof("Listening and serving HTTPS on http://%s:%s", ip, Port)
		}

		if index == 0 {
			openUrl = fmt.Sprintf("%s:%s", ip, Port)
		}
	}
	if isHttps {
		if OpenBrowser {
			open.Start(fmt.Sprintf("https://%s", openUrl))
		}
		if err := server.ListenAndServeTLS("ca.crt", "ca.key"); err != nil {
			if err == http.ErrServerClosed {
				log.Warn("Server closed under request")
			} else {
				log.Fatalf("Server closed unexpect %s", err.Error())
			}
		}
	} else {
		if OpenBrowser {
			open.Start(fmt.Sprintf("http://%s", openUrl))
		}
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Warn("Server closed under request")
			} else {
				log.Fatalf("Server closed unexpect %s", err.Error())
			}
		}
	}

	log.Warn("Server exiting")
}
