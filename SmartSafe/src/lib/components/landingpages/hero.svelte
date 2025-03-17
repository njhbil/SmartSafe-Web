<script lang="ts">
    import Header from "./header.svelte";
    import Page1 from "./page1.svelte";
    import { fade, fly } from "svelte/transition";
    import { onMount } from "svelte";
    export let onClick: () => void;

    let isVisible = false;
    onMount(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                isVisible = entry.isIntersecting;
            },
            { threshold: 0.3 },
        );

        const section = document.querySelector("#hero-section");
        if (section) observer.observe(section);

        return () => observer.disconnect();
    });
</script>

;
<Page1 id>
    <Header />
    <div
        id="hero-section"
        class="flex flex-col gap-6 items-center justify-center text-center py-20 pt-32"
    >
        {#if isVisible}
            <h2
                class="text-4xl sm:text-5xl md:text-6xl font-extrabold leading-tight max-w-[900px]"
                in:fade={{ duration: 800 }}
            >
                <span
                    class="bg-gradient-to-r from-blue-500 to-indigo-600 text-transparent bg-clip-text"
                >
                    Monitor, Be Alert, Stay Safe
                </span>
                <br />
                Security Radar in Your
                <span
                    class="bg-gradient-to-r from-indigo-600 to-purple-500 text-transparent bg-clip-text"
                >
                    Hands
                </span>
            </h2>

            <button
                on:click={onClick}
                class="SpesialBTNDark"
                in:fly={{ y: 20, duration: 500 }}
            >
                <a
                    href="/maps"
                    class="text-base sm:text-lg md:text-xl font-semibold"
                    >Get Started &rarr;</a
                >
            </button>
        {/if}
    </div>
</Page1>
