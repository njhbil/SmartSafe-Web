<script lang="ts">
    import { onMount } from "svelte";
    import { initMap } from "$lib/components/maps";
    let mapElement: HTMLElement;

    onMount(() => {
        initMap(mapElement);
    });
</script>

<main>
    <section class="h-screen">
        <div id="maps" class="h-full" bind:this={mapElement}></div>
    </section>
</main>
