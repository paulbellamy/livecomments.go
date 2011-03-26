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

var Comment = Backbone.Model.extend({});

var CommentStore = Backbone.Collection.extend({
  model: Comment,
  url: 'http://appdev.loc:3000/comments'
});
var comments = new CommentStore;


comments.bind('add', function(comment) {
  comments.fetch({success: function(){
    view.render();
    // Animate the new comment
    $('.comment#' + comment.get('Id')).slideUp(0);
    $('.comment#' + comment.get('Id')).slideDown();
  }});
});

var commentTemplate = '<div class="comment" id="{{id}}"><div class="commentAuthor">{{author}}:</div><div class="commentTimestamp">{{timestamp}}</div><div class="commentBody">{{body}}</div></div>';

var CommentView = Backbone.View.extend({
    initialize: function(options) {
      this.model.comments.bind('add', this.addComment);
      this.socket = options.socket;
    }

  , events: { "submit #commentForm" : "handleNewComment" }

  , addComment: function(comment) {
    var view = new CommentView({model: comment});
    $('#commentHistory').prepend(view.render().el);
  }

  , handleNewComment: function(data) {
      var author = $('input[name=newCommentAuthor]');
      var body = $('textarea[name=newCommentBody]');
      comments.create({ Author: author.val()
                      , Body: body.val()});
      body.val('');
    }

  , render: function() {
      var data = comments.map(function(comment) { return { id: comment.get('Id')
                                                         , author: comment.get('Author')
                                                         , body: comment.get('Body')
                                                         , timestamp: humanizeTimestamp(comment.get('CreatedAt')*1000)
                                                         };
                                                });
      var result = data.reduce(function(memo,commentData) { return memo + Mustache.to_html(commentTemplate, commentData); }, '');
      $("#commentHistory").html(result);
      return this;
    }
});

var view = new CommentView({el: $('#commentArea')});

comments.fetch({success: function(){view.render();}});
setInterval(function(){
  comments.fetch({success: function(){view.render();}});
},10000);
