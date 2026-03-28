
const POSTS_PAGE_SIZE = 10;
let currentPostsPage = 1;
let postsHasMore = false;

function bindPostsPagination() {
    const prev = document.getElementById("posts-prev");
    const next = document.getElementById("posts-next");

    if (prev && prev.dataset.bound !== "true") {
        prev.addEventListener("click", () => {
            if (currentPostsPage > 1) {
                getPosts(currentPostsPage - 1);
            }
        });
        prev.dataset.bound = "true";
    }

    if (next && next.dataset.bound !== "true") {
        next.addEventListener("click", () => {
            if (postsHasMore) {
                getPosts(currentPostsPage + 1);
            }
        });
        next.dataset.bound = "true";
    }

    updatePostsPaginationControls();
}

function updatePostsPaginationControls() {
    const prev = document.getElementById("posts-prev");
    const next = document.getElementById("posts-next");
    const indicator = document.getElementById("posts-page-indicator");

    if (indicator) {
        indicator.textContent = `Page ${currentPostsPage}`;
    }
    if (prev) {
        prev.disabled = currentPostsPage <= 1;
    }
    if (next) {
        next.disabled = !postsHasMore;
    }
}

// expandPostCreationArea manages the post composer interactions.
function expandPostCreationArea() {
    const post = document.getElementById("postCreationArea");
    const title = document.getElementById("titleCreationArea");
    const box = post.closest(".post-box");
    const cats = document.querySelectorAll(".cat");

    cats.forEach(btn => {
        btn.addEventListener("mousedown", (e) => {
            e.preventDefault();
            btn.classList.toggle("active");
        });
    });

    post.addEventListener("focus", () => {
        if (!currentUsername) {
            showToast("red", "login to create post");
            post.blur();
            return;
        }

        box.classList.add("expanded")
    });

    box.addEventListener("focusout", (e) => {
        setTimeout(() => {
            if (!box.contains(document.activeElement)) {
                if (!post.value.trim() && !title.value.trim()) {
                    box.classList.remove("expanded");
                }
            }
        }, 10);
    });

    const submit = () => {
        const content = post.value.trim();
        const titleText = title.value.trim();
        const selectedCats = [...document.querySelectorAll(".cat.active")]
            .map(b => b.dataset.cat);

        if (!titleText || !content) {
            showToast("red", "Title and content are required");
            return;
        }
        if (selectedCats.length === 0) {
            showToast("red", "Select at least one category");
            return;
        }

        createPost(titleText, content, selectedCats);

        post.value = "";
        title.value = "";
        box.classList.remove("expanded");
        cats.forEach(b => b.classList.remove("active"));
    };

    [post, title].forEach(el => {
        el.addEventListener("keydown", (e) => {
            if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                submit();
                post.blur();
            }
        });
    });
}

// createPost submits a new post to the backend.
function createPost(title, content, categories) {
    fetch("/api/posts/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            title: title,
            content: content,
            categories: categories
        }),
    })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showToast("green", "Post Created successfully");
                getPosts(1);
            } else {
                showToast("red", `Couldnt create post: ${data.error}`);
            }
        })
    .catch(err => `Login to create post`);
}

// getPosts fetches the post feed for a page.
function getPosts(page = currentPostsPage) {
    fetch(`/api/posts?page=${page}&limit=${POSTS_PAGE_SIZE}`)
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                if (page > 1 && Array.isArray(data.posts) && data.posts.length === 0) {
                    getPosts(page - 1);
                    return;
                }
                currentPostsPage = page;
                postsHasMore = Array.isArray(data.posts) && data.posts.length === POSTS_PAGE_SIZE;
                renderPosts(data.posts);
                updatePostsPaginationControls();
            } else {
                showToast("red", "couldn't load posts");
            }
        })
        .catch(err => console.log(err));
}

// renderPosts replaces the feed with the provided posts.
function renderPosts(posts) {
    const container = document.getElementById("posts-container");
    container.innerHTML = "";
    if (!posts || posts.length === 0) {
        container.innerHTML = '<div class="no-posts">No posts yet</div>';
        return;
    }

    posts.forEach(post => {
        container.appendChild(buildPostCard(post));
    });
}

// buildPostCard creates one feed card for a post.
function buildPostCard(post) {
    const div = document.createElement("div");
    div.className = "post-card";
    div.onclick = () => openPost(post.ID || post.id);
    const comments = post.Comments_num;
    const cats = post.Categories || [];
    const safeUsername = escapeHTML(post.Username);
    const safeTitle = escapeHTML(post.Title);
    const safeCats = cats.map((cat) => escapeHTML(cat)).join(" • ");

    div.innerHTML = `
        <div class="post-header">
            <div class="avatar">
                <img style="width: 40px" id="avatar" src="/static/images/avatar-white.png">
            </div>
            <div>
            <div class="username">${safeUsername}</div>
            <div class="timestamp">${new Date(post.CreatedAt).toLocaleString()}</div>
            </div>
        </div>
        <div class="post-body">${safeTitle}</div>
        <div class="post-stats">
            <span class="stats-left">
                ${comments} Comments
            </span>
            <span class="stats-right">
                ${safeCats}
            </span>
        </div>
    `;

    return div;
}

// addPostToFeed prepends a new post card to the feed.
function addPostToFeed(post) {
    const container = document.getElementById("posts-container");
    if (!container) return;
    const card = buildPostCard(post);
    container.prepend(card);
}

// openPost loads one post and switches to the detail view.
async function openPost(id) {
    try {
        const res = await fetch(`/api/posts/read?id=${id}`);
        const data = await res.json();

        if (!data.success) return showToast("red", "Failed to load post");

        renderPostDetail(data.post);

        document.getElementById("feed-view").classList.add("hidden");
        document.getElementById("post-view").classList.remove("hidden");
    } catch (err) {
        showToast("red", "Network error");
    }
}

// closePost returns from the detail view to the feed.
function closePost() {
    document.getElementById("post-view").classList.add("hidden");
    document.getElementById("feed-view").classList.remove("hidden");
}

let currentPostId = null;

// setCurrentPostId stores the currently open post ID.
function setCurrentPostId(postId) {
    currentPostId = postId;
}

// renderPostDetail renders the selected post and its comment area.
function renderPostDetail(post) {
    const container = document.getElementById("post-detail-container");
    currentPostId = post.ID || post.id;
    const comments = post.Comments_num || post.comments_count || 0;
    const cats = post.Categories || [];
    const safeUsername = escapeHTML(post.Username);
    const safeTitle = escapeHTML(post.Title);
    const safeContent = escapeHTML(post.Content);
    const safeCats = cats.map((cat) => escapeHTML(cat)).join(" • ");

    container.innerHTML = `
    <div style="display: flex; align-items: center; margin-bottom: 15px;">
      <button class="btn btn-back" onclick="closePost()">←</button>
      <h3>Post Details</h3>
    </div>

    <div class="post-header">
      <div class="avatar"><img src="/static/images/avatar-white.png" style="width:40px"></div>
      <div>
        <div class="username">${safeUsername}</div>
        <div class="timestamp">${new Date(post.CreatedAt).toLocaleString()}</div>
      </div>
    </div>

    <h2 class="post-title">${safeTitle}</h2>

    <br></br>
    <div class="post-body" style="font-size: 18px;">
      ${safeContent}
    </div>

    <div class="post-stats">
      <span>${comments} Comments</span>
      <span>${safeCats}</span>
    </div>

    <div class="comment-section">
      <h4>Comments</h4>
      <br>
      <div class="comment-input-area">
        <div class="avatar" style="width: 30px; height: 30px;"><img src="/static/images/avatar-white.png" style="width:100%"></div>
        <input type="text" id="commentInput" maxlength="500" placeholder="Write a comment..." onkeydown="if(event.key==='Enter'){event.preventDefault();submitComment();}">
        <button class="btn btn-primary" style="padding: 0 15px;" onclick="submitComment()">➤</button>
      </div>
      <div id="comments-container" style="margin-top: 20px;"></div>
    </div>
    `;

    loadComments(currentPostId);
}

// refreshCurrentPost reloads the currently open post.
function refreshCurrentPost() {
    if (!currentPostId) return;

    fetch(`/api/posts/read?id=${currentPostId}`)
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                renderPostDetail(data.post);
            }
        });
}
