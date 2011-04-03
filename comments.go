package main

import (
  "time";
  "fmt";
  "json";
  "strconv";
  "bytes";
  "os";
  "log";
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

func Create(j []byte) (c Comment, err os.Error) {
  c, err = New(j);
  if err != nil {
    log.Println(fmt.Sprintf("Error Parsing Comment: %s", err));
    return c, err;
  }

  if err = c.Save(); err != nil {
    log.Println(fmt.Sprintf("Error Saving Comment: %s", err));
    return c, err;
  }

  return c, nil;
}

func Find(id int64) (c Comment, err os.Error) {
  c.Id = id;
  js, _ := client.Get( fmt.Sprintf("comment:id:%i", id) );
  return New(js);
}

func PaginateFor(url string, start int, count int) (c []Comment) {
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

func (c *Comment) ToJson() (j string) {
  if j, err := json.Marshal(c); err != nil {
    log.Println("Error Encoding Comment to Json: ", err);
  } else {
    return string(j);
  }
  return "";
}

func (c *Comment) Save() (err os.Error) {

  newRecord := false;
  if (c.Id == 0) { newRecord = true; }

  if (newRecord) {
    // New record we should get an Id for it
    id, err := client.Incr("global:nextCommentId");
    if (err != nil) { return err; }
    c.Id = id;

    c.CreatedAt = time.Seconds();
  }

  // Store it by the primary key
  client.Set(fmt.Sprintf("comment:id:%i", c.Id), []uint8(c.ToJson()));
  if (err != nil) { return err; }

  if (newRecord) {
    // New record we should insert it into the page listing
    err :=client.Lpush(fmt.Sprintf("comment:page_url:%s", c.PageUrl), bytes.NewBufferString(strconv.Itoa64(c.Id)).Bytes());
    if (err != nil) { return err; }
  }

  return nil;
}
