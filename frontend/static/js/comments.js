// Comment Management Functions

function submitComment() {
  const input = document.getElementById("commentInput");
  const commentText = input.value.trim();

  if (!commentText) {
    showToast("red", "Comment cannot be empty");
    return;
  }

  if (!currentPostId) {
    showToast("red", "No post selected");
    return;
  }

  fetch("/api/comments/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      post_id: currentPostId,
      content: commentText
    }),
  })
    .then(res => res.json())
    .then(data => {
      if (data.success) {
        showToast("green", "Comment added successfully");
        input.value = "";
        loadComments(currentPostId);
        refreshCurrentPost();
      } else {
        showToast("red", `Couldn't add comment: ${data.error}`);
      }
    })
    .catch(err => {
      console.log(err);
      showToast("red", "Network error");
    });
}

function loadComments(postId) {
  fetch(`/api/comments?post_id=${postId}`)
    .then(res => res.json())
    .then(data => {
      if (data.success) {
        renderComments(data.comments || []);
      } else {
        showToast("red", "Couldn't load comments");
      }
    })
    .catch(err => {
      console.log(err);
      showToast("red", "Network error");
    });
}

function renderComments(comments) {
  const container = document.getElementById("comments-container");
  if (!container) return;

  if (!comments || comments.length === 0) {
    container.innerHTML = "<p style='color: #999;'>No comments yet</p>";
    return;
  }

  container.innerHTML = comments.map(comment => `
    <div style="border-left: 2px solid #ccc; padding-left: 10px; margin-bottom: 10px;">
      <div style="font-weight: bold; font-size: 14px;">${comment.username}</div>
      <div style="font-size: 12px; color: #666;">${new Date(comment.created_at).toLocaleString()}</div>
      <div style="margin-top: 5px;">${comment.content}</div>
      <div style="font-size: 12px; color: #999; margin-top: 5px;">👍 ${comment.likes_count} • 👎 ${comment.dislikes_count}</div>
    </div>
  `).join("");
}
function likeComment(commentId) {
  voteComment(commentId, true);
}

function dislikeComment(commentId) {
  voteComment(commentId, false);
}

function voteComment(commentId, isLike) {
  fetch("/api/comments/like", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      comment_id: commentId,
      is_like: isLike
    }),
  })
  .then(res => res.json())
  .then(data => {
    if (data.success) {
      loadComments(currentPostId);
    }
  });
}
function renderComments(comments) {
  const container = document.getElementById("comments-container");
  if (!container) return;

  if (!comments || comments.length === 0) {
    container.innerHTML = "<p style='color: #999; font-style: italic;'>No comments yet. Be the first!</p>";
    return;
  }

  container.innerHTML = comments.map(comment => `
    <div style="border-left: 3px solid #515253; padding: 10px; margin-bottom: 15px; background-color: #f9f9f9; border-radius: 0 5px 5px 0;">
      
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 5px;">
        <span style="font-weight: bold; font-size: 14px; color: #333;">${comment.username}</span>
        <span style="font-size: 11px; color: #888;">${new Date(comment.created_at).toLocaleString()}</span>
      </div>

      <div style="margin-bottom: 10px; font-size: 14px; line-height: 1.4;">${comment.content}</div>

      <div style="display: flex; gap: 10px; align-items: center;">
        
        <button 
          onclick="likeComment(${comment.id})" 
          style="cursor: pointer; background: white; border: 1px solid #ccc; border-radius: 4px; padding: 4px 10px; font-size: 12px; display: flex; align-items: center; gap: 5px; transition: background 0.2s;"
          onmouseover="this.style.background='#717070'" 
          onmouseout="this.style.background='white'">
          👍 <span>${comment.likes_count || 0}</span>
        </button>

        <button 
          onclick="dislikeComment(${comment.id})" 
          style="cursor: pointer; background: white; border: 1px solid #ccc; border-radius: 4px; padding: 4px 10px; font-size: 12px; display: flex; align-items: center; gap: 5px; transition: background 0.2s;"
          onmouseover="this.style.background='#717070'" 
          onmouseout="this.style.background='white'">
          👎 <span>${comment.dislikes_count || 0}</span>
        </button>

      </div>
    </div>
  `).join("");
}