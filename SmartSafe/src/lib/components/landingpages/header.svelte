<script lang="ts">
    export let title: string = "Lapor";
    export let subtitle: string = "In";
    import handleHistory from "../HandleHistory";

    let isMobileMenuOpen: boolean = false;
    const toggleMobileMenu = () => {
        isMobileMenuOpen = !isMobileMenuOpen;
    };
</script>

<header
    class="w-full bg-gray-900 text-white shadow-lg fixed top-0 left-0 z-50 mb-10"
>
    <div class="max-w-[1400px] mx-auto flex items-center justify-between p-4">
        <a href="/" class="text-2xl font-bold flex items-center ">
            <img src="/images/LAPORKAN.png" alt="Logo" class="h-12 w-auto"><span class="text-gray-300">{title}</span><span
                class="text-blue-500"
            >
                {subtitle}</span
            >
        </a>

        <nav class="hidden md:flex items-center gap-6">
            <a
                href="#smartsafe"
                class="hover:text-blue-400 transition font-bold"
                on:click={(e) => handleHistory(e, "smartsafe")}>Tentang Kita</a
            >
            <a
                href="#faq"
                class="hover:text-blue-400 transition font-bold"
                on:click={(e) => handleHistory(e, "faq")}>FAQ</a
            >
            <a
                href="/signin"
                class="bg-blue-500 px-5 py-2 rounded-full hover:bg-blue-600 transition"
                >Sign In</a
            >
        </nav>

        <button
            aria-label="Open menu"
            class="md:hidden text-white text-2xl"
            on:click={toggleMobileMenu}
        >
            <i class="fa-solid fa-bars"></i>
        </button>
    </div>

    {#if isMobileMenuOpen}
        <div class="md:hidden bg-gray-800 text-white p-4 flex flex-col gap-4">
            <a
                href="#smartsafe"
                class="block py-2 hover:text-blue-400 transition"
                on:click={(e) => handleHistory(e, "smartsafe")}>Tentang Kita</a
            >
            <a
                href="#faq"
                class="block py-2 hover:text-blue-400 transition"
                on:click={(e) => handleHistory(e, "faq")}>FAQ</a
            >
            <a
                href="/signin"
                class="block py-2 w-1/2 bg-blue-500 text-center text-white rounded-full hover:bg-blue-600 transition"
                >Sign In</a
            >
        </div>
    {/if}
</header>
