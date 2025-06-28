<script lang="ts">
    import { goto } from "$app/navigation";
</script>

<slot />
<nav
    class="bg-blue-500 text-white h-10 text-center fixed bottom-0 w-full overflow-hidden"
>
    <div class="flex">
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/3 text-white font-bold py-4"
            aria-label="Maps"
            on:click={() => goto("/maps")}
        >
            <i class="fa-solid fa-house"></i>
        </button>
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/3 text-white font-bold py-4"
            aria-label="Notifications"
            on:click={() => goto("/notifications")}
        >
            <i class="fa-solid fa-bell"></i>
        </button>
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/3 text-white font-bold py-4"
            aria-label="Statistik"
        >
            <i class="fa-solid fa-chart-simple"></i>
        </button>
        <button class="hover:bg-blue-300 bg-blue-500 w-1/3 text-white font-bold py-4" aria-label="Profile" on:click={() => goto("/profil")}>
            <i class="fa-solid fa-user"></i>
        </button>
    </div>
</nav>
