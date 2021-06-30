import Vue from 'vue'
import VueRouter from 'vue-router'
import VueMoment from 'vue-moment'
import App from '../vue/App.vue'
import Home from '../vue/Home.vue'
import Profile from '../vue/Profile.vue'
import Bookmarks from "../vue/Bookmarks.vue";
import Following from "../vue/Following.vue";
import Followers from "../vue/Followers.vue";
import ProfileEdit from "../vue/ProfileEdit.vue";
import Search from "../vue/Search.vue"
import http from './http.js'
import Vuex from 'vuex'
import "normalize.css"
import "../css/style.css"

Vue.use(VueRouter)
Vue.use(Vuex)
Vue.use(VueMoment);

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
                state.alert.type = "";
                state.alert.msg = "";
            }, 5000);
        }
    }
})

Vue.use(http, {store})

const routes = [
    {path: '/', component: Home},
    {path: '/user/edit', component: ProfileEdit},
    {path: '/user/:login', component: Profile},
    {path: '/user/:login/following', component: Following},
    {path: '/user/:login/followers', component: Followers},
    {path: '/bookmarks', component: Bookmarks},
    {path: '/search', component: Search}
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
    template: '<App/>',
    components: {
        App
    }
});
