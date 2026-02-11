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
      showToast("red", "Login to comment");
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
    </div>
  `).join("");
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
        

      </div>
    </div>
  `).join("");
}