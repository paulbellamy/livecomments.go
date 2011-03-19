package main

import (
	"github.com/hoisie/web.go";
  "fmt"
  "strconv";
  "json";
  "bytes";
  "os";
  //"io/ioutil";
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

  if start, err := strconv.Atoi(ctx.Params["start"]); err != nil {
    start = 0; // start defaults to 0
  }

  if count, err := strconv.Atoi(ctx.Params["count"]) err != nil {
    count = 10; // Count defaults to 10
   }

  comments := PaginateFor(url, start, count);
  j, _ := json.Marshal(comments);
  ctx.WriteString(string(j));
}

func create(ctx *web.Context, val string) { 
  allowCrossOriginResourceSharing(ctx);
  var comment Comment;
  var err os.Error;

  k, _ := json.Marshal(ctx.ParamData);
  fmt.Printf("ParamData: %s\n", string(k));

  if comment, err = New(k); err != nil {
    ctx.Abort(400, fmt.Sprintf("Error Parsing Comment: %", err));
    return;
  }

  if err = comment.Save(); err != nil {
    ctx.Abort(500, fmt.Sprintf("Error Saving Comment: %", err));
    return;
  }

  j, _ := json.Marshal(comment);
  ctx.WriteString(string(j));
}

var sio socketio.SocketIO;
func socketIOConnectHandler(c *socketio.Conn) {
  //sio.Broadcast(struct{ announcement string }{"connected: " + c.String()});
}

func socketIODisconnectHandler(c *socketio.Conn) {
  //sio.BroadcastExcept(c, struct{ announcement string }{"disconnected: " + c.String()});
}

func socketIOMessageHandler(c *socketio.Conn, msg socketio.Message) {
  if comment, err := New(bytes.NewBufferString(msg.Data()).Bytes()); err != nil {
    log.Println(fmt.Sprintf("Error Parsing Comment %s", err));
    c.Send(fmt.Sprintf("Error Parsing Comment %s", err));
    return;
  }

  if err = comment.Save(); err != nil {
    log.Println(fmt.Sprintf("Error Saving Comment %s", err));
    c.Send(fmt.Sprintf("Error Saving Comment %s", err));
    return;
  }

  if j, err := json.Marshal(comment); err != nil {
    log.Println(fmt.Sprintf("Error Saving Comment %s", err));
    c.Send(fmt.Sprintf("Error Saving Comment %s", err));
    return;
  }

  sio.Broadcast(struct{ announcement string }{string(j)});
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
