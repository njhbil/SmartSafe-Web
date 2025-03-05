<script lang="ts">
    export let title: string = "Smart";
    export let subtitle: string = "Safe";
    export const onClick: () => void = () => {};
</script>

<header class="w-full bg-gray-900 text-white shadow-lg fixed top-0 left-0 z-50 mb-10">
    <div class="max-w-[1400px] mx-auto flex items-center justify-between p-4">
        <a href="/" class="text-2xl font-bold">
            <span class="text-gray-300">{title}</span><span class="text-blue-500"> {subtitle}</span>
        </a>

        <nav class="hidden md:flex items-center gap-6">
            <a href="#smartsafe" class="hover:text-blue-400 transition">Tentang Kita</a>
            <a href="#faq" class="hover:text-blue-400 transition">FAQ</a>
            <a href="/LOGIN" class="bg-blue-500 px-5 py-2 rounded-full hover:bg-blue-600 transition">Login</a>
        </nav>

        <button 
            aria-label="Open menu" 
            class="md:hidden text-black text-2xl">
            <i class="fa-solid fa-bars"></i>
        </button>
    </div>
</header>