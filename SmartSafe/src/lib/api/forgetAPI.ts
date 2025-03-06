export default async function forgetAPI({ email }: { email: string }) {
	const localurl = import.meta.env.VITE_API_URL
		? import.meta.env.VITE_API_URL
		: "localhost";

	try {
		const response = await fetch(localurl + "/api/forgetpassword", {
			method: "POST",
			headers: {
				"Content-Type": "application/json; charset=UTF-8",
			},
			body: JSON.stringify({ email }),
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
