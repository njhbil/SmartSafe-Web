<script>
    import { fly, fade } from "svelte/transition";
    import Page3 from "./page3.svelte";
    let faqs = [
        {
            question: "Apa fungsi utama aplikasi ini?",
            answer: "Aplikasi ini membantu pengguna mengetahui wilayah rawan kejahatan berdasarkan data yang tersedia. Selain itu, aplikasi menyediakan fitur telepon darurat yang dapat langsung menghubungi pihak berwenang saat dalam keadaan darurat.",
            open: false,
        },
        {
            question:
                "Bagaimana cara mengetahui daerah rawan kejahatan di sekitar saya?",
            answer: "Anda dapat membuka peta dalam aplikasi yang menampilkan wilayah rawan kejahatan dengan indikator warna atau ikon tertentu. Data ini diperbarui secara berkala berdasarkan laporan pengguna dan sumber resmi.",
            open: false,
        },
        {
            question: "Bagaimana cara menggunakan fitur telepon darurat?",
            answer: "Fitur telepon darurat dapat diakses melalui tombol khusus dalam aplikasi. Saat ditekan, aplikasi akan langsung menghubungkan Anda dengan nomor darurat terdekat, seperti polisi atau layanan keamanan setempat.",
            open: false,
        },
        {
            question:
                " Apakah data wilayah rawan kejahatan di aplikasi ini akurat?",
            answer: "Akurasi data bergantung pada laporan pengguna serta data dari sumber resmi. Kami terus memperbarui informasi agar pengguna mendapatkan data terbaru mengenai tingkat kejahatan di suatu wilayah.",
            open: false,
        },
        {
            question: "Apakah aplikasi ini bisa digunakan secara offline?",
            answer: "Beberapa fitur, seperti peta wilayah rawan kejahatan, membutuhkan koneksi internet untuk memperbarui data. Namun, fitur telepon darurat tetap bisa digunakan selama perangkat memiliki sinyal seluler.",
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

