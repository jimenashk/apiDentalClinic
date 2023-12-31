package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"apiDentalClinic/cmd/server/handler"
	"apiDentalClinic/internal/appointments"
	"apiDentalClinic/internal/dentist"
	"apiDentalClinic/internal/patient"
	"apiDentalClinic/pkg/middleware"
	"apiDentalClinic/pkg/store"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("user"),
			os.Getenv("pass"),
			os.Getenv("hostdb"),
			os.Getenv("port"),
			os.Getenv("db_name"))
	)
	fmt.Print(ConnectionString)

	db, err := sql.Open("mysql", ConnectionString)
	if err != nil {
		log.Fatal("Error opening database")
	}
	// Store
	storeSQL := store.NewSQLStore(db)

	// Dentist
	repoDentist := dentist.NewRepository(storeSQL)
	serviceDentist := dentist.NewService(repoDentist)
	handlerDentist := handler.NewDentistHandler(serviceDentist)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	dentists := r.Group("/dentists")
	{
		dentists.GET("", handlerDentist.GetAll())
		dentists.GET(":id", handlerDentist.GetByID())
		dentists.POST("", middleware.AuthenticationMiddleware(), handlerDentist.Post())
		dentists.DELETE(":id", middleware.AuthenticationMiddleware(), handlerDentist.Delete())
		dentists.PUT(":id", middleware.AuthenticationMiddleware(), handlerDentist.Put())
		dentists.PATCH(":id", middleware.AuthenticationMiddleware(), handlerDentist.Patch())

	}
	// Patient
	repoPatient := patient.NewRepository(storeSQL)
	servicePatient := patient.NewService(repoPatient)
	handlerPatient := handler.NewPatientHandler(servicePatient)

	patients := r.Group("/patients")
	{
		patients.GET("", handlerPatient.GetAll())
		patients.GET(":id", handlerPatient.GetByID())
		patients.POST("", middleware.AuthenticationMiddleware(), handlerPatient.Post())
		patients.DELETE(":id", middleware.AuthenticationMiddleware(), handlerPatient.Delete())
		patients.PUT(":id", middleware.AuthenticationMiddleware(), handlerPatient.Put())
		patients.PATCH(":id", middleware.AuthenticationMiddleware(), handlerPatient.Patch())
	}
	// appointments
	repoAppointments := appointments.NewRepository(storeSQL)
	serviceAppointments := appointments.NewService(repoAppointments)
	handlerAppointments := handler.NewAppointmentHandler(serviceAppointments)

	appointments := r.Group("/appointments")
	{
		appointments.GET("", handlerAppointments.GetAll())
		appointments.POST("", middleware.AuthenticationMiddleware(), handlerAppointments.Post())
		appointments.GET(":id", handlerAppointments.GetByID())
		appointments.PUT(":id", middleware.AuthenticationMiddleware(), handlerAppointments.Put())
		appointments.PATCH(":id", middleware.AuthenticationMiddleware(), handlerAppointments.Patch())
		appointments.POST("/post", middleware.AuthenticationMiddleware(), handlerAppointments.PostxLicenseAndDni())
		appointments.DELETE(":id", middleware.AuthenticationMiddleware(), handlerAppointments.Delete())
		appointments.GET("/dni", middleware.AuthenticationMiddleware(), handlerAppointments.GetByDNI())
	}
	r.Run(":8080")
}
