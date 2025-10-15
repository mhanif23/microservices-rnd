package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	MembershipStatus string `json:"membership_status"`
}

func main() {
	r := gin.Default()
	port := getenv("PORT", "3000")
	dsn := getenv("DATABASE_URL", "postgres://library:library@localhost:5432/librarydb?sslmode=disable")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil { log.Fatal(err) }
	defer pool.Close()

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status":"ok"}) })

	// create
	r.POST("/users", func(c *gin.Context) {
		var in struct{ Email, Name string }
		if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, errorJSON("BAD_REQUEST", err.Error())); return }
		id := uuid.New().String()
		_, err := pool.Exec(ctx, `INSERT INTO users(id,email,name,membership_status) VALUES($1,$2,$3,'ACTIVE')`, id, in.Email, in.Name)
		if err != nil { c.JSON(http.StatusBadRequest, errorJSON("BAD_REQUEST", err.Error())); return }
		c.JSON(http.StatusCreated, gin.H{"id": id, "email": in.Email, "name": in.Name, "membership_status": "ACTIVE"})
	})

	// get by id
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var u User
		err := pool.QueryRow(ctx, `SELECT id::text,email,name,membership_status FROM users WHERE id=$1`, id).
			Scan(&u.ID,&u.Email,&u.Name,&u.MembershipStatus)
		if err != nil { c.JSON(http.StatusNotFound, errorJSON("NOT_FOUND","not found")); return }
		c.JSON(http.StatusOK, u)
	})

	// eligibility
	r.GET("/users/:id/eligibility", func(c *gin.Context) {
		id := c.Param("id")
		var status string
		err := pool.QueryRow(ctx, `SELECT membership_status FROM users WHERE id=$1`, id).Scan(&status)
		if err != nil { c.JSON(http.StatusNotFound, errorJSON("NOT_FOUND","not found")); return }
		c.JSON(http.StatusOK, gin.H{ "user_id": id, "can_borrow": status == "ACTIVE", "max_books": 3 })
	})

	log.Printf("user-svc running on :%s", port)
	r.Run(":" + port)
}

func errorJSON(code, msg string) gin.H { return gin.H{"code":code,"message":msg} }
func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }
