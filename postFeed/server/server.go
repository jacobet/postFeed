package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
	"postFeed/structs"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB

func main() {

	var fileName = "test.db"

	bNewFile, err := openDB(fileName)
	if err == nil && bNewFile {
		err = db.AutoMigrate(structs.Post{}, structs.Comment{})
	}
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to postFeed!")
	})

	e.POST("/posts", createPost)
	e.GET("/posts", getPost)
	e.GET("/posts/:id", getSinglePost)
	e.PUT("/posts/:id", updatePost)
	e.DELETE("/posts/:id", deletePost)

	e.POST("/comments", createComment)

	e.Logger.Fatal(e.Start(":1323"))
}
func openDB(fileName string) (bNewFile bool, err error) {
	bNewFile = false
	if _, err = os.Stat(fileName); err == nil {
		db, err = gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(fileName), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			bNewFile = true
		}
	}
	return bNewFile, err
}

func createComment(c echo.Context) error {
	comment := new(structs.Comment)
	if err := c.Bind(comment); err != nil {
		return c.JSON(http.StatusBadRequest, comment)
	}
	post := new(structs.Post)
	db.First(&post, comment.PostID)
	if post.ID == 0 {
		return c.JSON(http.StatusForbidden, comment.PostID)
	}
	comment.CreatedAt = time.Now()
	db.Create(&comment)
	return c.JSON(http.StatusCreated, comment)
}

func createPost(c echo.Context) error {
	post := new(structs.Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}
	post.CreatedAt = time.Now()
	db.Create(&post)
	return c.JSON(http.StatusCreated, post)
}

func checkParams(post *structs.Post) (model *gorm.DB) {
	model = db.Model(&structs.Post{})
	if post.Author != "" {
		model = model.Where("lower(author) = ?", strings.ToLower(post.Author))
	}
	if post.Content != "" {
		model = model.Where("lower(content) = ?", strings.ToLower(post.Content))
	}
	if post.Like > 0 {
		model = model.Where("like = ?", post.Like)
	}
	if post.Dislike > 0 {
		model = model.Where("dislike = ?", post.Dislike)
	}
	temp := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	if post.CreatedAt.After(temp) {
		model = model.Where("create_at = ?", post.CreatedAt)
	}
	if post.UpdatedAt.After(temp) {
		model = model.Where("update_at = ?", post.UpdatedAt)
	}
	return model
}

type paginate struct {
	skip  int
	limit int
	page  int
}

func (p *paginate) paginatedResult(db *gorm.DB) *gorm.DB {
	var limit int
	offset := p.skip * p.page
	if p.limit == -1 {
		limit = p.page
	} else {
		limit = p.limit
	}
	return db.Offset(offset).Limit(limit)
}

func getPost(c echo.Context) error {

	jsonMap := make(map[string]interface{})
	json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	var paginate paginate
	var page int
	var skip int
	var limit int
	var order string

	page = 5
	skip = 0
	limit = -1

	if jsonMap["skip"] != nil {
		skip, _ = strconv.Atoi(fmt.Sprintf("%v", jsonMap["skip"]))
	}
	paginate.skip = skip

	if jsonMap["page"] != nil {
		page, _ = strconv.Atoi(fmt.Sprintf("%v", jsonMap["page"]))
	}
	paginate.page = page

	if jsonMap["limit"] != nil {
		limit, _ = strconv.Atoi(fmt.Sprintf("%v", jsonMap["limit"]))
	}
	paginate.limit = limit
	if jsonMap["order"] != nil {
		order = fmt.Sprintf("%v", jsonMap["order"])
	}

	post := new(structs.Post)
	if jsonMap["author"] != nil {
		post.Author = fmt.Sprintf("%v", jsonMap["author"])
	}
	if jsonMap["content"] != nil {
		post.Content = fmt.Sprintf("%v", jsonMap["content"])
	}
	if jsonMap["create-at"] != nil {
		date, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", jsonMap["create-at"]))
		post.CreatedAt = date
	}

	var model = checkParams(post)
	var list []structs.Post
	model.Scopes(paginate.paginatedResult).Order(order).Scan(&list)

	return c.JSON(http.StatusOK, list)
}

func getSinglePost(c echo.Context) error {
	post := new(structs.Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}

	db.First(&post, id)
	if post.ID == 0 {
		return c.JSON(http.StatusNotFound, id)
	}
	return c.JSON(http.StatusOK, post)
}

func deletePost(c echo.Context) error {
	post := new(structs.Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}

	var id, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}

	db.First(&post, id)
	if post.ID == 0 {
		return c.JSON(http.StatusNotFound, id)
	}
	if err := db.Delete(&post, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, id)
	}
	return c.NoContent(http.StatusNoContent)
}

func updatePost(c echo.Context) error {
	post := new(structs.Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}
	var id, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, post)
	}

	var newPost structs.Post
	db.First(&newPost, id)
	if newPost.ID == 0 {
		return c.JSON(http.StatusNotFound, id)
	}
	newPost.UpdatedAt = time.Now()
	if post.Author != "" {
		newPost.Author = post.Author
	}
	if post.Content != "" {
		newPost.Content = post.Content
	}
	if post.Like != 0 {
		newPost.Like++
	}
	if post.Dislike != 0 {
		newPost.Dislike++
	}

	db.Model(&structs.Post{}).Where("id = ?", id).Updates(&newPost)
	return c.JSON(http.StatusOK, newPost)

}
