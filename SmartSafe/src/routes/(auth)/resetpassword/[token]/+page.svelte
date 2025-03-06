<script lang="ts">
    import { onMount } from "svelte";
    import resetAPI from "$lib/api/resetAPI";

    let error: any = null;
    function setErrors(errors: any) {
        error = errors;
    }

    let success: any = null;
    function setSuccess(message: any) {
        success = message;
    }

    let values = {
        password: "",
        confirmPassword: "",
        resetPWDtoken: "",
    };

    onMount(() => {
        if (typeof window !== "undefined") {
            const url = new URL(window.location.href);
            const pathSegments = url.pathname.split("/");
            values.resetPWDtoken = pathSegments[pathSegments.length - 1];
        }
    });

    const handleSubmit = async (values: {
        password: string;
        confirmPassword: string;
        resetPWDtoken: string;
    }) => {
        try {
            setErrors(null);
            setSuccess(null);
            const response = await resetAPI(values);

            if (values.password === "" || values.confirmPassword === "") {
                setErrors("Please fill in all fields.");
                return;
            }

            if (values.password !== values.confirmPassword) {
                setErrors("Passwords do not match.");
                return;
            }

            if (response.success) {
                setSuccess(response.message);
            } else {
                setErrors("Reset failed.");
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
                type="password"
                name="password"
                placeholder="Password"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.password}
            />
            <input
                type="password"
                name="confirmPassword"
                placeholder="Confirm Password"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.confirmPassword}
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
                Reset Password
            </button>
        </form>
    </section>
</main>
