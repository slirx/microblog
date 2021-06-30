<template>
  <div class="followers">
    <div class="login">@{{ login }}</div>
    <div class="tabs">
      <div class="tab">
        <router-link :to="'/user/' + login + '/following'">Following</router-link>
      </div>
      <div class="tab active">Followers</div>
    </div>

    <div class="empty" v-if="!followers.length">no followers yet..</div>
    <div class="item" v-for="user in followers">
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
      followers: [],
      latestFollowerId: 0,
      loadingTimeout: null,
      isAllFollowersFetched: false
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
        this.followers = [];
        this.latestFollowerId = 0;
        this.loadingTimeout = null;
        this.isAllFollowersFetched = false;
      }

      if (this.isAllFollowersFetched) {
        return;
      }

      this.$http.get(
          "/user/" + login + "/followers?lfid=" + this.latestFollowerId,
          this.accessToken,
          (response) => {
            if (this.followers.length === response.data.total) {
              this.isAllFollowersFetched = true;
              return;
            }

            this.followers.push(...response.data.followers);

            if (this.followers.length === response.data.total) {
              this.isAllFollowersFetched = true;
              return;
            }

            if (this.followers.length > 0) {
              this.latestFollowerId = this.followers[this.followers.length - 1].follower_id;
            }
          }, false);
    },
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
.followers {
  padding-top: 2em;
}
</style>
