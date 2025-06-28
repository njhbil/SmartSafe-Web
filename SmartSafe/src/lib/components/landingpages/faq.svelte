<script>
    import { fly, fade } from "svelte/transition";
    import Page3 from "./page3.svelte";
    let faqs = [
        {
            question: "Apa fungsi utama program ini?",
            answer: "Aplikasi ini membantu pengguna dalam mengidentifikasi wilayah yang rawan tindak kejahatan berdasarkan data yang tersedia. Selain itu, aplikasi ini juga menyediakan fitur panggilan darurat yang dapat langsung menghubungkan pengguna dengan pihak berwenang jika terjadi keadaan darurat.",
            open: false,
        },
        {
            question: "Bagaimana cara menemukan daerah rawan kejahatan di sekitar saya?",
            answer: "Kamu dapat membuka peta di dalam aplikasi yang menampilkan daerah rawan kejahatan dengan indikator warna atau ikon tertentu. Data ini diperbarui secara berkala berdasarkan laporan pengguna dan sumber resmi..",
            open: false,
        },
        {
            question: "Bagaimana cara saya menggunakan fitur panggilan darurat?",
            answer: "Fitur panggilan darurat dapat diakses melalui tombol khusus di dalam aplikasi. Saat ditekan, aplikasi akan langsung menghubungkan Anda ke nomor darurat terdekat, seperti kepolisian atau layanan keamanan setempat.",
            open: false,
        },
        {
            question: "Apakah data tentang daerah rawan kejahatan di aplikasi ini akurat?",
            answer: "TAkurasi data bergantung pada laporan pengguna dan sumber resmi. Kami terus memperbarui informasi agar pengguna mendapatkan data terbaru mengenai tingkat kejahatan di suatu wilayah.",
            open: false,
        },
        {
            question: "Apakah aplikasi ini dapat digunakan secara offline?",
            answer: "Beberapa fitur, seperti peta wilayah rawan kejahatan, memerlukan koneksi internet untuk memperbarui data. Namun, fitur panggilan darurat tetap dapat digunakan selama perangkat memiliki sinyal seluler..",
            open: false,
        },
    ];

    // @ts-ignore
    function toggleFAQ(index) {
        faqs = faqs.map((faq, i) => ({
            ...faq,
            open: i === index ? !faq.open : false,
        }));
    }
</script>

<Page3 id="faq">
    <div class="max-w-2xl mx-auto py-10">
        <h2
            class="text-4xl font-bold text-center mb-6"
            in:fade={{ duration: 600 }}
        >
            📌 FAQ
        </h2>

        <div class="space-y-4">
            {#each faqs as faq, i}
                <div
                    class="border border-gray-300 rounded-lg shadow-md transition-all duration-300 hover:shadow-lg"
                >
                    <button
                        class="w-full text-left p-4 flex justify-between items-center font-medium text-lg bg-gray-100 hover:bg-gray-200 transition duration-200 rounded-lg"
                        on:click={() => toggleFAQ(i)}
                    >
                        {faq.question}
                        <span
                            class="transform transition-transform duration-300"
                            class:rotate-180={faq.open}>▼</span
                        >
                    </button>

                    {#if faq.open}
                        <div
                            class="p-4 text-gray-700 bg-white border-t border-gray-300"
                            in:fly={{ y: 10, duration: 400 }}
                        >
                            {faq.answer}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</Page3>
