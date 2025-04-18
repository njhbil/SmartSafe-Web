export default async function verifyEmailOTP({
	email,
	oneTimeCode,
}: {
	email: string;
	oneTimeCode: string;
}) {
	try {
		const localurl = import.meta.env.VITE_API_URL
			? import.meta.env.VITE_API_URL
			: "localhost";
		const response = await fetch(localurl + "/api/verifyemailotp", {
			method: "POST",
			headers: {
				"Content-Type": "application/json; charset=UTF-8",
			},
			body: JSON.stringify({ email, oneTimeCode }),
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
