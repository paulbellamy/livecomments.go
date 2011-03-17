package main

import (
	"github.com/hoisie/web.go";
  "fmt"
  "strconv";
  "json";
  "bytes";
)

func list(ctx *web.Context, val string) { 
  url := string(ctx.Params["url"]);
  start, _ := strconv.Atoi(ctx.Params["start"]);
  count, _ := strconv.Atoi(ctx.Params["count"]);
  comments := PaginateFor(url, start, count);
  j, _ := json.Marshal(comments);
  ctx.WriteString(string(j));
}

func create(ctx *web.Context, val string) { 
  comment, err := New(bytes.NewBufferString(ctx.Params["body"]).Bytes());

  if (err != nil) {
    ctx.Abort(400, fmt.Sprintf("Error Parsing Comment: %", err));
    return;
  }

  err = comment.Save();
  if (err != nil) {
    ctx.Abort(500, fmt.Sprintf("Error Saving Comment: %", err));
    return;
  }

  j, _ := json.Marshal(comment);
  ctx.WriteString(string(j));
}

func main() {
  client.Addr = "127.0.0.1:6379";
  web.Get("/comments(.*)", list);
  web.Post("/comments(.*)", create);
	web.Run("0.0.0.0:3000");
}
