<template>
  <div class="new-post">
    <div class="new">
      <form action="#" v-on:submit.prevent="submitNewPost">
        <div class="text">
          <textarea name="message" id="message" rows="2" placeholder="what's happening?"
                    v-model="newPost.text"></textarea>
        </div>
        <div class="buttons">

          <input type="submit" :disabled="this.newPost.text.trim() === ''" value="post">
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import eventBus from "../js/event-bus";

export default {
  data() {
    return {
      newPost: {
        text: ""
      }
    }
  },
  computed: {
    accessToken() {
      return this.$store.state.jwt.access_token;
    }
  },
  methods: {
    submitNewPost() {
      if (this.newPost.text.trim() === '') {
        return;
      }

      let body = JSON.stringify({
        "text": this.newPost.text
      });

      this.newPost.text = ""

      this.$http.sendRequest("/post", "POST", this.accessToken, body, (response) => {
        // todo emit event to add this post to the opened feed
        // console.log("creating a new post", response);
        eventBus.$emit('new-post-created', response.data);
      }, undefined, undefined, true);
    }
  }
}
</script>

<style scoped>

.new-post .new {
  position: relative;
}

.new-post .new form {
  display: flex;
}

.new-post .new form .text {
  flex: 9;
}

.new-post .new form .buttons {
  flex: 1;
}

.new-post .new textarea {
  resize: none;
  width: 100%;
  padding: 0.5em;
  border: 1px solid #e8e8e8;
  margin: 0;
  outline: none;
}

.new-post .new .buttons input {
  border: none;
  color: #fff;
  background: #343434;
  width: 100%;
  height: 100%;
}

.new-post .new .buttons input[disabled=disabled] {
  color: #909090;
}

.new-post .new .buttons input:hover {
  cursor: pointer;
  background: #000000;
}

.new-post .new .buttons input[disabled=disabled]:hover {
  cursor: default;
  background: #343434;
}

</style>
