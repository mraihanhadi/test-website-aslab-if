package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataTP struct {
	gorm.Model
	Id        int       `gorm:"primaryKey" json:"id"`
	Judul     string    `json:"judul"`
	SubJudul  string    `json:"subjudul"`
	Kategori  string    `json:"kategori"`
	Tanggal   time.Time `json:"tanggal"`
	Deadline  string    `json:"deadline"`
	Deskripsi string    `json:"deskripsi"`
}

var db *gorm.DB

func connectDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_lab?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db = database

	err = db.AutoMigrate(&DataTP{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func main() {
	connectDatabase()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		var DataTP []DataTP
		db.Find(&DataTP)
		c.HTML(http.StatusOK, "index.html", gin.H{"data": DataTP})
	})
	r.GET("/view", func(c *gin.Context) {
		var DataTP []DataTP
		db.Find(&DataTP)
		c.HTML(http.StatusOK, "view.html", gin.H{"data": DataTP})
	})
	r.POST("/addTP", func(c *gin.Context) {
		Judul := c.PostForm("judul")
		SubJudul := c.PostForm("subJudul")
		Kategori := c.PostForm("kategori")
		Tanggal := time.Now()
		Deadline := c.PostForm("deadline")
		Deskripsi := c.PostForm("deskripsi")
		data := DataTP{
			Judul:     Judul,
			SubJudul:  SubJudul,
			Kategori:  Kategori,
			Tanggal:   Tanggal,
			Deadline:  strings.Replace(Deadline, "T", " ", -1),
			Deskripsi: Deskripsi,
		}
		db.Create(&data)
		c.Redirect(http.StatusFound, "/view")
	})
	r.POST("/update/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
		Judul := c.PostForm("judul")
		SubJudul := c.PostForm("subJudul")
		Kategori := c.PostForm("kategori")
		Tanggal := time.Now()
		Deadline := c.PostForm("deadline")
		Deskripsi := c.PostForm("deskripsi")

		updates := DataTP{
			Judul:     Judul,
			SubJudul:  SubJudul,
			Kategori:  Kategori,
			Tanggal:   Tanggal,
			Deadline:  strings.Replace(Deadline, "T", " ", -1),
			Deskripsi: Deskripsi,
		}
		if err := db.Model(&DataTP{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to update data")
			return
		}
		c.Redirect(http.StatusFound, "/view")
	})

	r.POST("delete/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		if err := db.Unscoped().Delete(&DataTP{}, id).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to delete data")
			return
		}
		c.Redirect(http.StatusFound, "/view")
	})
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
