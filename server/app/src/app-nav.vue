<!-- search -->
<template>
<nav>
    <img id="logo" src="img/grab_logo.png" alt="grab_logo">
    <ul>
        <li><a v-on:click.stop.prevent="cli(0)" v-bind:class="{active: activeItem == 0}">Cases</a></li>
        <li><a v-on:click.stop.prevent="cli(1)" v-bind:class="{active: activeItem == 1}">Helps</a></li>
    </ul>
</nav>
</template>

<script>

export default {
    data() {
        return {
            activeItem: 0
        }
    },
    methods: {
        cli: function (item) {
            this.activeItem = item
        }
    }
}

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

li a:hover, a.active {
    display: block;
    border-bottom: solid 4px gray;
    padding-bottom: 10px;
}
</style>
