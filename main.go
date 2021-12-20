package main

import (
    "net/http"
    "database/sql"
    "log"
    _"github.com/mattn/go-sqlite3"
    "github.com/gin-gonic/gin"
)


//Tract made to contain tract name, artist and album
type Tract struct {
    Name string `json:"name"`
    Artist string `json:"artist"`
    Album string `json:"album"`
}

//quries tracks in database by title and sends title, artist and album JSON to client
func getTractsByName(c *gin.Context) {

    var tracts []Tract
    var db, err = sql.Open("sqlite3", "./Chinook_Sqlite.sqlite")
    name := c.Param("name")
    log.Println("User searching for "+name)

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "issue connecting to database"})
        log.Println("Issue connecting to database")
        return
    }

    //confirming that connecting to the database works
    pingErr := db.Ping()
    if pingErr != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "issue pinging database"})
        log.Println("Issue pinging database")
        return
    }

    //SQL Search query: joining the album and artist table to return track name album title and artist
    rows, err := db.Query("SELECT t.name, ar.name, al.title FROM track AS t, album AS al, artist AS ar WHERE t.name LIKE '%" + name + "%' AND t.AlbumID = al.AlbumID AND al.Artistid = ar.Artistid")
    if err != nil{
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "issue querying database"})
        log.Println("Issue querying database")
        return
    }
    defer rows.Close()

    
    //Loop through all returned rows and scan sql output to Tract struct to add to track slice being returned
    for rows.Next(){
        var t Tract
        if err := rows.Scan(&t.Name, &t.Artist, &t.Album); err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "issue reading from database output"})
            return
        } 
        tracts = append(tracts, t)
    }
    
    if tracts == nil{
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "track not found"})
        log.Println("No trackt found")
        return
    }

    log.Println("Tracks found and returned successfully")
    c.IndentedJSON(http.StatusOK, tracts)
}

func main() {

    router := gin.Default()
    router.GET("/tracts/:name", getTractsByName)

    router.Run("localhost:4041")
}
