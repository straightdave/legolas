<!-- action listed item -->
<template>
<div id="action" v-on:click="$emit('actionClicked', action.name)">
    <div id="titlebar">
        <p>
            <span id="name">
                <span v-if="action.isNew" class="newitem">
                    <i class="fa fa-circle"></i>
                </span>
                <span v-else class="green">
                    <i class="fa fa-check"></i>
                </span>
                {{ action.name }}
            </span>
        </p>
        <p><span id="desc">{{ action.desc }}</span></p>
    </div>
</div>
</template>

<script>
import Vue from 'vue'

var AppAction = Vue.extend({
    props: {
        action: {
            type: Object,
            required: true
        }
    }
})
export default AppAction
</script>

<style scoped>
div#action {
    height: 40px;
    margin-bottom: 5px;
    padding:10px;
    cursor: pointer;
    background-color: #ececec;
    overflow: hidden;
    border-left: solid 10px #0366d6;
}

.green {
    color: #00B140;
}
.newitem {
    font-size: small;
    color: #4AACF2;
}

div#titlebar span#name {
    font-size: 20px;
}
div#titlebar span#desc {
    font-size: small;
    color: gray;
}
</style>
