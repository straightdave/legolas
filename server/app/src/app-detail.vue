<!-- app detail -->
<template>
<div id="app-detail">
    <div id="detail-header">
        <p><span class="group">$/{{ group }}</span></p>
        <p><span class="name">{{ name }}</span></p>
        <p><span class="desc">{{ desc }}</span></p>
    </div>
    <div id="nav_in_case">
        <ul>
            <li><a v-on:click.stop.prevent="cliTab(0)" v-bind:class="{active: activeItem == 0}">Actions</a></li>
            <li><a v-on:click.stop.prevent="cliTab(1)" v-bind:class="{active: activeItem == 1}">History</a></li>
            <li><a v-on:click.stop.prevent="cliTab(2)" v-bind:class="{active: activeItem == 2}">Tracing</a></li>
        </ul>
    </div>
    <div id="detail-panel">
        <div v-if="activeItem == 0">
            <div id="actionlist" v-if="actions.length">
                <AppAction
                    v-for="a in actions"
                    :key="a.id"
                    :action="a"
                    v-on:click="actionClicked(a.id)"
                />
            </div>
        </div>
        <div v-else-if="activeItem == 1">
            history
        </div>
        <div v-else-if="activeItem == 2">
            Data tracing
        </div>
    </div>
</div>
</template>

<script>
import AppAction from './app-action.vue'

var AppDetail = Vue.extend({
    components: { AppAction },
    props: {
        case: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            id: this.case.id,
            name: this.case.name,
            desc: this.case.desc,
            group: this.case.group,
            activeItem: 0,


            actions: [
                {
                    name: "action-1",
                    casePath: this.case.group,
                    caseName: this.case.name,
                    index: 1
                },
                {
                    name: "action-2",
                    casePath: this.case.group,
                    caseName: this.case.name,
                    index: 2
                },
                {
                    name: "action-3",
                    casePath: this.case.group,
                    caseName: this.case.name,
                    index: 3
                }
            ]
        }
    },
    methods: {
        cliTab: function (item) {
            this.activeItem = item
        },
        actionClicked: function (id) {
            alert('action clicked: ' + id)
        }
    }
})
export default AppDetail

</script>

<style scoped>

div#detail-header {
    height: 72px;
    overflow: hidden;
}

span.group, span.desc {
    color: gray;
}

span.name {
    font-size: 30px;
    font-weight: 300;
}

div#nav_in_case ul {
    display: inline-flex;
    list-style: none;
    margin-top: 15px;
    margin-bottom: 15px;
}

div#nav_in_case ul li {
    padding-right: 20px;
}

div#nav_in_case ul li a {
    font-weight: 300;
    text-decoration: none;
    color: gray;
    cursor: pointer;
    display: block;
}

div#nav_in_case li a:hover,
div#nav_in_case li a.active {
    border-bottom: solid 2px #00B140;
    padding-bottom: 8px;
}

div#actionlist {
    width: 300px;
}


</style>
