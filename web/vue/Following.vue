<template>
  <div class="following">
    <div class="login">@{{ login }}</div>
    <div class="tabs">
      <div class="tab active">Following</div>
      <div class="tab">
        <router-link :to="'/user/' + login + '/followers'">Followers</router-link>
      </div>
    </div>

    <div class="empty" v-if="!following.length">no following yet..</div>
    <div class="item" v-for="user in following">
      <div class="photo">
        <img :src="user.photo_url" :alt="user.login">
      </div>
      <div class="user">
        {{ user.name }} @
        <router-link :to="'/user/' + user.login">{{ user.login }}</router-link>
        <div class="bio">{{ user.bio }}</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      login: "",
      following: [],
      latestFollowingId: 0,
      loadingTimeout: null,
      isAllFollowingFetched: false
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
        _this.load(_this.login);
      }, 200);
    },
    load(login) {
      if (this.login !== login) {
        this.login = login;
        this.following = [];
        this.latestFollowerId = 0;
        this.loadingTimeout = null;
        this.isAllFollowingFetched = false;
      }

      if (this.isAllFollowingFetched) {
        return;
      }

      this.$http.get(
          "/user/" + login + "/following?lfid=" + this.latestFollowingId,
          this.accessToken,
          (response) => {
            if (this.following.length === response.data.total) {
              this.isAllFollowingFetched = true;
              return;
            }

            this.following.push(...response.data.following);

            if (this.following.length === response.data.total) {
              this.isAllFollowingFetched = true;
              return;
            }

            if (this.following.length > 0) {
              this.latestFollowingId = this.following[this.following.length - 1].follower_id;
            }
          }, false);
    },

    // follow() {
    //   if (!this.isFollowBtnEnabled) {
    //     return;
    //   }
    //
    //   this.isFollowBtnEnabled = false;
    //
    //   let body = JSON.stringify({
    //     "user_id": this.profile.id
    //   });
    //
    //   this.$http.sendRequest("/user/follow", "POST", this.accessToken, body, (response) => {
    //     this.profile.is_followed = true;
    //     this.profile.followers++;
    //     this.isFollowBtnEnabled = true;
    //   }, undefined, undefined, false);
    // },
    // unfollow() {
    //   if (!this.isUnfollowBtnEnabled) {
    //     return;
    //   }
    //
    //   this.isUnfollowBtnEnabled = false;
    //   let body = JSON.stringify({
    //     "user_id": this.profile.id
    //   });
    //
    //   // todo disable button when clicked
    //
    //   this.$http.sendRequest("/user/unfollow", "POST", this.accessToken, body, (response) => {
    //     this.profile.followers--;
    //     this.profile.is_followed = false;
    //     this.isUnfollowBtnEnabled = true;
    //   }, undefined, undefined, false);
    // }
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
.following {
  padding-top: 2em;
}
</style>
