<script>
    export let id;
</script>


<section {id} class={"min-h-screen p-10 pt-20 flex flex-col px-4 bg-blue-500 rounded-tl-[50px] rounded-tr-[50px]" }>
    <div class="flex flex-col flex-1 max-w-[1400px] mx-auto w-full">
        <slot></slot>
</section>