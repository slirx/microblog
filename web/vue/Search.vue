<template>
  <div class="search">
    <div class="query">
      <form action="#" v-on:submit.prevent="search">
        <div class="text">
          <input type="text" name="query" id="query" placeholder="search query.." v-model="query"/>
        </div>
        <div class="buttons">
          <input type="submit" :disabled="query.trim() === ''" value="ok">
        </div>
      </form>
    </div>

    <div class="posts">
      <div class="empty" v-if="!posts.length">No posts..</div>
      <div class="item" v-for="post in posts">
        <div class="head">
          <div class="photo">
            <img :src="post.user.photo_url" :alt="post.user.login">
          </div>
          {{ post.user.name }} @
          <router-link :to="'/user/' + post.user.login">{{ post.user.login }}</router-link>
          · {{ post.created_at | moment("from") }}
        </div>
        <div class="text">
          {{ post.text }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  data() {
    return {
      query: '',
      loadedQuery: '',
      posts: [],
      queryId: 0,
      offset: 0,
      loadingTimeout: null,
      isAllPostsFetched: false
    }
  },
  computed: {
    accessToken() {
      return this.$store.state.jwt.access_token;
    },
    currentUser() {
      return this.$store.state.currentUser;
    },
  },
  methods: {
    loadMore() {
      if (this.loadingTimeout != null) {
        clearTimeout(this.loadingTimeout);
      }

      let _this = this;
      this.loadingTimeout = setTimeout(function () {
        _this.search(true);
      }, 200);
    },
    search(isLoadingMore) {
      if (this.currentUser.id === 0) {
        return;
      }

      if (this.query !== this.loadedQuery || isLoadingMore !== true) {
        this.isAllPostsFetched = false;
        this.queryId = 0;
        this.offset = 0;
        this.posts = [];
      }

      if (this.isAllPostsFetched) {
        return;
      }

      let params = new URLSearchParams();
      params.append("offset", this.offset);
      params.append("query_id", this.queryId);
      params.append("query", this.query);

      this.$http.get(
          "/post/search?" + params.toString(),
          this.accessToken,
          (response) => {
            this.loadedQuery = this.query;
            this.queryId = response.data.query_id;

            this.offset += response.data.posts.length;

            if (this.posts.length === response.data.total || response.data.posts.length === 0) {
              this.isAllPostsFetched = true;
              return;
            }

            this.posts.push(...response.data.posts);

            if (this.posts.length === response.data.total) {
              this.isAllPostsFetched = true;
              return;
            }

            if (this.posts.length > 0) {
              this.latestPostId = this.posts[this.posts.length - 1].id;
            }
          }, false);
    }
  },
  created() {
    let _this = this;
    window.onscroll = function () {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight / 100 * 80) {
        _this.loadMore();
      }
    };
  }
}
</script>

<style scoped>

.search {
  margin-top: 1em;
}

.search .query {
  position: relative;
}

.search .query form {
  display: flex;
}

.search .query form .text {
  flex: 9;
}

.search .query form .buttons {
  flex: 1;
}

.search #query {
  width: 100%;
  padding: 0.5em;
  border: 1px solid #e8e8e8;
  margin: 0;
  outline: none;
}

.search .query .buttons input {
  border: none;
  color: #fff;
  background: #343434;
  width: 100%;
  height: 100%;
}

.search .query .buttons input[disabled=disabled] {
  color: #909090;
}

.search .query .buttons input:hover {
  cursor: pointer;
  background: #000000;
}

.search .query .buttons input[disabled=disabled]:hover {
  cursor: default;
  background: #343434;
}
</style>
