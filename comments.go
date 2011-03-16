package main

import (
  "github.com/hoisie/redis.go";
  "fmt";
  "json";
  "strconv";
  "bytes";
)

type Comment struct {
  Id        int64;
  Author    string;
  Body      string;
  CreatedAt int64;
  PageUrl   string;
}

func New(j []byte) (c Comment) {
  err := json.Unmarshal(j, &c);

  if (err != nil) {
    panic(err);
  }
  return c;
}

func Find(id int64) (c Comment) {
  var client redis.Client;
  c.Id = id;
  js, _ := client.Get( fmt.Sprintf("comment:id:%i", id) );
  return New(js);
}

func PaginateFor(url string, start int, count int) (c []Comment) {
  var client redis.Client;
  commentIds, _ := client.Lrange(fmt.Sprintf("comment:page_url:%s", url), start, count);
  for _, idString := range commentIds {
    id, _ := strconv.Atoi64(string(idString));
    c = append(c, Find(id));
  }
  return c;
}

func (c *Comment) Save() bool {
  var client redis.Client;

  newRecord := false;
  if (c.Id == -1) { newRecord = true; }

  if (newRecord) {
    // New record we should get an Id for it
    id, err := client.Incr("global:nextCommentId");
    if (err != nil) { return false; }
    c.Id = id;
  }

  // Store it by the primary key
  j, err := json.Marshal(c);
  client.Set(fmt.Sprintf("comment:id:%i", c.Id), j);
  if (err != nil) { return false; }

  if (newRecord) {
    // New record we should insert it into the page listing
    err :=client.Lpush(fmt.Sprintf("comment:page_url:%s", c.PageUrl), bytes.NewBufferString(strconv.Itoa64(c.Id)).Bytes());
    if (err != nil) { return false; }
  }

  return true;
}
