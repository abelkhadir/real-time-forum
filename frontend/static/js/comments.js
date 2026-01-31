// Comment Management Functions

let currentPostId = null;

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
    </div>
  `).join("");
}

function setCurrentPostId(postId) {
  currentPostId = postId;
}
