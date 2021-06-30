<template>
  <div class="profile">
    <div class="images">
      <div class="banner"></div>
      <div class="photo">
        <img v-if="profile.photo_url !== ''" :src="profile.photo_url" :alt="profile.login">
        <img v-if="profile.photo_url === ''" src="../images/404-avatar.png" alt="404">
      </div>
    </div>
    <div class="info">
      <div class="left">
        <div class="name">{{ profile.name }}</div>
        <div class="login">@{{ profile.login ? profile.login : "Profile not found" }}</div>
      </div>
      <div class="center">
        <router-link to="/user/edit" class="btn" v-if="currentUser.id === profile.id">edit profile</router-link>
        <a href="#" id="follow-btn" class="btn"
           v-if="currentUser.id !== profile.id && !profile.is_followed && profile.id"
           @click.prevent="follow" :disabled="!isFollowBtnEnabled">follow</a>
        <a href="#" id="unfollow-btn" class="btn" v-if="currentUser.id !== profile.id && profile.is_followed"
           @click.prevent="unfollow" :disabled="!isUnfollowBtnEnabled">unfollow</a>
      </div>
      <div class="right">
        <div class="following">
          <router-link :to="'/user/'+profile.login+'/following'">
            Following: <span>{{ profile.following }}</span>
          </router-link>
        </div>
        <div class="followers">
          <router-link :to="'/user/'+profile.login+'/followers'">
            Followers: <span>{{ profile.followers }}</span>
          </router-link>
        </div>
      </div>
    </div>
    <div class="bio">{{ profile.bio }}</div>
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
      isFollowBtnEnabled: true,
      isUnfollowBtnEnabled: true,
      profile: {
        "id": 0,
        "login": "",
        "name": "",
        "photo_url": "",
        "bio": "",
        "following": 0,
        "followers": 0,
        "is_followed": false
      },
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
    }
  },
  methods: {
    loadMore() {
      if (this.loadingTimeout != null) {
        clearTimeout(this.loadingTimeout);
      }

      let _this = this;
      this.loadingTimeout = setTimeout(function () {
        _this.loadPosts();
      }, 200);
    },
    loadPosts() {
      if (!this.profile.id) {
        return
      }

      if (this.isAllPostsFetched) {
        return
      }

      this.$http.get(
          "/post/user/" + this.profile.id + "?lpid=" + this.latestPostId,
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
          }, false
      );
    },
    load(login) {
      if (login !== this.profile.login) {
        this.profile.id = 0;
        this.profile.login = "";
        this.profile.name = "";
        this.profile.photo_url = "";
        this.profile.bio = "";
        this.profile.following = 0;
        this.profile.followers = 0;
        this.profile.is_followed = false;
        this.posts = [];
      }

      this.$http.get("/user/" + login, this.accessToken, (response) => {
        this.profile = response.data;
        this.loadPosts();
      }, true);
    },
    follow() {
      if (!this.isFollowBtnEnabled) {
        return;
      }

      this.isFollowBtnEnabled = false;

      let body = JSON.stringify({
        "user_id": this.profile.id
      });

      this.$http.sendRequest("/user/follow", "POST", this.accessToken, body, (response) => {
        this.profile.is_followed = true;
        this.profile.followers++;
        this.isFollowBtnEnabled = true;
      }, undefined, undefined, false);
    },
    unfollow() {
      if (!this.isUnfollowBtnEnabled) {
        return;
      }

      this.isUnfollowBtnEnabled = false;
      let body = JSON.stringify({
        "user_id": this.profile.id
      });

      // todo disable button when clicked

      this.$http.sendRequest("/user/unfollow", "POST", this.accessToken, body, (response) => {
        this.profile.followers--;
        this.profile.is_followed = false;
        this.isUnfollowBtnEnabled = true;
      }, undefined, undefined, false);
    }
  },
  created() {
    let _this = this;
    window.onscroll = function () {
      if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight / 100 * 80) {
        _this.loadMore();
      }
    };
  },
  watch: {
    $route(to, from) {
      console.log("Watch");
      this.load(to.params.login);
    }
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.load(to.params.login);
    })
  }
}
</script>

<style scoped>
/*@import url('https://fonts.googleapis.com/css?family=Solway&display=swap');*/
/*@import url('https://fonts.googleapis.com/css?family=Share+Tech+Mono&display=swap');*/

.images {
  position: relative;
}

.banner {
  height: 10em;
  overflow: hidden;
  background: #017eec33;
}

.banner img {
  width: 100%;
}

.images .photo {
  border-radius: 50%;
  width: 128px;
  height: 128px;
  position: absolute;
  left: 50%;
  top: 50%;
  margin-left: -64px;
  margin-top: -64px;
}

.images .photo img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
}

.info {
  display: flex;
  padding: 1em 2em 1em;
  border-left: 1px solid #e8e8e8;
  border-right: 1px solid #e8e8e8;
}

.info .left {
  flex: 1;
  text-align: right;
}

.info .center {
  flex: 1;
  text-align: center;
}

.info .right {
  flex: 1;
  text-align: left;
}

.name {
  font-weight: bold;
  font-size: 1.2em;
}

.login {
  color: #494949;
  font-size: 1.1em;
}

.bio {
  border-left: 1px solid #e8e8e8;
  border-right: 1px solid #e8e8e8;
  padding: 0 2em 1em;
  text-align: center;
}

.followers {
  margin-top: 0.3em;
}

</style>
