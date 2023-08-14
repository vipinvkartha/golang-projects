package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name        string `gorm:"not null;default:null"`
	Designation string `gorm:"not null;default:null"`
}

var (
	db  *gorm.DB
	err error
)

func initDB() {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&Employee{})
}

func main() {
	initDB()
	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()
	r := gin.Default()

	r.POST("/employee", createEmployee)
	r.POST("/employees", createEmployees)
	r.GET("/employees", getEmployees)
	r.GET("/employees/:id", getEmployee)
	r.PUT("/employees/:id", updateEmployee)
	r.DELETE("/employees/:id", deleteEmployee)

	fmt.Println("Server started on :8080")
	r.Run(":8080")
}

// createEmployee add a new employee
func createEmployee(c *gin.Context) {
	var employee Employee
	// check for errors
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// if employee.Name == "" || employee.Designation == "" {
	// 	c.JSON(400, gin.H{"error": "Name and Designation are required"})
	// 	return
	// }
	emp := db.Create(&employee)
	if emp.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": emp.Error.Error()})
		return
	}
	c.JSON(200, employee)
}

// getEmployees fetch all employees
func getEmployees(c *gin.Context) {
	var employees []Employee
	db.Find(&employees)
	c.JSON(200, employees)
}

// getEmployee fetch single employee by id
func getEmployee(c *gin.Context) {
	var employee Employee
	id := c.Param("id")
	emp := db.First(&employee, id)
	if emp.Error != nil {
		c.JSON(400, gin.H{"error": emp.Error.Error()})
		return
	}
	c.JSON(200, employee)
}

// putEmployee update employee
func updateEmployee(c *gin.Context) {
	var employee Employee
	id := c.Param("id")
	emp := db.First(&employee, id)
	// add error handling
	if emp.Error != nil {
		c.JSON(400, gin.H{"error": emp.Error})
		return
	}
	// update employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, employee)
}

func deleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee Employee
	emp := db.First(&employee, id)
	if emp.Error != nil {
		c.JSON(400, gin.H{"error": emp.Error.Error()})
		return
	}
	db.Delete(&employee)
	c.JSON(200, gin.H{"message": "Employee deleted successfully"})
}

// batch create employees
func createEmployees(c *gin.Context) {
	var employees []Employee
	if err := c.ShouldBindJSON(&employees); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	emp := db.Create(&employees)
	if emp.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": emp.Error.Error()})
		return
	}
	c.JSON(200, employees)
}
