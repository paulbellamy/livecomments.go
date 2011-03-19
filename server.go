package main

import (
	"github.com/hoisie/web.go";
  "fmt"
  "strconv";
  "json";
  "bytes";
  "os";
  "io/ioutil";
  "github.com/madari/go-socket.io";
  "http";
  "log";
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

  k, _ := json.Marshal(ctx);
  fmt.Printf("Context: %s\n", string(k));

  b := bytes.NewBufferString("");
  for k, v := range ctx.Request.Params {
    b.WriteString(fmt.Sprintf("%s=%s&", k, v));
  }
  ctx.Request.Body = b;

  k, _ = ioutil.ReadAll(ctx.Request.Body);
  fmt.Printf("Request Body: %s\n", string(k));

  comment, err = New(k);

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

var sio socketio.SocketIO;
func socketIOConnectHandler(c *socketio.Conn) {
  sio.Broadcast(struct{ announcement string }{"connected: " + c.String()});
}

func socketIODisconnectHandler(c *socketio.Conn) {
  sio.BroadcastExcept(c, struct{ announcement string }{"disconnected: " + c.String()});
}

func socketIOMessageHandler(c *socketio.Conn, msg socketio.Message) {
  sio.BroadcastExcept(c, struct{ message []string }{[]string{c.String(), msg.Data()}});
}

func main() {
  client.Addr = "127.0.0.1:6379";
  web.Get("/comments(.*)", list);
  web.Post("/comments(.*)", create);
	web.Run("0.0.0.0:3000");

  // create the socket.io server and mux it to /socket.io/
  config := socketio.DefaultConfig
  config.Origins = []string{"localhost:8080"}
  sio := socketio.NewSocketIO(&config)
  
  go func() {
    if err := sio.ListenAndServeFlashPolicy(":843"); err != nil {
      log.Println(err)
    }
  }()

  sio.OnConnect(socketIOConnectHandler);
  sio.OnDisconnect(socketIODisconnectHandler);
  sio.OnMessage(socketIOMessageHandler);

  mux := sio.ServeMux();
  mux.Handle("/", http.FileServer("www/", "/"))
  if err := http.ListenAndServe(":8080", mux); err != nil {
    log.Fatal("ListenAndServe: ", err);
  }
}
