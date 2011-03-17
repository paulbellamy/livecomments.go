package main

import (
	"github.com/hoisie/web.go";
  "fmt"
  "strconv";
  "json";
  "bytes";
  "os";
  //"io/ioutil";
)

func allowCrossOriginResourceSharing(ctx *web.Context) {
  ctx.SetHeader("Access-Control-Allow-Origin", "*", true);
}

func list(ctx *web.Context, val string) { 
  allowCrossOriginResourceSharing(ctx);
  url := string(ctx.Params["url"]);

  start, err := strconv.Atoi(ctx.Params["start"]);
  if (err != nil) { start = 0; } // start defaults to 0

  count, err := strconv.Atoi(ctx.Params["count"]);
  if (err != nil) { count = 10; } // Count defaults to 10

  comments := PaginateFor(url, start, count);
  j, _ := json.Marshal(comments);
  ctx.WriteString(string(j));
}

func create(ctx *web.Context, val string) { 
  allowCrossOriginResourceSharing(ctx);
  var comment Comment;
  var err os.Error;
  for k, _ := range ctx.Request.Params {
    comment, err = New(bytes.NewBufferString(k).Bytes());
  }


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
