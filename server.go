package main

import (
	"github.com/hoisie/web.go";
  "strconv";
  "json";
  "bytes";
)

func list(ctx *web.Context, val string) { 
  url := string(ctx.Params["url"]);
  if (url != "") {
    start, _ := strconv.Atoi(ctx.Params["start"]);
    count, _ := strconv.Atoi(ctx.Params["count"]);
    comments := PaginateFor(url, start, count);
    j, _ := json.Marshal(comments);
    ctx.WriteString(string(j));
  }
}

func create(ctx *web.Context, val string) { 
  comment, err := New(bytes.NewBufferString(ctx.Params["body"]).Bytes());

  if (err != nil) {
    ctx.Abort(400, "Error Parsing Comment");
    return;
  }

  if (!comment.Save()) {
    ctx.Abort(500, "Error Saving Comment");
    return;
  }

  j, _ := json.Marshal(comment);
  ctx.WriteString(string(j));
}

func main() {
  web.Get("/comments(.*)", list);
  web.Post("/comments(.*)", create);
	web.Run("0.0.0.0:3000");
}
