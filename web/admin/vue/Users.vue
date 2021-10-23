<template>
  <div class="users">
    <h3 class="center">Users</h3>
    <div v-if="users.total > 0">
      <div class="item" v-for="user in users.users">
        <div class="photo">
          <img :src="user.photo_url" :alt="user.login">
        </div>
        {{ user.name }} @
        <router-link :to="'/user/' + user.login">{{ user.login }}</router-link>
        · {{ user.created_at | moment("from") }}
      </div>
    </div>
    <div class="load-more">
      <a v-if="!isAllUsersFetched" href="#" class="btn" id="btn-load-more" @click.prevent="loadMore">load more</a>
    </div>
  </div>
</template>

<script>
import gql from 'graphql-tag';

export default {
  data() {
    return {
      users: {
        users: [],
        total: 0
      },
      usersList: {
        users: [],
        total: 0
      },
      latestUserId: 2147483647,
      loadingTimeout: null,
      isAllUsersFetched: false
    }
  },
  methods: {
    loadMore() {
      if (this.isAllUsersFetched) {
        return;
      }

      if (this.loadingTimeout != null) {
        clearTimeout(this.loadingTimeout);
      }

      let _this = this;

      this.loadingTimeout = setTimeout(function () {
        let latestUserId = _this.users.users[_this.users.users.length - 1].id;
        console.log("latestUserId:", latestUserId);

        _this.$apollo.queries.users.fetchMore({
          variables: {
            latestUserId: _this.users.users[_this.users.users.length - 1].id,
          },
          updateQuery: (previousResult, {fetchMoreResult}) => {
            if (fetchMoreResult.users.users.length === 0) {
              _this.isAllUsersFetched = true;
            }

            return {
              users: {
                __typename: previousResult.users.__typename,
                users: [...previousResult.users.users, ...fetchMoreResult.users.users],
                total: fetchMoreResult.users.total
              },
            }
          },
        })
      }, 150);
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
  apollo: {
    users: {
      query: gql`query users($latestUserId: Int!) {
        users(latestUserID: $latestUserId) {
          total,
          users {
            id
            login
            name
            photo_url
          }
        }
      }
       `,
      variables() {
        return {
          latestUserId: this.latestUserId,
        }
      }
    }
  }
}
</script>

<style scoped>
.users .item {
  display: flex;
  align-items: center;
}

.users .item .photo {
  margin-right: 0.5em;
  margin-top: 0.5em;
}

.users .item .photo img {
  width: 2em;
}

.load-more {
  text-align: center;
  margin: 1em 0;
}
</style>
