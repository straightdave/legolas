// legolas Vue app
// Dave @ Nov-2017

var nav_tmpl = `
<nav>
    <img id="logo" src="img/grab_logo.png" alt="grab_logo">
    <ul>
        <li><a v-on:click.stop.prevent="cli(0)" v-bind:class="{active: activeItem == 0}">Cases</a></li>
        <li><a v-on:click.stop.prevent="cli(1)" v-bind:class="{active: activeItem == 1}">Helps</a></li>
    </ul>
</nav>
`

Vue.component('app-nav', {
  template: nav_tmpl,
  data: function () {
    return {
        activeItem: 0
    }
  },
  methods: {
    cli: function (item) {
        this.activeItem = item
    }
  }
})

var search_tmpl = `
<div class="search-wrapper">
    <input type="text" v-model="word" placeholder="Search">
    <button v-on:click.stop.prevent="go"><i class="fa fa-search"></i></button>
</div>
`

Vue.component('app-search', {
    template: search_tmpl,
    data: function () {
        return {
            word: ""
        }
    },
    methods: {
        go: function () {
            alert('you are search: ' + this.word)
        }
    }
})

var app = new Vue({
    el: '#app'
})
