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

<main class="bg-blue-500">
  <section class="w-full container mx-auto h-screen flex items-center justify-center px-4">
    <div class=" max-w-[1000px] flex flex-col gap-4 items-center bg-gradient-to-r from-blue-800 to-blue-600 p-10 rounded-2xl text-white shadow-lg w-full">
      <form class="flex flex-col items-center w-80 space-y-4">
        <input 
          type="email"
          name="email"
          placeholder="Email"
          class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
          bind:value={values.email}
        />
        <input 
          type="password"
          name="password"
          placeholder="Password"
          class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
          bind:value={values.password}
        />
        <p class="text-sm text-gray-200">
          Don't have an account?
          <a href="/signup" class="text-sm text-white font-semibold hover:underline">Sign Up</a>
        </p>
        <p class="text-sm text-gray-200">
          Forgot your password?
          <a href="/forgotpassword" class="text-sm text-white font-semibold hover:underline"> Reset Password</a>
        </p>
        <label class="flex items-center space-x-2 text-gray-200">
          <input 
            type="checkbox" 
            name="rememberMe" 
            class="accent-blue-500"
            bind:checked={values.rememberMe}
          />
          <span>Remember Me</span>
        </label>
        {#if error}
          <p class="text-red-500 text-sm">{error}</p>
        {/if}
        {#if success}
          <p class="text-green-500 text-sm">{success}</p>
        {/if}
        <button 
          type="submit"
          class="bg-blue-500 hover:bg-blue-700 text-white rounded-md p-3 w-full transition-all duration-300"
          on:click|preventDefault={() => handleSubmit(values)}
        >
          Sign Up
        </button>
      </form>
    </div>
  </section>
</main>