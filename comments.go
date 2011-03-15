package main

import (
  "github.com/hoisie/redis.go";
  "fmt";
  "json";
)

type Comment struct {
  Id        int;
  Author    string;
  Body      string;
  CreatedAt int;
  PageUrl   string;
}

func NewComment(j []byte) (c Comment) {
  err := json.Unmarshal(j, &c);

  if (err != nil) {
    panic(err);
  }
  return c;
}

func FindComment(id int) (c Comment) {
  var client redis.Client;
  js, _ := client.Get( fmt.Sprintf("comment:id:%i", id) );
  return NewComment(js);
}

func (c *Comment) Save() bool {
  return false;
}
