var humanizeTimestamp = function(timestamp) {
  var now = new Date();
  var then = new Date(timestamp);

  if (now.getFullYear() != then.getFullYear()) {
    return then.toLocaleDateString();
  } else if (now.getMonth() != then.getMonth()) {
    return then.toLocaleDateString();
  } else if (now.getDate() != then.getDate()) {
    return then.toLocaleDateString();
  } else if (now.getHours() != then.getHours()) {
    return then.toLocaleTimeString();
  } else if (now.getMinutes() != then.getMinutes()) {
    var diff = now.getMinutes() - then.getMinutes();
    if (diff == 1) {
      return "a minute ago";
    } else {
      return diff + " minutes ago";
    }
  } else {
    return "just now";
  }
};

var models = this.models = {};

models.Comment = Backbone.Model.extend({});

models.CommentCollection = Backbone.Model.extend({
  model: models.Comment
});

models.LiveCommentsModel = Backbone.Collection.extend({
  initialize: function() {
    this.comments = new models.CommentCollection(); 
  }
});

var commentTemplate = '<div class="comment" id="{{id}}"><div class="commentAuthor">{{author}}:</div><div class="commentTimestamp">{{timestamp}}</div><div class="commentBody">{{body}}</div></div>';

var CommentView = Backbone.View.extend({
  tagName: 'li',

  initialize: function(options) {
    _.bindAll(this, 'render');
    this.model.bind('all', this.render);
  },

  render: function() {
    var data = { id: this.model.get('Id')
               , author: this.model.get('Author')
               , body: this.model.get('Body')
               , timestamp: humanizeTimestamp(this.model.get('CreatedAt')*1000)
               };
    $(this.el).html(Mustache.to_html(commentTemplate, commentData));
    return this;
  }
});

var LiveCommentsView = Backbone.View.extend({
    initialize: function(options) {
      this.model.comments.bind('add', this.addComment);
      this.socket = options.socket;
    }

  , events: { "submit #commentForm" : "postComment" }

  , addComment: function(comment) {
    var view = new CommentView({model: comment});
    var el = view.render().el;
    $('#commentHistory').prepend(el);
    // Animate the new comment
    $(el).slideUp(0);
    $(el).slideDown();
  }

  , msgReceived: function(message){
    console.log("RECEIVED: " + JSON.stringify(message));
    switch(message.event) {
      case 'initial':
        this.model.mport(message.data);
        break;
      case 'comment':
        var newComment = new models.Comment($.parseJSON(message.data));
        this.model.comments.add(newComment);
        break;
    }
  }

  , postComment: function(){
    var author = $('input[name=newCommentAuthor]');
    var body = $('textarea[name=newCommentBody]');
    var newComment= new models.Comment({ Author: author.val()
                                       , Body: body.val()});
    this.socket.send(newComment.toJSON());
    body.val('');
  }
});

var LiveCommentsController = {
  init: function() {
    this.socket = new io.Socket(null);

    this.model = new models.LiveCommentsModel();
    var model = this.model;
    this.view = new LiveCommentsView({model: this.model, socket: this.socket, el: $('#commentArea')});
    var view = this.view;

    this.socket.on('message', function(msg) {view.msgReceived(msg)});
    this.socket.connect();

    this.view.render();

    return this;
  }
};

$(document).ready(function () {
  window.app = LiveCommentsController.init({});
});
