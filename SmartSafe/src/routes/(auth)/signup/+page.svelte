<script lang="ts">
    import { enhance } from "$app/forms";
    import sendEmailOTP from "$lib/api/sendEmailOTP";

    type FormData = {
        error?: string;
    };

    export let form: FormData | null = null;

    let error = form?.error || null;

    let success: any = null;

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
            result: { type: string; data?: FormData };
        }) => {
            error = null;
            success = null;

            if (result.type === "failure") {
                error = result.data?.error || null;
            } else if (result.type === "success") {
                success = "Sign up successful!";
                setTimeout(() => {
                    window.location.href = "/signin";
                }, 1000);
            }
        };
    }

    async function sendOTP() {
        error = null;
        success = null;
        const emailInput = document.querySelector(
            'input[name="email"]',
        ) as HTMLInputElement | null;
        const email = emailInput?.value;
        if (email) {
            try {
                const response = await sendEmailOTP({ email });
                if (response.success) {
                    success = "OTP sent successfully!";
                    disableButton = true;
                    countdown();
                    setTimeout(() => {
                        timer = false;
                        setTimerMessage("");
                    }, time * 1000);
                } else {
                    error = response.error || "Failed to send OTP.";
                }
            } catch (err) {
                error = "Failed to send OTP. Please try again.";
            }
        } else {
            error = "Please enter your email first.";
        }
    }
</script>

<main class="bg-blue-500">
    <section
        class="w-full container mx-auto h-screen flex items-center justify-center px-4"
    >
        <div
            class=" max-w-[1000px] flex flex-col gap-4 items-center bg-gradient-to-r from-blue-800 to-blue-600 p-10 rounded-2xl text-white shadow-lg w-full"
        >
            <form
                class="flex flex-col items-center w-80 space-y-4"
                method="POST"
                action="?/signup"
                use:enhance={handleEnhance}
            >
                <input
                    type="text"
                    name="username"
                    placeholder="Username (3 - 20 characters)"
                    class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
                    required
                />
                <input
                    type="email"
                    name="email"
                    placeholder="Email"
                    class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
                    required
                />
                <input
                    type="password"
                    name="password"
                    placeholder="Password (6 characters minimum)"
                    class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
                    required
                />
                <input
                    type="password"
                    name="confirmPassword"
                    placeholder="Confirm Password (6 characters minimum)"
                    class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
                    required
                />
                <div class="flex flex-col gap-2">
                    <div class="flex flex-row gap-2">
                        <input
                            type="one-time-code"
                            name="oneTimeCode"
                            placeholder="OTP (6 Digits)"
                            class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
                            required
                        />
                        <button
                            type="button"
                            class="bg-blue-500 hover:bg-blue-700 text-white rounded-md p-3 w-1/2 transition-all duration-300"
                            on:click={sendOTP}
                            disabled={disableButton}
                        >
                            {disableButton
                                ? `Resend OTP (${time})`
                                : "Send OTP"}
                        </button>
                    </div>
                </div>
                <p class="text-sm m-2">
                    <a
                        href="/signin"
                        class="flex items-center space-x-2 text-gray-200"
                        >Already have an account?</a
                    >
                </p>
                {#if error}
                    <p class="text-red-500">{error}</p>
                {/if}

                {#if success}
                    <p class="text-green-500">{success}</p>
                {/if}
                <button
                    type="submit"
                    class="bg-blue-500 hover:bg-blue-700 text-white rounded-md p-3 w-full transition-all duration-300"
                >
                    Sign Up
                </button>
            </form>
        </div>
    </section>
</main>
