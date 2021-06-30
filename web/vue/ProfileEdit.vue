<template>
  <div class="profile-edit">
    <h3>Edit profile</h3>
    <form action="#" method="POST" v-on:submit.prevent="submitProfile">
      <div>
        <div class="input-container">
          <input placeholder="Name" class="input" type="text" autocomplete="off" v-model="profile.name"/>
          <label class="label">Name</label>
        </div>
        <div class="input-container">
          <input placeholder="Avatar" class="input" type="file" ref="avatar" autocomplete="off"/>
          <label class="label">Avatar</label>
        </div>
        <div class="input-container">
          <textarea placeholder="Bio" class="input" rows="3" v-model="profile.bio"></textarea>
          <label class="label">Bio</label>
        </div>
      </div>
      <div><input type="submit" value="save" class="btn"></div>
    </form>
  </div>
</template>

<script>
import eventBus from "../js/event-bus";

export default {
  data() {
    return {
      profile: {
        name: "",
        bio: "",
      }
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
    submitProfile() {
      let req = JSON.stringify({
        name: this.profile.name,
        bio: this.profile.bio
      });

      if (this.$refs.avatar.files.length > 0) {
        let data = new FormData();
        data.append('image', this.$refs.avatar.files[0]);
        data.append('service', 'user');

        this.$http.sendRequest("/media", "POST", this.accessToken, data, (response) => {
          console.log("avatar has been changed");
          // todo what to do here?
          // this.$store.commit('alert', {
          //   type: response.type,
          //   message: response.message
          // });
        }, undefined, undefined, true);
      }

      this.$http.sendRequest("/user", "PATCH", this.accessToken, req, (response) => {
        this.$store.commit('alert', {
          type: response.type,
          message: response.message
        });

        this.$router.push('/user/' + this.currentUser.login);
      }, undefined, undefined, true);
    }
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.$http.get("/user/me", vm.accessToken, (response) => {
        vm.profile = response.data;
      }, true);
    })
  },
}
</script>

<style scoped>

</style>
