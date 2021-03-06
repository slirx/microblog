import Vue from 'vue'
import VueRouter from 'vue-router'
import VueMoment from 'vue-moment'
import App from '../vue/App.vue'
import Users from '../vue/Users.vue'
import http from './http.js'
import Vuex from 'vuex'
import "normalize.css"
import "../css/style.css"
import VueApollo from 'vue-apollo'
import ApolloClient from 'apollo-boost'

Vue.use(VueRouter);
Vue.use(Vuex);
Vue.use(VueMoment);
Vue.use(VueApollo);

const apolloClient = new ApolloClient({
    // You should use an absolute URL here
    uri: 'http://microblog.local:8080/graphql'
})

const apolloProvider = new VueApollo({
    defaultClient: apolloClient,
})

const store = new Vuex.Store({
    state: {
        jwt: {
            access_token: ""
        },
        currentUser: {
            "id": 0,
            "login": "",
            "name": "",
            "photo_url": "",
            "bio": "",
            "following": 0,
            "followers": 0
        },
        isLoading: true,
        alert: {
            type: "",
            msg: ""
        },

    },
    mutations: {
        setJWT(state, value) {
            localStorage.setItem("jwt", JSON.stringify({
                access_token: value.access_token
            }));
            state.jwt = value;
        },
        deleteJWT(state) {
            localStorage.removeItem("jwt");
            state.jwt.access_token = "";
        },
        setCurrentUser(state, value) {
            state.currentUser = value;
        },
        isLoading(state, value) {
            state.isLoading = value;
        },
        alert(state, data) {
            state.alert.type = data.type;
            state.alert.msg = data.message;

            setTimeout(function () {
                state.alert.msg = "";
            }, 5000);
        }
    }
})

Vue.use(http, {store})

const routes = [
    {path: '/users', component: Users},
    // {path: '/user/edit', component: ProfileEdit},
    // {path: '/user/:login', component: Profile},
    // {path: '/user/:login/following', component: Following},
    // {path: '/user/:login/followers', component: Followers},
    // {path: '/search', component: Search}
    // path: '*' todo add 404 page
]

const router = new VueRouter({
    mode: 'history',
    routes: routes
})

const app = new Vue({
    el: "#app",
    store,
    router,
    apolloProvider,
    template: '<App/>',
    components: {
        App
    }
});
