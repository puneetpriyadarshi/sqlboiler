package main

import (
	// Log items to the terminal

	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"root/models"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

const (
	host   = "localhost"
	port   = 8080
	user   = "postgres"
	pass   = "Singhasan26!"
	dbname = "postgres"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Successfully connected")

	return db, nil
}

// func CreateStudent(c *gin.Context) {
// 	var student models.Student
// 	err := c.BindJSON(&student)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ID := student.NAME
// 	Firstname := student.Firstname
// 	insertError := student.Insert(&models.Student{
// 		ID:        ID,
// 		Firstname: Firstname,
// 	})
// 	if insertError != nil {
// 		log.Printf("Error while inserting new tenant into db, Reason: %v\n", insertError)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  http.StatusInternalServerError,
// 			"message": "Something went wrong",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"status":  http.StatusCreated,
// 		"message": "Student created Successfully",
// 	})
// }

func main() {
	log.Printf("HELLO WORLD")
	sqlxConn, err := Connect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed to connect to db: %+v\n", err)
	}
	conn := sqlxConn
	router := gin.Default()
	// router.POST("/student", CreateStudent)
	router.NoRoute(func(c *gin.Context) {
		// In gin this is how you return a JSON response
		c.JSON(404, gin.H{"message": "Not found"})
	})

	router.POST("/student", func(c *gin.Context) {
		boil.DebugMode = true
		var studentReq models.Student

		err := c.BindJSON(&studentReq)
		if err != nil {
			log.Fatal(err)
		}

		student := models.Student{
			ID:        int64(studentReq.ID),
			Firstname: null.StringFrom(studentReq.Firstname.String),
		}

		err = student.Insert(c, conn, boil.Infer())
		if err != nil {
			log.Fatal(err)
		}

		log.Println("student*****************", student)
		fmt.Println("err*****************", err)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Student created Successfully",
		})
	})
	router.GET("/student", func(c *gin.Context) {
		boil.DebugMode = true
		students, err := models.Students().All(c, conn)
		if err != nil {
			log.Println("****", err)
		} // handle err
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "All Students",
			"data":    students,
		})
	})
	router.DELETE("/student/:id", func(c *gin.Context) {
		boil.DebugMode = true
		s := c.Param("id")
		studentId, err1 := strconv.Atoi(s)
		if err1 != nil {
			log.Println("*****", err1)
		}
		student, err := models.FindStudent(c, conn, int64(studentId))
		if err != nil {
			log.Println("*****", err)
		} // handle err

		_, err2 := student.Delete(c, conn)
		if err2 != nil {
			log.Println("****", err2)
		} // handle err
		c.JSON(http.StatusOK, nil)
	})
	router.PUT("/student/:id", func(c *gin.Context) {
		boil.DebugMode = true
		s := c.Param("id")
		studentId, err1 := strconv.Atoi(s)
		if err1 != nil {
			log.Println("*****", err1)
		}
		student, err := models.FindStudent(c, conn, int64(studentId))
		if err != nil {
			log.Println("*****", err)
		} // handle err
		var studentReq models.Student

		err3 := c.BindJSON(&studentReq)
		if err3 != nil {
			log.Fatal(err3)
		}
		student.Firstname = studentReq.Firstname
		_, err2 := student.Update(c, conn, boil.Infer())
		if err2 != nil {
			log.Println("****", err2)
		} // handle err
		c.JSON(http.StatusOK, nil)
	})
	// Init our server
	router.Run(":5000")
}
