package main

import (
	"fmt"
	"food-delivery/component/appctx"
	"food-delivery/middleware"
	"food-delivery/module/restaurant/transport/ginrestaurant"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

func main() {
	fmt.Println("hello")
	//dsn := "food_delivery:sapassword@tcp(127.0.0.1:3306)/food_delivery?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("MYSQL_CONN_STRING")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db = db.Debug()
	appContext := appctx.NewAppContext(db)
	r := gin.Default()
	r.Use(middleware.Recover(appContext))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")
	restaurants := v1.Group("/restaurants")
	restaurants.POST("/", ginrestaurant.CreateRestaurant(appContext))
	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var data Restaurant
		db.Where("id = ?", id).First(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})
	restaurants.GET("", ginrestaurant.ListRestaurant(appContext))
	restaurants.PATCH("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var data RestaurantUpdate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		db.Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})
	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))
	r.Run()

	//newRestaurant := Restaurant{Name: "nam", Addr: "ho chi minh"}
	//if err := db.Create(&newRestaurant).Error; err != nil {
	//	log.Println(err)
	//}
	//log.Println("new id:", newRestaurant.Id)

	//var myRestaurant Restaurant

	//if err := db.Where("id=?", 3).First(&myRestaurant).Error; err != nil {
	//	log.Println(err)
	//}
	//log.Println(myRestaurant)

	//newName := ""
	//updateData := RestaurantUpdate{Name: &newName}
	//if err := db.Where("id=?", 4).Updates(&updateData).Error; err != nil {
	//	log.Println(err)
	//}
	//log.Println(myRestaurant)

	//myRestaurant.Name = "200 nam"
	//if err := db.Table(Restaurant{}.TableName()).Where("id=?", 1).Delete(nil).Error; err != nil {
	//	log.Println(err)
	//}
	//log.Println(myRestaurant)

}
