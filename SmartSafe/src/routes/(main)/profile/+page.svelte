<script lang="ts">
    import DefaultProfile from "$lib/assets/profile.jpg";
    import { enhance } from "$app/forms";

    function handleEnhance() {
        return async ({ result }: { result: { type: string; data?: any } }) => {
            if (result.type === "failure") {
                console.error(result.data);
            } else if (result.type === "success") {
                console.log("Form submitted successfully");
                setTimeout(() => {
                    window.location.href = "/signin";
                }, 1000);
            }
        };
    }
</script>

<main>
    <div class="flex flex-col min-h-screen px-4 py-12">
        <div
            class="w-full bg-white shadow-lg rounded-lg overflow-hidden p-8 transition-all duration-300 hover:shadow-xl"
        >
            <div class="flex items-center justify-center text-center">
                <img
                    src={DefaultProfile}
                    alt="PlaceholderImage"
                    class="rounded-full w-24 h-24 mb-4"
                />
            </div>
            <h1 class="text-sm text-center font-bold text-gray-800">Name</h1>
            <div class=" h-1 bg-blue-500 my-3 rounded"></div>
        </div>
        <div
            class="flex flex-col w-full my-2 bg-white shadow-lg rounded-lg overflow-hidden p-8 transition-all duration-300 hover:shadow-xl"
        >
            <div class="flex flex-row items-center my-1">
                <i
                    class="fa-solid fa-address-card h-4 w-4 mx-2"
                    style="color: #000000;"
                ></i>
                <h1 class="text-sm text-center text-gray-800 mx-2">
                    Change Account Information
                </h1>
            </div>
            <div class="flex flex-row items-center my-1">
                <i
                    class="fa-solid fa-language mx-2 h-4 w-4"
                    style="color: #000000;"
                ></i>
                <h1 class="text-sm text-center text-gray-800 mx-2">Language</h1>
            </div>
            <div class="flex flex-row items-center my-1">
                <i class="fa-solid fa-gear mx-2 h-4 w-4" style="color: #000000;"
                ></i>
                <h1 class="text-sm text-center text-gray-800 mx-2">Settings</h1>
            </div>
        </div>
        <div
            class="w-full bg-white shadow-lg rounded-lg overflow-hidden p-8 transition-all duration-300 hover:shadow-xl"
        >
            <div class="flex flex-row items-center my-1">
                <form
                    action="?/logout"
                    method="POST"
                    use:enhance={handleEnhance}
                >
                    <button
                        type="submit"
                        class="flex flex-row items-center my-1"
                    >
                        <i
                            class="fa-solid fa-right-from-bracket mx-2 h-4 w-4"
                            style="color: #000000;"
                        ></i>
                        <span class="text-sm text-center text-gray-800 mx-2"
                            >Logout</span
                        >
                    </button>
                </form>
            </div>
        </div>
    </div>
</main>
