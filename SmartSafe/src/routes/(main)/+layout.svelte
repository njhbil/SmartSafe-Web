<script lang="ts">
    import { goto } from "$app/navigation";
</script>

<slot />
<nav class="bg-blue-500 text-white text-center fixed bottom-0 w-full">
    <div class="flex">
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/3 rounded text-white font-bold py-4"
            on:click={() => goto("/maps")}
        >
            Maps
        </button>
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/3 rounded text-white font-bold py-4"
            on:click={() => goto("/notifications")}
        >
            Notifications
        </button>
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/3 rounded text-white font-bold py-4"
            on:click={() => goto("/profile")}
        >
            Profile
        </button>
    </div>
</nav>
