package main

import (
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Datamahasiswa struct {
	gorm.Model
	Id         int    `gorm:"primaryKey" json:"id"`
	Nama       string `json:"nama"`
	Nim        string `json:"nim"`
	Matakuliah string `json:"matakuliah"`
	Sks        int    `json:"sks"`
	Nilai      int    `json:"nilai"`
	Indeks     string `json:"indeks"`
}

var db *gorm.DB

func connectDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_cruddy?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db = database

	err = db.AutoMigrate(&Datamahasiswa{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func main() {
	connectDatabase()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		var Datamahasiswa []Datamahasiswa
		db.Find(&Datamahasiswa)
		for i := range Datamahasiswa {
			fmt.Println(Datamahasiswa[i].Id)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"data": Datamahasiswa})
	})
	r.GET("/login", func(c *gin.Context) {
		var Datamahasiswa []Datamahasiswa
		db.Find(&Datamahasiswa)
		for i := range Datamahasiswa {
			fmt.Println(Datamahasiswa[i].Id)
		}
		c.HTML(http.StatusOK, "login.html", gin.H{"data": Datamahasiswa})
	})
	r.POST("/sendData", func(c *gin.Context) {
		nama := c.PostForm("nama")
		nim := c.PostForm("nim")
		matakuliah := c.PostForm("matakuliah")
		sks := c.PostForm("sks")
		nilai := c.PostForm("nilai")
		sksInt, sksError := strconv.Atoi(sks)
		nilaiInt, nilaiError := strconv.Atoi(nilai)
		indeks := setIndeks(nilaiInt)
		if sksError != nil || nilaiError != nil {
			c.String(http.StatusBadRequest, "Invalid input for SKS or Nilai")
			return
		}
		data := Datamahasiswa{
			Nama:       nama,
			Nim:        nim,
			Matakuliah: matakuliah,
			Sks:        sksInt,
			Nilai:      nilaiInt,
			Indeks:     indeks,
		}
		db.Create(&data)
		c.Redirect(http.StatusFound, "/login")
	})

	r.POST("/update/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		nama := c.PostForm("nama")
		nim := c.PostForm("nim")
		matakuliah := c.PostForm("matakuliah")
		sksStr := c.PostForm("sks")
		nilaiStr := c.PostForm("nilai")

		sksInt, sksErr := strconv.Atoi(sksStr)
		nilaiInt, nilaiErr := strconv.Atoi(nilaiStr)
		if sksErr != nil || nilaiErr != nil {
			c.String(http.StatusBadRequest, "Invalid input for SKS or Nilai")
			return
		}

		updates := map[string]interface{}{
			"Nama":       nama,
			"Nim":        nim,
			"Matakuliah": matakuliah,
			"Sks":        sksInt,
			"Nilai":      nilaiInt,
			"Indeks":     setIndeks(nilaiInt),
		}

		if err := db.Model(&Datamahasiswa{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to update data")
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	r.POST("/delete/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		if err := db.Unscoped().Delete(&Datamahasiswa{}, id).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to delete data")
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func setIndeks(nilai int) string {
	switch {
	case nilai > 85:
		return "A"
	case nilai > 75:
		return "AB"
	case nilai > 65:
		return "B"
	case nilai > 55:
		return "BC"
	case nilai > 45:
		return "C"
	case nilai > 35:
		return "D"
	default:
		return "E"
	}
}
