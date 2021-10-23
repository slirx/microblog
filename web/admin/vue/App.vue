<template>
  <div class="app">
    <div class="container" v-show="accessToken !== ''">
      <div class="menu" v-show="alert.type !== 'network-error'">
        <ul>
          <li>
            <router-link to="/">home</router-link>
          </li>
          <li><router-link to="/users">users</router-link></li>
          <li><a href="#" v-on:click="logout">logout</a></li>
        </ul>
      </div>
      <div class="content">
        <spinner v-show="isLoading"/>
        <div v-show="!isLoading">
          <router-view></router-view>
        </div>
      </div>
    </div>
    <auth @logout="logout" v-show="!isLoading && accessToken === ''"/>
    <spinner v-show="isLoading && accessToken === ''"/>
    <alert v-bind:message="alert.msg" v-bind:type="alert.type"/>
  </div>
</template>

<script>
import Auth from './Auth.vue'
import Alert from "../vue/Alert.vue"
import Spinner from './Spinner.vue'
import eventBus from "../js/event-bus.js";

export default {
  components: {
    "auth": Auth,
    "alert": Alert,
    "spinner": Spinner
  },
  computed: {
    accessToken() {
      return this.$store.state.jwt.access_token;
    },
    currentUser() {
      return this.$store.state.currentUser;
    },
    isLoading() {
      return this.$store.state.isLoading;
    },
    alert() {
      return this.$store.state.alert;
    },
  },
  methods: {
    logout(e) {
      e.preventDefault();
      eventBus.$emit("logout");
    }
  },
  created() {
    let jwt = JSON.parse(localStorage.getItem('jwt'));
    if (jwt && jwt.access_token) {
      this.$store.commit('setJWT', jwt)

      this.$http.get("/user/me", jwt.access_token, (response) => {
        this.$store.commit('setCurrentUser', response.data);
        console.log("setCurrentUser called");
      }, true);
    } else {
      this.$store.commit('isLoading', false);
    }
  },
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Share+Tech+Mono&display=swap');

.app {
  margin: 0 auto;
  width: 50%;
  font-family: 'Share Tech Mono', monospace;
}

.container {
  display: flex;
}

.menu {
  flex: 1;
}

.menu ul {
  list-style: none;
  border-left: 1px solid #e8e8e8;
  padding: 1em 0;
}

.menu ul li {
  padding: 0.5em 1em;
}

.content {
  flex: 3;
  padding-bottom: 2em;
}

</style>
