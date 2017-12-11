<template>
<div>
    <nav>
        <img id="logo" src="img/grab_logo.png" @click.stop.prevent="logoClick">
        <ul>
            <li><router-link to="/">Cases</router-link></li>
            <li><router-link to="/templates">Templates</router-link></li>
            <li><router-link to="/help">Help</router-link></li>
        </ul>

        <div id="user-profile" v-if="isSignedIn">
            <span @click.stop.prevent="signOut">{{ userProfile.getName() }}</span>
            <img :src="userProfile.getImageUrl()" />
        </div>
        <div id="user-signin" v-else>
            <span>Sign-In</span>
        </div>
    </nav>
    <router-view></router-view>
</div>
</template>

<script>
import Vue       from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

import AppCaseMain     from './app-case-main.vue'
import AppTemplateMain from './app-template-main.vue'
import AppHelpMain     from './app-help-main.vue'

const router = new VueRouter({
    routes: [
        {path: '/', component: AppCaseMain},
        {path: '/templates', component: AppTemplateMain},
        {path: '/help', component: AppHelpMain}
    ]
})

const App = Vue.extend({
    router: router,

    data() {
        return {
            authObj: null,
            currentUser: null
        }
    },

    computed: {
        userProfile() {
            if (this.currentUser) {
                return this.currentUser.getBasicProfile()
            }
            return null
        },

        isSignedIn() {
            if (this.authObj) {
                return (this.authObj.isSignedIn.get() && this.currentUser)
            }
            return false
        }
    },

    created() {
        console.log('created the app')

        if (!window.gapi) {
            console.log('cannot load gapi, in <script>?')
            return
        }

        window.gapi.load('auth2', () => {
            this.initAuthObject()
            this.initSignInHandler()
        })
    },

    methods: {
        initAuthObject() {
            var _auth2 = window.gapi.auth2.init({
                client_id: '345267627311-sslkr40iac0fhmrsci7lfg7ja5trjrbm.apps.googleusercontent.com',
                cookiepolicy: 'single_host_origin'
            })
            this.authObj = _auth2
            console.log('successfully init the auth2 object')
        },

        initSignInHandler() {
            if (!this.authObj) {
                console.log('authObj is not ready, cannot render sign-in button')
                return
            }

            var signInElement = document.getElementById('user-signin')
            this.authObj.attachClickHandler(
                signInElement, {prompt: 'select_account'},
                googleUser => {
                    console.log('login as ' + googleUser.getBasicProfile().getName())
                    this.currentUser = googleUser
                },
                error => alert(JSON.stringify(error, undefined, 2))
            )
            console.log('successfully attached sign-in handler')
        },

        signOut() {
            if (!this.authObj) {
                console.log("no authobj for now, no need to sign out")
                return
            }

            console.log('signing out...')
            console.log('log out current user: ' + this.currentUser.getBasicProfile().getName())

            window.gapi.auth2.getAuthInstance().signOut().then(() => {
                console.log('User signed out hopefully.')
                this.currentUser = null
            })
        },

        logoClick() {
            if (this.currentUser) {
                console.log('current user: ' + this.currentUser.getBasicProfile().getEmail())
            }
        }
    }
})
export default App
</script>

<style scoped>
nav {
    height: 50px;
    width: 100%;
    background-color: #ececec;
    position: relative;
}

div#user-signin, div#user-profile {
    position: absolute;
    right: 0;
    top: 0;

    font-size: 18px;
    color: gray;
    cursor: pointer;
}

div#user-signin span, div#user-profile span {
    display: block;
    float: left;

    margin-top: 15px;
    margin-right: 10px
}

div#user-profile img {
    display: block;
    float: left;
    width: 40px;
    height: 40px;

    margin: 5px;
}

img#logo {
    height: 34px;
    margin: 8px 0 3px 40px;
}

nav ul {
    display: inline-flex;
    position: absolute;
    list-style: none;
    margin-left: 40px;
    margin-top: 15px;
    font-size: 1.1em;
}

nav ul li {
    padding-right: 20px;
}

nav ul li a {
    text-decoration: none;
    color: gray;
    cursor: pointer;
}

a.router-link-exact-active {
    display: block;
    border-bottom: solid 4px gray;
    padding-bottom: 10px;
}
</style>
