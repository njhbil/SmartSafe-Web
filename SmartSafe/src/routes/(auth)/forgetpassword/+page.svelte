<script lang="ts">
    import forgetAPI from "$lib/api/forgetAPI";

    let error: any = null;
    function setErrors(errors: any) {
        error = errors;
    }

    let success: any = null;
    function setSuccess(message: any) {
        success = message;
    }

    let values = {
        email: "",
    };

    const handleSubmit = async (value: { email: string }) => {
        try {
            setErrors(null);
            setSuccess(null);
            const response = await forgetAPI(value);

            if (value.email === "") {
                setErrors("Please fill in all fields.");
                return;
            }

            if (response.success) {
                setSuccess(response.message);
            } else {
                setErrors("Email is invalid.");
            }
        } catch (error) {
            console.error(error);
            setErrors("An error occurred, please try again later.");
        }
    };
</script>

<main>
    <section class="h-screen flex items-center justify-center">
        <form class="flex flex-col items-center">
            <input
                type="email"
                name="email"
                placeholder="Email"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.email}
            />
            {#if error}
                <p class="text-red-500">{error}</p>
            {/if}
            {#if success}
                <p class="text-green text-center">{success}</p>
            {/if}
            <button
                type="submit"
                class="bg-blue-500 text-white rounded p-2 m-2"
                on:click|preventDefault={() => handleSubmit(values)}
            >
                Send Email
            </button>
        </form>
    </section>
</main>
