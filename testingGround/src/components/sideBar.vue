<template>
    <div class="w-[400px] h-full bg-gray-200 text-white text-center">
        <div class="h-[50px] bg-gray-900 flex justify-start items-center">
           <div class="px-[20px]">
                <h3 class="font-bold text-2xl">
                    WIRA Dashboard
                </h3>
           </div>
        </div>
        <div class="h-[calc(100vh-50px)] bg-gray-800 py-[20px]">
            <div class="px-[20px] flex flex-col justify-between space-y-[10px]">
                <div class="h-full">
                    <div class="h-[calc(100vh-50px)]">
                        <router-link to="/home" class="inline-flex relative items-center py-[10px] px-[10px] w-full text-sm font-medium rounded-md border-gray-200 hover:bg-gray-100 hover:text-gray-300">
                            Home
                        </router-link>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
</script>

<style>
</style>