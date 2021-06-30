<template>
  <div v-if="message !== ''" class="alert">
    <span v-bind:class="type"><span class="close" v-on:click="close">&#10007;</span> {{ message }}</span>
  </div>
</template>

<script>
export default {
  name: "Alert",
  props: {
    message: String,
    type: String // success or error
  },
  methods: {
    close() {
      this.$store.commit('alert', {
        type: "",
        message: ""
      });
    }
  }
}
</script>

<style scoped>
.alert {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
}

.success,
.network-error,
.error {
  display: inline-block;
  margin-top: 0.5em;
  font-size: 0.9em;
  padding: 0.2em;
  border-radius: 5px;
}

.success {
  background-color: #96f98f;
}

.network-error,
.error {
  background-color: #f98f8f;
}

.close {
  cursor: pointer;
  font-size: 1.5em;
  color: #3c3c3c;
  line-height: 1.5em;
  vertical-align: middle;
}
</style>
