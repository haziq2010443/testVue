<template>
    <div class="w-screen h-screen flex">
        <sideBar />
        <dashboardHeader />
    </div>
</template>

<script>
    import sideBar from '../../components/sideBar.vue'
    import dashboardHeader from '../../components/Header.vue'

    export default{
        components:{
            sideBar,
            dashboardHeader
        }
    }
</script>

<style>
</style>