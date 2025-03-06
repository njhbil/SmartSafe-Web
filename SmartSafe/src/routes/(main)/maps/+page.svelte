<script lang="ts">
    import { onMount } from "svelte";
    import { initMap } from "$lib/components/maps/maps";
    let mapElement: HTMLElement;

    function updateViewPortHeight() {
        const vh = window.innerHeight * 0.01;
        document.documentElement.style.setProperty("--vh", `${vh}px`);
    }

    onMount(() => {
        initMap(mapElement);
        updateViewPortHeight();
        window.addEventListener("resize", updateViewPortHeight);
        window.addEventListener("orientationchange", updateViewPortHeight);

        return () => {
            window.removeEventListener("resize", updateViewPortHeight);
            window.removeEventListener(
                "orientationchange",
                updateViewPortHeight,
            );
        };
    });
</script>

<main class="map-container absolute top-0 w-full overflow-hidden">
    <section class="map-container absolute w-full">
        <div id="maps" class="map-container" bind:this={mapElement}></div>
    </section>
</main>
