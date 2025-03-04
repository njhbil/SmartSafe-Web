export default async function loginAPI({
  email,
  password,
  rememberMe,
}: {
  email: string;
  password: string;
  rememberMe: boolean;
}) {
  const localurl = import.meta.env.VITE_API_URL
    ? import.meta.env.VITE_API_URL
    : "localhost";

  try {
    const response = await fetch(localurl + "/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json; charset=UTF-8",
      },
      body: JSON.stringify({ email, password, rememberMe }),
    });

    const json = await response.json();

    if (!response.ok) {
      return {
        success: false,
        message: json.message || "Failed to authenticate",
      };
    }

    return { success: true, ...json };
  } catch (error) {
    console.error("AuthToken error:", error);
    return { success: false, message: "Unexpected error occurred" };
  }
}
