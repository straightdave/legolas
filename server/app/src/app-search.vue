<template>
<div id="search-wrapper">
    <input type="text" v-model.trim="word" placeholder="Commands" @keyup.enter="$emit('do-filter', word)">
    <button @click="$emit('do-filter', word)"><i class="fa fa-search"></i></button>
</div>
</template>

<script>
export default {
    data() {
        return {
            word: ""
        }
    }
}
</script>

<style scoped>
div#search-wrapper {
    display: block;
    width: 300px;
    margin: 30px auto;
    border: 2px solid #ececec;
    border-radius: 5px;
    position: relative;
}
div#search-wrapper input,
div#search-wrapper button {
    background-color: transparent;
    border: 0;
    font-size: 20px;
    padding: 5px;
}
div#search-wrapper button {
    position: absolute;
    right: 0;
    top: 0;
    cursor: pointer;
    color: #ececec;
}
div#search-wrapper input {
    width: 250px;
}
</style>
