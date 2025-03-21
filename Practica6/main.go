package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"name"`
	Email  string `json:"email"`
}

var usuarios []Usuario

func main() {
	// Gin
	r := gin.Default()

	usuarios = append(usuarios, Usuario{
		ID:     1,
		Nombre: "Alfredo",
		Email:  "Alfredo@mail.com",
	})

	// Rutas
	r.GET("/ping", Ping)
	r.GET("/", Index)
	r.GET("/v1/users", GetUsers)
	r.POST("/v1/users", CreateUser)
	r.PUT("/v1/users/:id", UpdateUser)
	r.DELETE("/v1/users/:id", DeleteUser)

	fmt.Println("Servidor escuchando en el puerto 3000")
	r.Run(":3000")
}


func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Index(c *gin.Context) {
	content, err := os.ReadFile("./public/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "error leyendo el html")
		return
	}
	c.String(http.StatusOK, string(content))
}

func GetUsers(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("test-header", "header")
	c.JSON(http.StatusOK, usuarios)
}

func CreateUser(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el body"})
		return
	}

	var user Usuario
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parseando el JSON"})
		return
	}

	user.ID = len(usuarios) + 1
	usuarios = append(usuarios, user)

	c.Header("Content-Type", "application/json")
	c.Header("test-header", "header")
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el body"})
		return
	}

	var updatedUser Usuario
	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parseando el JSON"})
		return
	}

	for i, user := range usuarios {
		if user.ID == id {
			updatedUser.ID = id
			usuarios[i] = updatedUser
			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
}

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for i, user := range usuarios {
		if user.ID == id {
			usuarios = append(usuarios[:i], usuarios[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
}
