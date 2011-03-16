package main

import (
  "github.com/hoisie/redis.go";
  "time";
  "fmt";
  "json";
  "strconv";
  "bytes";
  "os";
)

type Comment struct {
  Id        int64;
  Author    string;
  Body      string;
  CreatedAt int64;
  PageUrl   string;
}

func New(j []byte) (c Comment, err os.Error) {
  err = json.Unmarshal(j, &c);

  return c, err;
}

func Find(id int64) (c Comment, err os.Error) {
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
    comment, err := Find(id);
    if (err == nil) {
      c = append(c, comment);
    }
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

    c.CreatedAt = time.Seconds();
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
