<script lang="ts">
  import { enhance } from "$app/forms";

  type FormData = {
    error?: string;
  };

  export let form: FormData | null = null;
  let error = form?.error || null;
  let success: string | null = null;

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
        success = "Login successful!";
        return setTimeout(() => {
          window.location.href = "/maps";
        }, 1000);
      }
    };
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
        action="?/signin"
        use:enhance={handleEnhance}
      >
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
          placeholder="Password"
          class="border border-gray-300 rounded-md p-3 w-full text-white focus:ring-2 focus:ring-blue-500"
          required
        />
        <p class="text-sm text-gray-200">
          Don't have an account?
          <a
            href="/signup"
            class="text-sm text-white font-semibold hover:underline">Sign Up</a
          >
        </p>
        <p class="text-sm text-gray-200">
          Forgot your password?
          <a
            href="/forgetpassword"
            class="text-sm text-white font-semibold hover:underline"
          >
            Reset Password</a
          >
        </p>
        <label class="flex items-center space-x-2 text-gray-200">
          <span>Remember Me</span>
          <input type="checkbox" name="rememberMe" class="accent-blue-500" />
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
        >
          Sign In
        </button>
      </form>
    </div>
  </section>
</main>
