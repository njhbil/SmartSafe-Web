<script lang="ts">
    export let socialMedia = [
        { name: "Facebook", icon: "facebook", url: "#" },
        { name: "Twitter", icon: "twitter", url: "#" },
        { name: "Instagram", icon: "instagram", url: "#" },
        { name: "YouTube", icon: "youtube", url: "#" },
    ];

    export let copyright = "© 2025 SMARTSAFE LTD. All rights reserved.";

    export let supportLinks = [
        { name: "Privacy Notice", url: "#" },
        { name: "Terms & Conditions", url: "#" },
        { name: "Cookie Settings", url: "#" },
    ];
</script>

<footer class="bg-gray-900 text-gray-300 py-10 px-6">
    <div
        class="max-w-7xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6"
    >
        <!-- Brand -->
        <div class="col-span-1">
            <h1 class="text-2xl font-bold text-white">SmartSafe</h1>
            <p>{copyright}</p>
        </div>

        <!-- Support Links -->
        <div>
            <h2 class="font-semibold text-white">Support Us</h2>
            <ul class="mt-2 space-y-2">
                {#each supportLinks as link}
                    <li>
                        <a href={link.url} class="hover:text-white"
                            >{link.name}</a
                        >
                    </li>
                {/each}
            </ul>
        </div>

        <!-- Social Media -->
        <div>
            <h2 class="font-semibold text-white">Connect with Us</h2>
            <div class="flex space-x-4 mt-2">
                {#each socialMedia as social}
                    <a
                        href={social.url}
                        aria-label={social.name}
                        class="hover:text-white text-lg"
                    >
                        <i class="fab fa-{social.icon}"></i>
                    </a>
                {/each}
            </div>
        </div>
    </div>

    <!-- Footer Bottom -->
    <div
        class="border-t border-gray-700 mt-8 pt-6 flex flex-col sm:flex-row items-center justify-between text-sm"
    ></div>
</footer>
