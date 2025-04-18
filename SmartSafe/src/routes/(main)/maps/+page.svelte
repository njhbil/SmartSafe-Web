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

        mapElement.addEventListener("touchmove", (event) => {
            event.preventDefault();
        });

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
    <section class="map-wrapper absolute w-full overflow-hidden">
        <div id="maps" class="map-element" bind:this={mapElement}></div>
    </section>
</main>
