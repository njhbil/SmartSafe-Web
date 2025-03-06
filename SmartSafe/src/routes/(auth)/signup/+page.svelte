<script lang="ts">
    import registerApi from "$lib/api/registerAPI";

    let error: any = null;
    function setErrors(errors: any) {
        error = { ...error, ...errors };
    }

    let success: any = null;
    function setSuccess(message: any) {
        success = message;
    }

    const specialChars = new RegExp(/[!@#$%^&*(),.?":{}|<>]/);

    function validateEmail(email: string) {
        const emailcheck = /\S+@\S+\.\S+/;
        return emailcheck.test(email);
    }

    function validateInput(value: string, type: string) {
        if (type === "username") {
            if (!value.length) {
                setErrors({ username: "Username is required." });
                return;
            }

            if (value.length < 3) {
                setErrors({
                    username: "Username must be at least 3 characters.",
                });
                return;
            }

            if (value.includes(" ")) {
                setErrors({ username: "Username cannot contain spaces." });
                return;
            }

            if (specialChars.test(value)) {
                setErrors({
                    username: "Username cannot contain special characters.",
                });
                return;
            }

            if (value.length > 50) {
                setErrors({
                    username: "Username must be less than 50 characters.",
                });
                return;
            }
            setErrors({ username: null });
        }

        if (type === "email") {
            if (!value.length) {
                setErrors({ email: "Email is required." });
                return;
            }

            if (value.length > 200) {
                setErrors({ email: "Email must be less than 200 characters." });
                return;
            }

            if (!validateEmail(value)) {
                setErrors({ email: "Invalid email address." });
                return;
            }
            setErrors({ email: null });
        }

        if (type === "password") {
            if (!value.length) {
                setErrors({ password: "Password is required." });
                return;
            }

            if (value.length < 8) {
                setErrors({
                    password: "Password must be at least 8 characters.",
                });
                return;
            }

            if (value.includes(" ")) {
                setErrors({ password: "Password cannot contain spaces." });
                return;
            }

            if (value.length > 200) {
                setErrors({
                    password: "Password must be less than 200 characters.",
                });
                return;
            }
            setErrors({ password: null });
        }

        if (type === "passwordConfirmation") {
            if (!value.length) {
                setErrors({
                    confirmPassword: "Password confirmation is required.",
                });
                return;
            }

            if (value !== values.password) {
                setErrors({ confirmPassword: "Passwords do not match." });
                return;
            }
            setErrors({ confirmPassword: null });
        }
    }

    let values = {
        username: "",
        email: "",
        password: "",
        passwordConfirmation: "",
    };

    const handleSubmit = async (values: {
        username: string;
        email: string;
        password: string;
        passwordConfirmation: string;
    }) => {
        try {
            setErrors(null);
            setSuccess(null);
            validateInput(values.username, "username");
            validateInput(values.email, "email");
            validateInput(values.password, "password");
            validateInput(values.passwordConfirmation, "passwordConfirmation");

            const response = await registerApi(values);
            if (!response.success) {
                setErrors(response.message);
            } else {
                setSuccess(
                    "Account created successfully, redirecting to login...",
                );
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
                type="text"
                name="username"
                placeholder="Username"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.username}
            />
            {#if error && error.username}
                <p class="text-red-500">{error.username}</p>
            {/if}
            <input
                type="email"
                name="email"
                placeholder="Email"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.email}
            />
            {#if error && error.email}
                <p class="text-red-500">{error.email}</p>
            {/if}
            <input
                type="password"
                name="password"
                placeholder="Password"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.password}
            />
            {#if error && error.password}
                <p class="text-red-500">{error.password}</p>
            {/if}
            <input
                type="password"
                name="passwordConfirmation"
                placeholder="Confirm Password"
                class="border border-gray-300 rounded p-2 m-2"
                bind:value={values.passwordConfirmation}
            />
            <p class="text-sm m-2">
                <a href="/signin" class="text-blue-500"
                    >Already have an account?</a
                >
            </p>
            {#if error && error.confirmPassword}
                <p class="text-red-500">{error.confirmPassword}</p>
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
