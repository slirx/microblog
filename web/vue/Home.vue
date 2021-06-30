<template>
  <div class="home">
    <new-post/>
    <div class="feed">
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
  </div>
</template>

<script>
import NewPost from './NewPost.vue';
import eventBus from "../js/event-bus";

export default {
  components: {
    "new-post": NewPost,
  },
  data() {
    return {
      posts: [],
      latestPostId: 0,
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
        _this.load();
      }, 200);
    },
    load() {
      if (this.currentUser.id === 0) {
        return;
      }

      if (this.isAllPostsFetched) {
        return;
      }

      this.$http.get(
          "/post/feed/" + this.currentUser.id + "?lpid=" + this.latestPostId,
          this.accessToken,
          (response) => {
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
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.load();
    });
  },
  created() {
    let _this = this;
    window.onscroll = function () {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight / 100 * 80) {
        _this.loadMore();
      }
    };

    this.$watch(
        () => this.currentUser.id,
        (to, from) => {
          this.load();
        }
    );

    eventBus.$on("new-post-created", (data) => {
      this.posts.unshift(data);
    });
  }
}
</script>

<style scoped>

</style>
