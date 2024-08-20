package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Students struct {
	//gorm.Model - will generate id, created/updated/deleted at fields
	gorm.Model
	Name string `json:"name"`
}

func DBConnection() (*gorm.DB, error) {
	//connects with DB
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return db, err
	}

	//migration of schema
	err = db.AutoMigrate(&Students{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetStudent(ctx *gin.Context) {
	var student Students
	name := ctx.Param("student")

	db, err := DBConnection()
	if err != nil {
		log.Println(err)
	}

	//WHERE action only search for first data equal to the name passed
	if err := db.Where("name= ?", name).First(&student).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Failed": "Student not found"})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

func PostStudent(ctx *gin.Context) {
	var student Students

	//bind -> will keep the http request and convert to student fields
	if err := ctx.ShouldBindJSON(&student); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//after recieve the POST request, parses into newStudent struct
	newStudent := Students{
		Name: student.Name,
	}

	db, err := DBConnection()
	if err != nil {
		log.Println(err)
	}

	//declares the err (create action) and check if there is any error
	if err := db.Create(&newStudent).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Migration error"})
	}

	ctx.JSON(http.StatusOK, student)
}

func UpdateStudent(ctx *gin.Context) {
	var student Students
	name := ctx.Param("student")

	db, err := DBConnection()
	if err != nil {
		log.Println(err)
	}

	//searches student in the DB
	if err := db.Where("name= ?", name).First(&student).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	//convert for student struct
	if err := ctx.ShouldBindJSON(&student); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//update action
	if err := db.Model(&student).Updates(Students{
		Name: student.Name,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

func DeleteStudent(ctx *gin.Context) {
	var student Students
	name := ctx.Param("student")

	db, err := DBConnection()
	if err != nil {
		log.Println(err)
	}

	if err := db.Where("name= ?", name).First(&student).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if err := db.Delete(&student).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
}

func main() {
	r := gin.Default() //creates a server

	//command to get a student at the terminal: curl -X GET "http://localhost:8080/Name"
	r.GET("/:student", GetStudent)

	//command to create a student at the terminal: curl -X POST -H "Content-Type: application/json" -d '{"name": "Name"}' "http://localhost:8080/student"
	r.POST("/student", PostStudent)

	//command to update a student at the terminal: curl -X PUT -H "Content-Type: application/json" -d '{"name": "newName"}' "http://localhost:8080/oldName"
	r.PUT("/:student", UpdateStudent)

	//command to delete at the terminal: curl -X DELETE "http://localhost:8080/Name"
	r.DELETE("/:student", DeleteStudent)

	//runs server at 8080 port (default)
	r.Run()
}
