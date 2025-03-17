import { fail } from "@sveltejs/kit";
import resetAPI from "$lib/api/resetAPI";
import type { Actions } from "./$types";

export const actions: Actions = {
  resetpassword: async ({ request }) => {
    try {
      const formData = await request.formData();
      const password = formData.get("password") as string;
      const confirmPassword = formData.get("confirmPassword") as string;
      const resetPWDtoken = formData.get("token") as string;

      if (!password || !confirmPassword) {
        return fail(400, { error: "Please fill out all fields." });
      }

      if (password.length < 6) {
        return fail(400, { error: "Password must be at least 6 characters." });
      }

      if (password !== confirmPassword) {
        return fail(400, { error: "Passwords do not match." });
      }

      const response = await resetAPI({ resetPWDtoken, password });

      if (response.success) {
        return {
          status: 302,
          headers: { Location: "/signin" },
        };
      }

      return fail(400, { error: response.message });
    } catch (error) {
      console.error("Reset Password Error:", error);
      return fail(500, { error: "An error occurred, please try again later." });
    }
  },
};
