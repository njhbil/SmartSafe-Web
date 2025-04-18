export default async function sendEmailOTP({ email }: { email: string }) {
	try {
		const localurl = import.meta.env.VITE_API_URL
			? import.meta.env.VITE_API_URL
			: "localhost";
		const response = await fetch(localurl + "/api/sendemailotp", {
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
		console.error("postFetch error:", error);
		return { success: false, message: "Unexpected error occurred" };
	}
}
