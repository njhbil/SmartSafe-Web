<script lang="ts">
    import { enhance } from "$app/forms";
    import { page } from "$app/stores";

    type FormData = {
        error?: string;
    };

    export let form: FormData | null = null;
    let error = form?.error || null;
    let success: any = null;

    const token = $page.url.searchParams.get("token") as string;

    function handleEnhance() {
        return async ({
            result,
        }: {
            result: { type: string; data?: FormData };
        }) => {
            error = null;
            success = null;

            if (result.type === "failure") {
                error = result.data?.error || null;
            } else if (result.type === "success") {
                success = "Reset password successful!";
                setTimeout(() => {
                    window.location.href = "/signin";
                }, 1000);
            }
        };
    }
</script>

<main>
    <section class="h-screen flex items-center justify-center">
        <form
            class="flex flex-col items-center"
            method="POST"
            action="?/resetpassword"
            use:enhance={handleEnhance}
        >
            <input type="hidden" name="token" value={token} />
            <input
                type="password"
                name="password"
                placeholder="Password"
                class="border border-gray-300 rounded p-2 m-2"
                required
            />
            <input
                type="password"
                name="confirmPassword"
                placeholder="Confirm Password"
                class="border border-gray-300 rounded p-2 m-2"
                required
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
            >
                Reset Password
            </button>
        </form>
    </section>
</main>
