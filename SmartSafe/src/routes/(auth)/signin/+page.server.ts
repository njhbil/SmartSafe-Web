import { fail } from "@sveltejs/kit";
import type { Actions } from "./$types";
import loginAPI from "$lib/api/loginAPI";

export const actions: Actions = {
  signin: async ({ request, cookies }) => {
    try {
      const formData = await request.formData();
      const email = formData.get("email") as string;
      const password = formData.get("password") as string;
      const rememberMe = formData.has("rememberMe");

      if (!email || !password) {
        return fail(400, { error: "Please fill in all fields." });
      }

      const response = await loginAPI({ email, password, rememberMe });

      if (response.success) {
        if (rememberMe) {
          cookies.set("refreshToken", response.refreshToken, {
            path: "/",
            httpOnly: true,
            secure: true,
            sameSite: "lax",
            maxAge: 60 * 60 * 24 * 30,
          });
        }
        if (!rememberMe) {
          cookies.set("refreshToken", response.refreshToken, {
            path: "/",
            httpOnly: true,
            secure: true,
            sameSite: "lax",
            maxAge: 60 * 60 * 24,
          });
        }

        return {
          status: 302,
          headers: { Location: "/maps" },
        };
      }

      if (!response.success) {
        return fail(400, { error: "Email or password is invalid!" });
      }
    } catch (error) {
      console.error("Login Error:", error);
      return fail(500, { error: "An error occurred, please try again later." });
    }
  },
};
