// Post Creation and Display Functions

function expandPostCreationArea() {
    const post = document.getElementById("postCreationArea");
    const title = document.getElementById("titleCreationArea");
    const box = post.closest(".post-box");
    const cats = document.querySelectorAll(".cat");

    // Multi-select logic
    cats.forEach(btn => {
        btn.addEventListener("mousedown", (e) => {
            // Use mousedown to prevent focus loss issues
            e.preventDefault();
            btn.classList.toggle("active");
        });
    });

    post.addEventListener("focus", () => box.classList.add("expanded"));

    // Fix: Check if focus moved inside the same box
    box.addEventListener("focusout", (e) => {
        // Delay check to see where focus went
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

        if (!content) return;

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
            }
        });
    });
}

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
                getPosts()
            } else {
                showToast("red", `Couldnt create post: ${data.error}`);
            }
        })
        .catch(err => console.log(err));
}

function getPosts(page = 1) {
    fetch(`/api/posts?page=${page}`)
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                renderPosts(data.posts);
            } else {
                showToast("red", "couldn't load posts");
            }
        })
        .catch(err => console.log(err));
}

function renderPosts(posts) {
    const container = document.getElementById("posts-container");
    container.innerHTML = "";

    posts.forEach(post => {
        const div = document.createElement("div");
        div.className = "post-card";
        div.onclick = () => openPost(post.ID);

        div.innerHTML = `
        <div class="post-header">
            <div class="avatar">
                <img style="width: 40px" id="avatar" src="/static/images/avatar-white.png">
            </div>
            <div>
            <div class="username">${post.Username}</div>
            <div class="timestamp">${new Date(post.CreatedAt).toLocaleString()}</div>
            </div>
        </div>
        <div class="post-body">${post.Title}</div>
        <div class="post-stats">
            <span class="stats-left">
                ${post.Comments_num} Comments • ${post.Likes_num} Likes
            </span>
            <span class="stats-right">
                ${post.Categories.join(" • ")}
            </span>
        </div>
    `;

        container.appendChild(div);
    });
}

async function openPost(id) {
    try {
        const res = await fetch(`/api/posts/read?id=${id}`);
        const data = await res.json();

        if (!data.success) return showToast("red", "Failed to load post");

        renderPostDetail(data.post);

        document.getElementById("feed-view").classList.add("hidden");
        document.getElementById("post-view").classList.remove("hidden");
    } catch (err) {
        console.log(err)
        showToast("red", "Network error");
    }
}

function closePost() {
    document.getElementById("post-view").classList.add("hidden");
    document.getElementById("feed-view").classList.remove("hidden");
}

function renderPostDetail(post) {
    const container = document.getElementById("post-detail-container");

    container.innerHTML = `
    <div style="display: flex; align-items: center; margin-bottom: 15px;">
      <button class="btn btn-back" onclick="closePost()">←</button>
      <h3>Post Details</h3>
    </div>

    <div class="post-header">
      <div class="avatar"></div>
      <div>
        <div class="username">${post.Username}</div>
        <div class="timestamp">${new Date(post.CreatedAt).toLocaleString()}</div>
      </div>
    </div>

    <h2 class="post-title">${post.Title}</h2>

    <div class="post-body" style="font-size: 18px;">
      ${post.Content}
    </div>

    <div class="post-stats">
      <span>${post.Comments_num} Comments • ${post.Likes_num} Likes</span>
      <span>${post.Categories.join(" • ")}</span>
    </div>

    <div class="comment-section">
      <h4>Comments</h4>
      <br>
      <div class="comment-input-area">
        <div class="avatar" style="width: 30px; height: 30px;"></div>
        <input type="text" id="commentInput" placeholder="Write a comment...">
        <button class="btn btn-primary" style="padding: 0 15px;" onclick="submitComment()">➤</button>
      </div>
      <div id="comments-container" style="margin-top: 20px;"></div>
    </div>
  `;

  // Set the current post ID and load comments
  setCurrentPostId(post.ID);
  loadComments(post.ID);
}

