<script lang="ts">
    import loginAPI from "$lib/api/loginAPI";

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
        password: "",
        rememberMe: false,
    };

    const handleSubmit = async (values: {
        email: string;
        password: string;
        rememberMe: boolean;
    }) => {
        try {
            setErrors(null);
            setSuccess(null);
            const response = await loginAPI(values);

            if (values.email === "" || values.password === "") {
                setErrors("Please fill in all fields.");
                return;
            }

            if (response.success) {
                setSuccess("Login successful. Redirecting...");
                localStorage.setItem("token", response.token);
            } else {
                setErrors("Email or password is invalid.");
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
            <input
                type="password"
                name="password"
                placeholder="Password"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.password}
            />
            <p class="text-sm m-2">
                Don't have an account?
                <a href="/signup" class="text-sm text-blue-700">Sign Up</a>
            </p>
            <p class="text-sm m-2">
                Forgot your password?
                <a href="/forgot-password" class="text-sm text-blue-700"
                    >Reset Password</a
                >
            </p>
            <label>
                <input
                    type="checkbox"
                    name="rememberMe"
                    class="m-2"
                    bind:checked={values.rememberMe}
                />
                Remember Me
            </label>
            {#if error}
                <p class="text-red-500">{error}</p>
            {/if}
            {#if success}
                <p class="text-green-500">{success}</p>
            {/if}
            <button
                type="submit"
                class="bg-blue-500 text-white rounded p-2 m-2"
                on:click|preventDefault={() => handleSubmit(values)}
            >
                Sign Up
            </button>
        </form>
    </section>
</main>
