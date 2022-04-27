package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	// Выполнить подключение к базе данных
	e := Connect()
	if e != nil {
		fmt.Println(e)
		return
	}

	// Создать переменную для обработки запросов
	router := gin.Default()

	// Указать URI "/assets" для отдачи сервером
	// файлов из папки "assets"
	router.Static("/assets", "assets")

	// Указать URI "./image" для сохранения файлов изображений
	// в папку "image"
	router.Use(static.Serve("/image", static.LocalFile("image", false)))

	// Загрузить HTML страницы
	router.LoadHTMLGlob("html/*.html")

	// Обработчики запросов...

	router.GET("/", indexPage)
	router.PUT("/product", createProduct)
	router.GET("/product/:category", selectProduct)

	// Обработчик для сохранения файла на сервере
	router.POST("/upload", upload)

	// Запуск сервера
	_ = router.Run("127.0.0.1:8800")
}

func indexPage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func upload(c *gin.Context) {
	form, e := c.MultipartForm()
	if e != nil {
		fmt.Println(e)
		c.JSON(400, nil)
		return
	}

	files := form.File["MyFiles"]

	for _, file := range files {
		e = c.SaveUploadedFile(file, "image/"+file.Filename)
		if e != nil {
			fmt.Println(e)
		}
	}

	c.JSON(200, nil)
}

type Product struct {
	ID    int
	Name  string
	Price float64
	Image string
}

func createProduct(c *gin.Context) {
	var product Product

	e := c.BindJSON(&product)
	if e != nil {
		fmt.Println(e)
		c.Status(400)
		return
	}

	_, e = Connector.Exec(
		`INSERT INTO product (name, price, image) VALUES ($1, $2, $3)`,
		product.Name,
		product.Price,
		product.Image,
	)
	if e != nil {
		fmt.Println(e)
		c.Status(400)
		return
	}

	c.JSON(200, nil)
}

func selectProduct(c *gin.Context) {
	category := c.Param("category")
	r, e := Connector.Query(`SELECT id, price, name, image FROM product WHERE category=$1 ORDER BY price DESC`, category)
	if e != nil {
		fmt.Println(e)
		c.Status(400)
		return
	}

	products := make([]Product, 0)
	var id int
	var name string
	var price float64
	var image sql.NullString

	for r.Next() {
		e = r.Scan(&id, &price, &name, &image)
		if e != nil {
			fmt.Println(e)
			c.Status(400)
			return
		}

		products = append(products, Product{
			ID:    id,
			Name:  name,
			Price: price,
			Image: image.String,
		})
	}
	c.JSON(200, products)
}
