<template>
  <div class="auth">
    <div class="left">
      <img src="../images/bird.jpg" alt="bird"/>
    </div>
    <div class="right">
      <h1>Micro Blog</h1>
      <form action="#" v-on:submit.prevent="submitSignIn" v-if="openedForm === 'sign-in'">
        <h3>Sign In</h3>
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

      <form action="#" v-on:submit.prevent="submitSignUp" v-if="openedForm === 'sign-up'">
        <h3>Sign Up</h3>
        <div>
          <div class="input-container">
            <input placeholder="Login" class="input" type="text" autocomplete="off" v-model="signUp.login"/>
            <label class="label">Login</label>
          </div>
          <div class="input-container">
            <input placeholder="Email" class="input" type="email" autocomplete="off" v-model="signUp.email"/>
            <label class="label">Email</label>
          </div>
        </div>
        <div>
          <input type="submit" value="Sign Up"
                 :disabled="signUp.disabled || signUp.login.trim() === '' || signUp.email.trim() === ''" class="btn">
        </div>
      </form>

      <form action="#" v-on:submit.prevent="submitSignUpConfirmation" v-if="openedForm === 'sign-up-confirmation'">
        <h3>Sign Up Confirmation</h3>
        <div>
          <div class="input-container">
            <input placeholder="Email confirmation code" class="input" type="text" autocomplete="off"
                   v-model="signUpConfirmation.code"/>
            <label class="label">Email confirmation code</label>
          </div>
          <div class="input-container">
            <input placeholder="Password" class="input" type="password" autocomplete="off"
                   v-model="signUpConfirmation.password"/>
            <label class="label">Password</label>
          </div>
          <div class="input-container">
            <input placeholder="Password confirmation" class="input" type="password" autocomplete="off"
                   v-model="signUpConfirmation.passwordConfirmation"/>
            <label class="label">Password confirmation</label>
          </div>
        </div>
        <div>
          <input type="submit" value="Confirm"
                 :disabled="this.signUpConfirmation.code.trim() === '' || this.signUpConfirmation.password.trim() === '' || this.signUpConfirmation.password.trim() !== this.signUpConfirmation.passwordConfirmation.trim()"
                 class="btn">
        </div>
      </form>
      <div>
        or
        <a href="#sign-in" v-on:click.prevent="tabSignIn" v-if="openedForm == 'sign-up'">Sign In</a>
        <a href="#sign-up" v-on:click.prevent="tabSignUp" v-if="openedForm == 'sign-in'">Sign Up</a>
      </div>
    </div>
  </div>
</template>

<script>
import eventBus from "../js/event-bus.js";

export default {
  data() {
    return {
      openedForm: "sign-in",
      isSignUpEnabled: true,
      signIn: {
        login: "",
        password: ""
      },
      signUp: {
        login: "",
        email: "",
        disabled: false
      },
      signUpConfirmation: {
        email: "",
        code: "",
        password: "",
        passwordConfirmation: ""
      }
    }
  },
  methods: {
    submitSignIn() {
      let body = JSON.stringify({
        "login": this.signIn.login,
        "password": this.signIn.password,
      });

      // this.signIn.login = ""
      this.signIn.password = ""

      this.$http.sendRequest("/auth/login", "POST", "", body, (response) => {
        this.$store.commit('setJWT', response.data);
        this.$http.get("/user/me", response.data.access_token, (response) => {
          this.$store.commit('setCurrentUser', response.data);
        }, true);
      }, undefined, undefined, true);
    },
    submitSignUp() {
      let body = JSON.stringify({
        "login": this.signUp.login,
        "email": this.signUp.email,
      });

      this.signUpConfirmation.email = this.signUp.email;
      // this.signUp.login = "";
      // this.signUp.email = "";
      this.signUp.disabled = true;

      this.$http.sendRequest("/registration", "POST", "", body, (response) => {
        this.$store.commit('alert', {
          type: response.type,
          message: response.message
        });
        this.openedForm = 'sign-up-confirmation';
        this.signUp.disabled = false;
      }, (response) => {
        this.signUp.disabled = false;
      }, undefined, true);
    },
    submitSignUpConfirmation() {
      let body = JSON.stringify({
        "email": this.signUpConfirmation.email,
        "code": parseInt(this.signUpConfirmation.code),
        "password": this.signUpConfirmation.password,
        "password_confirmation": this.signUpConfirmation.passwordConfirmation,
      });

      this.signUpConfirmation.code = "";
      this.signUpConfirmation.password = "";
      this.signUpConfirmation.passwordConfirmation = "";

      this.$http.sendRequest("/registration/confirm", "POST", "", body, (response) => {
        this.$store.commit('alert', {
          type: response.type,
          message: response.message
        });
        this.openedForm = 'sign-in';
        this.signIn.login = this.signUp.login;
      }, undefined, undefined, true);
    },
    tabSignIn() {
      this.openedForm = "sign-in";
    },
    tabSignUp() {
      this.openedForm = "sign-up";
    }
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

/*.auth form > div {*/
/*  margin-bottom: 0.5em;*/
/*}*/

.auth a:hover,
.auth a {
  color: black;
}

.auth a:hover {
  text-decoration: none;
}
</style>

