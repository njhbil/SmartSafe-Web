import { fail } from "@sveltejs/kit";
import registerAPI from "$lib/api/registerAPI";
import verifyEmailOTP from "$lib/api/verifyEmailOTP";
import type { Actions } from "../signin/$types";

export const actions: Actions = {
  signup: async ({ request }) => {
    try {
      const formData = await request.formData();
      const username = formData.get("username") as string;
      const email = formData.get("email") as string;
      const password = formData.get("password") as string;
      const confirmPassword = formData.get("confirmPassword") as string;
      const oneTimeCode = formData.get("oneTimeCode") as string;

      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      const specialChars = /[!@#$%^&*(),.?":{}|<>]/;

      if (!username || !email || !password || !confirmPassword) {
        return fail(400, { error: "Please fill out all fields." });
      }

      if (username.length < 3 || username.length > 20) {
        return fail(400, {
          error: "Username must be between 3 and 20 characters.",
        });
      }

      if (specialChars.test(username)) {
        return fail(400, {
          error: "Username cannot contain special characters.",
        });
      }

      if (!emailRegex.test(email)) {
        return fail(400, { error: "Please enter a valid email address." });
      }

      if (password.length < 6) {
        return fail(400, { error: "Password must be at least 6 characters." });
      }

      if (password !== confirmPassword) {
        return fail(400, { error: "Passwords do not match." });
      }

      if (!oneTimeCode) {
        return fail(400, { error: "Please enter the one-time code." });
      }

      const verifyResponse = await verifyEmailOTP({ email, oneTimeCode });

      if (!verifyResponse.success) {
        return fail(400, { error: "Invalid one-time code." });
      }

      const response = await registerAPI({ username, email, password });

      if (response.success) {
        return {
          status: 302,
          headers: { Location: "/signin" },
        };
      }
    } catch (error) {
      console.error("Register Error:", error);
      return fail(500, { error: "An error occurred, please try again later." });
    }
  },
};
