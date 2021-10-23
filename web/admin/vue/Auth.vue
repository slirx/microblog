<template>
  <div class="auth">
    <div class="left">
      <img src="../images/bird.jpg" alt="bird"/>
    </div>
    <div class="right">
      <h1>Micro Blog Admin Panel</h1>
      <form action="#" v-on:submit.prevent="submitSignIn">
        <div>
          <div class="input-container">
            <input placeholder="Login" class="input" type="text" autocomplete="off" v-model="signIn.login"/>
            <label class="label">Login</label>
          </div>
          <div class="input-container">
            <input placeholder="Password" class="input" type="password" autocomplete="off" v-model="signIn.password"/>
            <label class="label">Password</label>
          </div>
        </div>
        <div><input type="submit" class="btn" value="Sign In"
                    :disabled="this.signIn.login.trim() === '' || this.signIn.password.trim() === ''"></div>
      </form>
    </div>
  </div>
</template>

<script>
import eventBus from "../js/event-bus.js";

const adminBaseURL = 'http://microblog.local:9007';

export default {
  data() {
    return {
      signIn: {
        login: "",
        password: ""
      },
    }
  },
  methods: {
    submitSignIn() {
      let body = JSON.stringify({
        "login": this.signIn.login,
        "password": this.signIn.password,
      });

      this.signIn.password = ""

      this.$http.sendRequest(adminBaseURL + "/auth/admin/login", "POST", "", body, (response) => {
        this.$store.commit('setJWT', response.data);
        this.$http.get("/user/me", response.data.access_token, (response) => {
          this.$store.commit('setCurrentUser', response.data);
          this.$router.push('users');
        }, true);
      }, undefined, undefined, true);
    },
  },
  created() {
    eventBus.$once("logout", () => {
      this.$store.commit('deleteJWT');
    });
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Solway&display=swap');
@import url('https://fonts.googleapis.com/css?family=Share+Tech+Mono&display=swap');

.auth {
  display: flex;
  font-family: 'Share Tech Mono', monospace;
}

.auth input {
  font-family: 'Share Tech Mono', monospace;
  padding: 0.3em;
}

.auth input:not([type=submit]) {
  width: 15em;
}

.auth img {
  width: 100%;
}

.auth h1 {
  font-family: 'Solway', serif;
}

.auth .left,
.auth .right {
  flex: 1;
  padding: 3em 2em 0;
}

.auth a:hover,
.auth a {
  color: black;
}

.auth a:hover {
  text-decoration: none;
}
</style>

