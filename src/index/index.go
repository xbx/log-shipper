package index

import (
  "fmt"
  "time"
  elastigo "github.com/mattbaird/elastigo/lib"

)

var conn *elastigo.Conn
var ElasticHost *string
var ElasticIndex *string

func getConn() *elastigo.Conn {
  if conn == nil {
    conn = elastigo.NewConn()
    if *ElasticHost == "" {
      panic("Undefined Elasticsearch host.")
    }
    conn.Domain = *ElasticHost
  }
  return conn
}

func getToday() string {
  today := time.Now().Format("2006-01-02")
  return today
}

func Index(document map[string]interface{}) {
  conn := getConn()
  name_index := *ElasticIndex + "-" + getToday()
  response, err := conn.Index(name_index, "nginx", "", nil, document)
  if err != nil {
    fmt.Println("Error indexing! " + err.Error())
  }
  fmt.Println(response)
}
