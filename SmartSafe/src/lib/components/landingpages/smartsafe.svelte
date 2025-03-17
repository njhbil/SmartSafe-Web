<script>
    import Page2 from "./page2.svelte";
    import { fade, scale } from "svelte/transition";
    import { onMount } from "svelte";

    let isVisible = false;

    onMount(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                isVisible = entry.isIntersecting;
            },
            { threshold: 0.3 },
        );

        const section = document.querySelector("#animated-section");
        if (section) observer.observe(section);

        return () => observer.disconnect();
    });
</script>

<Page2 id="smartsafe">
    <div class="w-full container mx-auto px-4 md:px-8">
        <div
            id="animated-section"
            class="flex flex-col gap-6 items-center bg-gradient-to-r from-blue-500 to-blue-700 text-white p-10 rounded-2xl shadow-lg"
        >
            {#if isVisible}
                <h3
                    class="text-4xl sm:text-5xl md:text-6xl font-bold"
                    in:fade={{ delay: 200, duration: 600 }}
                >
                    🚀 SMART SAFE
                </h3>
                <p
                    class="text-lg max-w-2xl font-medium text-justify leading-relaxed"
                    in:fade={{ delay: 400, duration: 800 }}
                >
                    SmartSafe is a web-based application designed to provide
                    information about crime-prone areas. Additionally, SmartSafe
                    is equipped with an <span class="font-bold text-yellow-300"
                        >emergency call</span
                    >
                    and <span class="font-bold text-yellow-300">SOS</span> feature
                    that can be used in case of crime-related issues. This application
                    is created as a modern solution to tackle the increasing rate
                    of crime.
                </p>
                <button
                    class="bg-yellow-300 text-black font-semibold py-2 px-6 rounded-lg shadow-md hover:bg-yellow-400 transition duration-300 transform hover:scale-110"
                    in:scale={{ start: 0.8, duration: 500, delay: 600 }}
                >
                    Learn More
                </button>
            {/if}
        </div>
    </div>
</Page2>
