package main

import (
      elastic "gopkg.in/olivere/elastic.v3"

      "encoding/json"
      "fmt"
      "net/http"
      "reflect"
      "regexp"
      "time"

      "github.com/dgrijalva/jwt-go"
)

const (
      TYPE_USER = "user"
)

var (
      usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

type User struct {
      Username string `json:"username"`
      Password string `json:"password"`
      Age int `json:”age”`
      Gender string `json:”gender”`
}

// checkUser checks whether user is valid
func checkUser(username, password string) bool {
      es_client, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
      if err != nil {
             fmt.Printf("ES is not setup %v\n", err)
             return false
      }

      // Search with a term query
      termQuery := elastic.NewTermQuery("username", username)
      queryResult, err := es_client.Search().
             Index(INDEX).
             Query(termQuery).
             Pretty(true).
             Do()
      if err != nil {
             fmt.Printf("ES query failed %v\n", err)
             return false
      }

      var tyu User
      for _, item := range queryResult.Each(reflect.TypeOf(tyu)) {
             u := item.(User)
             return u.Password == password && u.Username == username
      }
      // If no user exist, return false.
      return false
}
// add user adds a new user
