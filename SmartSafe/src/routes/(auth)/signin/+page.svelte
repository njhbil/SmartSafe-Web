<script lang="ts">
  import { enhance } from "$app/forms";
  let showpassword =  false

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
        setTimeout(() => {
          window.location.href = "/maps";
        }, 1000);
      }
    };
  }
</script>

<main class="bg-blue-500">
  <section
        class="flex flex-col md:flex-row min-h-screen w-full"
    >
        <div class="md:w-1/2 bg-blue-500 text-white flex flex-col items-center justify-center p-10">
            <img src="/images/LAPORKAN.png" alt="Logo" class="h-24 w-auto " />
            <h1 class="text-3xl font-bold">LaporIn</h1>
            <p class="text-lg mt-2 text-center">Pantau,Waspada dan Tetap Aman</p>
            <p class="text-lg mt-2 text-center">Radar Keamanan di Genggamanmu</p>
        </div>
    <div class="md:w-1/2 bg-white flex items-center justify-center p-8 rounded-t-4xl md:rounded-tl-4xl md:rounded-bl-4xl md:rounded-tr-none md:rounded-br-none">

    <div class="w-full max-w-md">
      <form
        class="flex flex-col items-center w-80 space-y-6"
        method="POST"
        action="?/signin"
        use:enhance={handleEnhance}
      >
        <p class="text-xl text-gray-700 font-semibold">
            SIGN IN
        </p>
        <input
          type="email"
          name="email"
          placeholder="Email"
          class="border border-gray-700 rounded-sm p-3 w-full text-gray-400 focus:ring-2 focus:ring-blue-500"
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          class="border border-gray-700 rounded-sm p-3 w-full text-gray-400 focus:ring-2 focus:ring-blue-500"
          required
        />
        <label class="flex items-center space-x-2 text-gray-700">
          <span>Ingat Saya</span>
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
          class="bg-blue-500 hover:bg-blue-700 text-white p-3 rounded-xl w-full transition-all duration-300"
        >
          Sign In
        </button>
        <p class="text-sm text-gray-700">
          Belum punya akun?
          <a
            href="/signup"
            class="text-sm text-blue-600 font-semibold hover:underline">Sign Up</a
          >
        </p>
        <p class="text-sm text-gray-700">
          Lupa password?
          <a
            href="/changepassword"
            class="text-sm text-blue-600 font-semibold hover:underline"
          >
            Reset Password</a
          >
        </p>
      </form>
      </div>
    </div>
  </section>
</main>
