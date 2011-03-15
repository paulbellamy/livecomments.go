package main

import (
 // "github.com/hoisie/redis.go";
 // "strings";
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

func (c *Comment) Save() bool {
  return false;
}
