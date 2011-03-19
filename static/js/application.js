var Comment = Backbone.Model.extend({});

var CommentStore = Backbone.Collection.extend({
  model: Comment,
  url: 'http://appdev.loc:3000/comments'
});
var comments = new CommentStore;


comments.bind('add', function(comment) {
  comments.fetch({success: function(){view.render();}});
});

var commentTemplate = '<div class="comment" id="{{id}}"><div id="commentAuthor">{{author}}</div><div id="commentBody">{{body}}</div></div>';

var CommentView = Backbone.View.extend({
    events: { "submit #commentForm" : "handleNewComment" }

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
