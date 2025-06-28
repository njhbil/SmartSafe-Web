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
        class="flex flex-col md:flex-row min-h-screen w-full"
    >
        <div class="md:w-1/2 bg-blue-500 text-white flex flex-col items-center justify-center p-10">
            <img src="/images/LAPORKAN.png" alt="Logo" class="h-24 w-auto mb-6" />
            <h1 class="text-3xl font-bold">LaporIn</h1>
            <p class="text-lg mt-2 text-center">Pantau,Waspada dan Tetap Aman</p>
            <p class="text-lg mt-2 text-center">Radar Keamanan di Genggamanmu</p>
        </div>

       <div class="md:w-1/2 bg-white flex items-center justify-center p-8 rounded-t-4xl md:rounded-tl-4xl md:rounded-bl-4xl md:rounded-tr-none md:rounded-br-none">
            <div class="w-full max-w-md">
            <h2 class="text-2xl font-bold text-gray-800 mb-6 text-center">Sign Up</h2>
            <form class="space-y-4" method="POST" action="?/signup" use:enhance={handleEnhance}>
                <input type="text" name="username" placeholder="Username" required class="w-full p-3 border border-gray-300 rounded" />
                <input type="email" name="email" placeholder="Email" required class="w-full p-3 border border-gray-300 rounded" />
                <input type="password" name="password" placeholder="Password" required class="w-full p-3 border border-gray-300 rounded" />
                <input type="password" name="confirmPassword" placeholder="Confirm Password" required class="w-full p-3 border border-gray-300 rounded" />
                
                <div class="flex gap-2">
                <input type="text" name="oneTimeCode" placeholder="OTP" class="w-full p-3 border border-gray-300 rounded" required />
                <button
                    type="button"
                    class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                    on:click={sendOTP}
                    disabled={disableButton}
                >
                    {disableButton ? `Request (${time})` : "Obtain"}
                </button>
                </div>

                {#if error}
                <p class="text-red-500 text-sm">{error}</p>
                {/if}
                {#if success}
                <p class="text-green-500 text-sm">{success}</p>
                {/if}

                <button type="submit" class="w-full bg-blue-500 text-white p-3 rounded hover:bg-blue-600 transition">
                Sign Up
                </button>

                <p class="text-sm text-center mt-4">
                Sudah punya akun? <a href="/signin" class="text-blue-600 font-semibold hover:underline">Sign In</a>
                </p>
            </form>
            </div>
        </div>
    </section>
</main>
