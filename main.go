package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	_ "github.com/BeHappych/lastPlease/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type List struct {
	Id        string `json:"id"`
	Full_name string `json:"full_name"`
	Birthday  string `json:"birthday"`
	Address   string `json:"address"`
}

var database *sql.DB

// @Summary          Show all
// @Description      Get all
// @Tags             lists
// @Accept           json
// @Produce          json
// @Success          200
// @Router           /lists [get]
func getLists(c *gin.Context) {

	rows, err := database.Query("select * from Lists")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	lists := []List{}

	for rows.Next() {
		p := List{}
		err := rows.Scan(&p.Id, &p.Full_name, &p.Birthday, &p.Address)
		if err != nil {
			fmt.Println(err)
			continue
		}
		lists = append(lists, p)
	}

	c.IndentedJSON(http.StatusOK, lists)
}

// @Summary          New list
// @Description      Add new list
// @Tags             lists
// @Accept           json
// @Produce          json
// @Success          200
// @Param            newList body List true "Birthday input format - 1999-01-02"
// @Router           /lists [post]
func postList(c *gin.Context) {
	var newList List

	if err := c.BindJSON(&newList); err != nil {
		return
	}
	database.Exec("insert into Lists (Full_name, Birthday, Address) values ($1, $2, $3)", newList.Full_name, newList.Birthday, newList.Address)

	c.IndentedJSON(http.StatusCreated, newList)
}

// @Summary          Delete by Id
// @Description      Delete list
// @Tags             lists
// @Accept           json
// @Produce          json
// @Param            id path string true "List ID"
// @Success          204
// @Router           /lists/{id} [delete]
func deleteById(c *gin.Context) {
	id := c.Param("id")
	database.Exec("delete from Lists where id = $1", id)
	c.IndentedJSON(http.StatusNotFound, id)
}

// @Summary          Update by Id
// @Description      Update list
// @Tags             lists
// @Accept           json
// @Produce          json
// @Param            id path string true "List ID"
// @Param            updateList body List true "Birthday input format - 1999-01-02"
// @Success          201 {array} List
// @Router           /lists/{id} [PUT]
func updateById(c *gin.Context) {
	id := c.Param("id")
	var updateList List

	if err := c.BindJSON(&updateList); err != nil {
		return
	}

	database.Exec("update Lists set Full_name = $1, Birthday = $2, Address = $3 where id = $4", updateList.Full_name, updateList.Birthday, updateList.Address, id)
	c.IndentedJSON(http.StatusCreated, id)

}

// @title           Swagger API
// @version         1.0
// @description     Swagger API for Golang Project.
// @host            localhost:8080

func main() {

	db, err := sql.Open("postgres", "user=user password=pass dbname=betadb sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/lists", getLists)
	router.POST("/lists", postList)
	router.DELETE("/lists/:id", deleteById)
	router.PUT("/lists/:id", updateById)

	router.Run("localhost:8080")

}
