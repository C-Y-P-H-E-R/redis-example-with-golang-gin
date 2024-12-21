package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "home/kushagra/Desktop/Nw_Folder/backend/redisclient"
  "context"
  "fmt"
  "time"
  "strings"
  "encoding/json"
  "github.com/redis/go-redis/v9"
)

// base structure for data
// album represents data about a record album.
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
    {ID: "4", Title: "Tu Aake Dekh le", Artist: "King", Price: 40.00},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    client := redisclient.SetupRedisCaching()
    ctx := context.Background()
    var retrievedAlbum *album

    var url string = c.Request.URL.String()
    pos := strings.Index(url, "id")
    id := string(url[pos+3])

    // key to be checked
    key := "album:" + id

    // get that key
    result, err := client.Get(ctx, key).Result()
    if err == redis.Nil {
      fmt.Println("Key does not exist in Redis")
      // key fetch from other database
      retrievedAlbum = GetDataRemoteServer(id, client)
    } else if err != nil {
      fmt.Println("Error getting data from Redis: ->", err)
    } else {
      // Deserialize the JSON back into a struct
      err = json.Unmarshal([]byte(result), &retrievedAlbum)
      if err != nil {
        fmt.Println("Error deserializing data: ->", err)
      }

      // Print the retrieved album
      fmt.Printf("Retrieved album: %+v\n", retrievedAlbum)
    }
    c.IndentedJSON(http.StatusOK, retrievedAlbum)
}

// fetch data from remote system (assume)
func GetDataRemoteServer(id string, client *redis.Client) *album {
  var targetData *album
  ctx := context.Background()
  time.AfterFunc(1*time.Second, func() {
    // fmt.Println("This message is printed after 3 seconds!")
    for _, i := range albums {
      if i.ID == id {
        targetData = &i
        break;
      }
    }
    // insert that key to redis
    marshalled, err := json.Marshal(targetData)
    if err != nil {
      panic(err)
    }
    key := "album:" + targetData.ID
    err = client.Set(ctx, key, marshalled, 0).Err()
    if err != nil {
      fmt.Println("Error storing data in Redis: ->", err)
    }
  })
  time.Sleep(2 * time.Second)
  return targetData
}


// server setup
func main() {
  router := gin.Default()
  router.GET("/albums", getAlbums)
  router.Run("localhost:6060")
}
