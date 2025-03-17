<script lang="ts">
    import { goto } from "$app/navigation";
</script>

<slot />
<nav class="bg-blue-500 text-white text-center fixed bottom-0 w-full">
    <div class="flex">
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/2 rounded text-white font-bold py-4"
            on:click={() => goto("/signup")}
        >
            Sign Up
        </button>
        <button
            class="hover:bg-blue-300 bg-blue-500 w-1/2 rounded text-white font-bold py-4"
            on:click={() => goto("/signin")}
        >
            Sign In
        </button>
    </div>
</nav>
