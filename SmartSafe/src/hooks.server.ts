import type { Handle } from "@sveltejs/kit";
import refreshTokenAPI from "$lib/api/refreshTokenAPI";
import verifyTokenAPI from "$lib/api/verifyTokenAPI";

export const handle: Handle = async ({ event, resolve }) => {
  const refreshToken = event.cookies.get("refreshToken") as string;
  const loginToken = event.cookies.get("loginToken") as string;

  const publicRoutes = [
    "/",
    "/signup",
    "/signin",
    "/forgetpassword",
    "/resetpassword",
  ];

  if (publicRoutes.includes(event.url.pathname)) {
    if (loginToken) {
      try {
        const response = await verifyTokenAPI({ loginToken, refreshToken });
        if (response.success) {
          if (
            event.url.pathname === "/signin" ||
            event.url.pathname === "/signup"
          ) {
            return new Response(null, {
              status: 302,
              headers: { Location: "/maps" },
            });
          }
        }
      } catch (error) {
        console.error("Invalid token:", error);
      }
    }
    return resolve(event);
  }

  if (!loginToken) {
    try {
      const response = await refreshTokenAPI({ refreshToken });
      if (response.success) {
        event.cookies.set("loginToken", response.refreshToken, {
          path: "/",
          httpOnly: true,
          secure: true,
          sameSite: "lax",
          maxAge: 60 * 60 * 24,
        });

        return resolve(event);
      } else {
        event.cookies.delete("refreshToken", { path: "/" });
        event.cookies.delete("loginToken", { path: "/" });

        console.error("Refresh token failed:", response.message);
        return new Response(null, {
          status: 302,
          headers: { Location: "/signin" },
        });
      }
    } catch (error) {
      event.cookies.delete("refreshToken", { path: "/" });
      event.cookies.delete("loginToken", { path: "/" });

      console.error("Refresh token error:", error);
      return new Response(null, {
        status: 302,
        headers: { Location: "/signin" },
      });
    }
  }

  try {
    const response = await verifyTokenAPI({ loginToken, refreshToken });
    if (!response.success) {
      event.cookies.delete("refreshToken", { path: "/" });
      event.cookies.delete("loginToken", { path: "/" });

      return new Response(null, {
        status: 302,
        headers: { Location: "/signin" },
      });
    }
    return resolve(event);
  } catch (error) {
    console.error("Token validation error:", error);
    event.cookies.delete("refreshToken", { path: "/" });
    event.cookies.delete("loginToken", { path: "/" });

    return new Response(null, {
      status: 302,
      headers: { Location: "/signin" },
    });
  }
};
