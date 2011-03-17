var Comment = Backbone.Model.extend({});

var CommentStore = Backbone.Collection.extend({
  model: Comment,
  url: 'http://appdev.loc:3000/comments'
});
var comments = new CommentStore;


comments.bind('add', function(comment) {
  comments.fetch({success: function(){view.render();}});
});

var commentTemplate = '<div id="comment"><div id="commentAuthor">{{author}}</div><div id="commentBody">{{body}}</div></div>';

var CommentView = Backbone.View.extend({
    events: { "submit #commentForm" : "handleNewComment" }

  , handleNewComment: function(data) {
      var inputField = $('input[name=newCommentBody]');
      comments.create({Body: inputField.val()});
      inputField.val('');
    }

  , render: function() {
      var data = comments.map(function(comment) { return {body: comment.get('Body')}; });
      var result = data.reduce(function(memo,commentData) { return memo + Mustache.to_html(commentTemplate, commentData); }, '');
      $("#commentHistory").html(result);
      return this;
    }
});

var view = new CommentView({el: $('#commentArea')});

setInterval(function(){
  comments.fetch({success: function(){view.render();}});
},10000);
