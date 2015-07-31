package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
	"log"
	"strconv"
)

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "user=postgres dbname=dataitem sslmode=disable password=password")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(DataItem{}, "DataItem").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}



type DataItem struct {
	Id   int64  `db:"id" json:"id"`
	Uuid string `db:"uuid" json:"uuid"`
	Name string `db:"name" json:"name"`
	Data string `db:"data" json:"data"`
}

func main() {
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.GET("/data", GetDataItems)
		v1.GET("/data/:uuid", GetDataItem)
		v1.POST("/data", PostDataItem)
		v1.PUT("/data/:uuid", UpdateDataItem)
		v1.DELETE("/data/:uuid", DeleteDataItem)
	}
	r.Run(":8080")
}

func GetDataItems(c *gin.Context) {
	var dataItems []DataItem
	_, err := dbmap.Select(&dataItems, "SELECT * FROM dataitem")
	if err == nil {
		c.JSON(200, dataItems)
	} else {
		c.JSON(404, gin.H{"error": "no dataitem(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/data
}

func GetDataItem(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
	if uuid == "uuid1" {
		content := gin.H{
			"id":   1,
			"uuid": "uuid1",
			"name": "name1",
			"data": `
			"key" : "value"
		`}
		c.JSON(200, content)
	} else if uuid == "uuid2" {
		content := gin.H{
			"id":   2,
			"uuid": "uuid2",
			"name": "name2",
			"data": `
			"key" : "value"
		`}
		c.JSON(200, content)
	} else {
		content := gin.H{"error": "dataItem with uuid#" + uuid + " not found"}
		c.JSON(404, content)
	}
	// curl -i http://localhost:8080/api/v1/users/1
}
func PostDataItem(c *gin.Context) {
	var dataItem DataItem
	fmt.Println("binding")
	c.Bind(&dataItem)
	fmt.Println("bound")
	if dataItem.Name != "" {
		fmt.Println("inserting")
		fmt.Printf("%+v\n", dataItem)
		var id int64
		err := dbmap.Db.QueryRow(`INSERT INTO dataItem (uuid, name, data) VALUES ($1, $2, $3) returning id;`, dataItem.Uuid, dataItem.Name, dataItem.Data).Scan(&id)
		fmt.Println(err)
		//if insert != "" {
		fmt.Println("inserted")
		//id, err := insert.LastInsertId()
		if err == nil {
			fmt.Println("returning")
			content := &DataItem{
				Id:   id,
				Name: dataItem.Name,
				Uuid: dataItem.Uuid,
				Data: dataItem.Data,
			}
			c.JSON(201, content)
		} else {
			checkErr(err, "Insert failed")
		}
		//}
	} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users

}


func UpdateDataItem(c *gin.Context) {
	id := c.Params.ByName("id")
	var dataItem DataItem
	err := dbmap.SelectOne(&dataItem, "SELECT * FROM dataitem WHERE id=?", id)
	if err == nil {
		var json DataItem
		c.Bind(&json)
		dataItem_id, _ := strconv.ParseInt(id, 0, 64)
		dataItem := DataItem{
			Id: dataItem_id,
			Name: dataItem.Name,
			Uuid: dataItem.Uuid,
			Data: dataItem.Data,
		}
		if dataItem.Name != "" && dataItem.Uuid != "" {
			_, err = dbmap.Update(&dataItem)
			if err == nil {
				c.JSON(200, dataItem)
			} else {
				checkErr(err, "Updated failed")
			}
		} else {
			c.JSON(422, gin.H{"error": "fields are empty"})
		}
		} else {
		c.JSON(404, gin.H{"error": "user not found"})
		}
		// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
	}
func DeleteDataItem(c *gin.Context) {
	// The futur code…
}
