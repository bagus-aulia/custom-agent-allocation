package main

import (
	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/middleware"
	"github.com/bagus-aulia/pondok-sinau-golang/routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.InitDB()
	defer config.DB.Close()

	route := gin.Default()

	api := route.Group("/api/v1/")
	{
		// api.GET("/", routes.Home)
		// api.GET("/dashboard", routes.Home)
		api.GET("/auth/:provider", routes.RedirectHandler)          //done
		api.GET("/auth/:provider/callback", routes.CallbackHandler) //done

		api.POST("/check-username/:role", routes.CheckUsername) //done

		member := api.Group("/member")
		{
			member.GET("/", routes.MemberIndex)                    //done
			member.GET("/profile/:username", routes.MemberProfile) //done

			member.POST("/", middleware.IsAuth(), routes.MemberCreate) //done

			member.PUT("/update/:id", middleware.IsAuth(), routes.MemberUpdate) //done

			member.DELETE("/delete/:id", middleware.IsAuth(), routes.MemberDelete) //done
		}

		admin := api.Group("/admin")
		{
			admin.GET("/", routes.AdminIndex)                    //done
			admin.GET("/profile/:username", routes.AdminProfile) //done

			admin.PUT("/update/:id", middleware.IsAuth(), routes.AdminUpdate)           //
			admin.PUT("/toggle-role/:id", middleware.IsAdmin(), routes.AdminToggleRole) //done
		}

		book := api.Group("/book")
		{
			book.GET("/", routes.BookIndex)              //done
			book.GET("/detail/:code", routes.BookDetail) //done

			book.POST("/", middleware.IsAuth(), routes.BookCreate) //done

			book.PUT("/update/:id", middleware.IsAuth(), routes.BookUpdate) //done

			book.DELETE("/delete/:id", middleware.IsAuth(), routes.BookDelete) //done
		}

		borrow := api.Group("/borrow")
		{
			borrow.GET("/", middleware.IsAuth(), routes.BorrowIndex)              //done
			borrow.GET("/detail/:code", middleware.IsAuth(), routes.BorrowDetail) //done

			borrow.POST("/", middleware.IsAuth(), routes.BorrowCreate) //done

			borrow.PUT("/update/:id", middleware.IsAuth(), routes.BorrowUpdate) //done

			borrow.DELETE("/delete/:id", middleware.IsAuth(), routes.BorrowDelete) //done
		}

		returns := api.Group("/return")
		{
			returns.PUT("/update/:id", middleware.IsAuth(), routes.Return) //done
		}
	}

	route.Run()
}
