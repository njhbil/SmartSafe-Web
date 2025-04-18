<script lang="ts">
    import { enhance } from "$app/forms";

    type FormData = {
        error?: string;
        success?: string;
    };

    export let form: FormData | null = null;
    let error = form?.error || null;
    let success = form?.success || null;
    let isSubmitting = false;

    let timerMessage: any = null;
    function setTimerMessage(message: any) {
        timerMessage = message;
    }

    let timer: boolean = true;
    let time = 30;
    let disableButton = false;

    function countdown() {
        if (time > 0) {
            setTimerMessage(
                `Email sent successfully! You can try again in ${time} seconds.`,
            );
            time--;
            setTimeout(countdown, 1000);
        } else {
            timer = false;
            disableButton = false;
            setTimerMessage("");
        }
    }

    function handleEnhance() {
        return async ({
            result,
        }: {
            result: {
                type: string;
                data?: { error?: string; success?: string };
            };
        }) => {
            isSubmitting = false;
            error = null;
            success = null;

            if (result.type === "failure") {
                error = result.data?.error || null;
            } else if (result.type === "success") {
                success = result.data?.success || null;
                disableButton = true;
                countdown();
                setTimeout(() => {
                    timer = false;
                    setTimerMessage("");
                }, time * 1000);
            }
        };
    }
</script>

<main>
    <section class="h-screen flex items-center justify-center">
        <form
            class="flex flex-col items-center"
            method="POST"
            action="?/forgetpassword"
            use:enhance={({}) => {
                isSubmitting = true;

                return handleEnhance();
            }}
        >
            <input
                type="email"
                name="email"
                placeholder="Email"
                class="border border-gray-300 rounded p-2 m-2"
                required
            />
            {#if error}
                <p class="text-red-500 text-center">{error}</p>
            {/if}
            {#if success}
                <p class="text-green-500 text-center">{success}</p>
            {/if}
            <p class="text-black text-center">{timerMessage}</p>
            <button
                type="submit"
                class="bg-blue-500 text-white rounded p-2 m-2"
                disabled={isSubmitting || disableButton}
            >
                {#if isSubmitting}
                    Sending email...
                {:else}
                    Send Email
                {/if}
            </button>
        </form>
    </section>
</main>
