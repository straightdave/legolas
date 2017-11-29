<!-- case variable -->
<template>
<div id="app-case-variable">

</div>
</template>

<script>
export default {
    props: {
        caseInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            localCaseInfo: JSON.parse(JSON.stringify(this.caseInfo)),

        }
    }
}
</script>

<style>





</style>
