import { fail } from "@sveltejs/kit";
import forgetAPI from "$lib/api/forgetAPI";
import type { Actions } from "./$types";

export const actions: Actions = {
  forgetpassword: async ({ request }) => {
    try {
      const formData = await request.formData();
      const email = formData.get("email") as string;

      if (!email) {
        return fail(400, { error: "Please fill out all fields." });
      }

      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

      if (!emailRegex.test(email)) {
        return fail(400, { error: "Please enter a valid email address." });
      }

      const response = await forgetAPI({ email });

      if (response.success) {
        return {
          status: 302,
          message: "Email sent successfully",
        };
      }

      return fail(400, { error: response.message });
    } catch (error) {
      console.error("Forgot Password Error:", error);
      return fail(500, { error: "An error occurred, please try again later." });
    }
  },
};
