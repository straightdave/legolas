<template>
<div id="app-sidebar">
    <app-search @do-filter="doFilter" />

    <div class="titlebar">
        <span>{{ title }}</span>
        <a id="create" v-on:click.stop.prevent="createItem"><i class="fa fa-plus"></i></a>
    </div>

    <div id="list" v-if="items != null && items.length">
        <app-sidebar-item
            v-for="item in items"
            :key="item._id"
            :item-object="item"
            @view-item-detail="$emit('item-clicked', item)" />
    </div>
</div>
</template>

<script>
import Vue         from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

import AppSearch      from './app-search.vue'
import AppSidebarItem from './app-sidebar-item.vue'

export default {
    components: {AppSearch, AppSidebarItem},
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
    mounted() {
        this.getItems('')
    },
    methods: {
        getItems(filter) {
            console.log('filter: ' + filter)

            const iType = this.whatFor

            var url = `/${iType}`
            if (filter && filter.length > 2) {
                url += `/f/${encodeURI(filter)}`
            }

            console.log('sidebar >>> url: ' + url)
            this.$http.get(url).then(
                resp => {
                    var data = resp.body
                    console.log(`got by ${url}: ${ JSON.stringify(data) }`)
                    if (data && data.length > 0) {
                        this.items = data
                    }
                    else {
                        this.items = []
                    }
                },
                resp => {
                    console.log(`calling ${url} failed. Resp: ${ JSON.stringify(resp) }`)
                    this.items = []
                }
            )
        },

        createItem() {
            var newOne = {
                path: 'default',
                name: `${ this.whatFor }-new`,
                desc: `new item of ${ this.whatFor }`,
                created_at: new Date(),
                isNew: true
            }
            this.items.unshift(newOne)
            // pop-up event to parent, let new case content show in panel
            this.$emit('item-clicked', newOne)
        },

        doFilter(word) {
            console.log('filter word: ' + word)
            this.items = null

            // for now, use the word as a simple name-based filter
            this.getItems(word)

            // ask outter component to close right panel
            this.$emit('close-right-panel')
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
