export default async function resetAPI({
	resetPWDtoken,
	password,
}: {
	resetPWDtoken: string;
	password: string;
}) {
	const localurl = import.meta.env.VITE_API_URL
		? import.meta.env.VITE_API_URL
		: "localhost";

	try {
		const response = await fetch(localurl + "/api/resetpassword", {
			method: "POST",
			headers: {
				"Content-Type": "application/json; charset=UTF-8",
			},
			body: JSON.stringify({ resetPWDtoken, password }),
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
		console.error("resetpw error: ", error);
		return { success: false, message: "Unexpected error occurred" };
	}
}
