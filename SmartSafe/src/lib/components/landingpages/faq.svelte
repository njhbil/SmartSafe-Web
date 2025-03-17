<script>
    import { fly, fade } from "svelte/transition";
    import Page3 from "./page3.svelte";
    let faqs = [
        {
            question: "What is the main function of this application?",
            answer: "This application helps users identify crime-prone areas based on available data. Additionally, the app provides an emergency call feature that can directly connect to authorities in case of an emergency.",
            open: false,
        },
        {
            question: "How can I find crime-prone areas around me?",
            answer: "You can open the map in the app which displays crime-prone areas with specific color indicators or icons. This data is regularly updated based on user reports and official sources.",
            open: false,
        },
        {
            question: "How do I use the emergency call feature?",
            answer: "The emergency call feature can be accessed through a special button in the app. When pressed, the app will directly connect you to the nearest emergency number, such as the police or local security services.",
            open: false,
        },
        {
            question: "Is the data on crime-prone areas in this app accurate?",
            answer: "The accuracy of the data depends on user reports and official sources. We continuously update the information to ensure users have the latest data on crime levels in an area.",
            open: false,
        },
        {
            question: "Can this app be used offline?",
            answer: "Some features, such as the crime-prone area map, require an internet connection to update data. However, the emergency call feature can still be used as long as the device has a cellular signal.",
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
