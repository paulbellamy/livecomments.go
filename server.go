package main

import (
	"github.com/hoisie/web.go"
  "strconv"
  "json";
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

func main() {
  web.Get("/comments/(.*)", list);
	web.Run("0.0.0.0:3000");
}
