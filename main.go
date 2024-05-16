package main

import (
    "database/sql"
    "log"
    "net/http"
	"fmt"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbName     = "demo"
	dbUser     = "postgres"
	dbPassword = "1234"

)

type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}


var Dbase *sql.DB



func getAlbums(c *gin.Context) {
    var albums []album
    rows, err := Dbase.Query("SELECT id, title, artist, price FROM albums")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var album album
        if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
            log.Fatal(err)
        }
        albums = append(albums, album)
    }
    c.IndentedJSON(http.StatusOK, albums)
}



func postAlbums(c *gin.Context) {
    var new_album album
    if err := c.BindJSON(&new_album); err != nil {
        return
    }
    _, err := Dbase.Exec("INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3)",
        new_album.Title, new_album.Artist, new_album.Price)
    if err != nil {
        log.Fatal(err)
    }
    c.IndentedJSON(http.StatusCreated, new_album)
}



func initDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
    Dbase, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    if err = Dbase.Ping(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Database connected!")
}




func deleteAlbums(c *gin.Context) {
    id := c.Param("id")
    if _, err := Dbase.Exec("DELETE FROM albums WHERE id = $1", id); err != nil {
        log.Fatal(err)
    }
    c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted successfully"})
}



func updateAlbums(c *gin.Context) {
    id := c.Param("id")
    var AlbumUpdate album
    if err := c.BindJSON(&AlbumUpdate); err != nil {
        return
    }

    if _, err := Dbase.Exec("UPDATE albums SET title = $1, artist = $2, price = $3 WHERE id = $4",
        AlbumUpdate.Title, AlbumUpdate.Artist, AlbumUpdate.Price, id); err != nil {
        log.Fatal(err)
    }
    c.IndentedJSON(http.StatusOK, gin.H{"message": "album updated successfully"})
}

func main() {
    router := gin.Default()
    initDB()

    router.GET("/albums", getAlbums)
    router.POST("/albums", postAlbums)
    router.DELETE("/albums/:id", deleteAlbums)
    router.PUT("/albums/:id", updateAlbums)

    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
