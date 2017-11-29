<template>
<div>
    <nav>
        <img id="logo" src="img/grab_logo.png" alt="grab_logo">
        <ul>
            <li><router-link to="/">Cases</router-link></li>
            <li><router-link to="/templates">Templates</router-link></li>
            <li><router-link to="/help">Help</router-link></li>
        </ul>
    </nav>
    <router-view></router-view>
</div>
</template>

<script>
import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

import AppCaseMain from './app-case-main.vue'
import AppTemplateMain from './app-template-main.vue'

const AppHelpMain = {
    template: '<h1>fucking help</h1>'
}

const router = new VueRouter({
    routes: [
        {path: '/', component: AppCaseMain},
        {path: '/templates', component: AppTemplateMain},
        {path: '/help', component: AppHelpMain}
    ]
})

var App = Vue.extend({
    router: router
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

nav img#logo {
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
