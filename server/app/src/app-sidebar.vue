<template>
<div id="app-sidebar">
    <app-search @do-filter="doFilter" />

    <div class="titlebar">
        <span>{{ title }}</span>
        <a id="create" v-on:click.stop.prevent="createItem"><i class="fa fa-plus"></i></a>
    </div>

    <div id="list" v-if="whatFor === 'cases' && items != null && items.length">
        <app-case
            v-for="item in items"
            :key="item.name"
            :case="item"
            @view-case="$emit('item-clicked', item)"
        />
    </div>

    <div id="list" v-if="whatFor === 'templates'">
        template list
    </div>

</div>
</template>

<script>
import AppSearch from './app-search.vue'
import AppCase from './app-case.vue'
import $ from 'jquery'

export default {
    components: {AppSearch, AppCase},
    props: {
        whatFor: {
            type: String,
            required: true
        }
    },
    data() {
        return {
            items: []
        }
    },
    computed: {
        title() {
            var result = ""
            switch (this.whatFor) {
                case "cases":
                    result = "My Cases"
                    break
                case "templates":
                    result = "My Templates"
                    break
                default:
                    result = "unknown"
            }
            return result
        }
    },
    mounted: function () {
        var self = this // critical!
        if (this.whatFor === 'cases') {
            $.get('/cases', function (data) {
                self.items = data
            })
        }
        else if (this.whatFor === 'templates') {
        }
    },
    methods: {
        createItem() {
            if (this.whatFor === 'cases') {
                var newOne = {
                    path: 'default',
                    name: 'case-new',
                    desc: 'new case',
                    isNew: true
                }
                this.items.unshift(newOne)
                // pop-up event to parent, let new case content show in panel
                this.$emit('item-clicked', newOne)
            }
            else if (this.whatFor === 'templates') {

            }
        },

        doFilter(word) {
            console.log('filter word: ' + word)
        }
    }
}
</script>


<style scoped>
div#app-sidebar {
    width: 340px;
    border-right: solid 4px #ececec;
    float: left;
}

div.titlebar {
    width: 300px;
    margin: 0 auto;
}

div.titlebar span {
    display: inline-block;
    margin-right: 10px;
    padding-bottom: 8px;
}

a#create {
    display: inline-block;
    text-decoration: none;
    color: gray;
    cursor: pointer;
    font-size: 0.7em;
    font-weight: 200px;
}

div#list {
    background-color: #ececec;
    padding: 4px 0;
}
</style>
