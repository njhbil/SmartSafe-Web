export default async function logoutAPI({
	refreshToken,
}: {
	refreshToken: string;
}) {
	const localurl = import.meta.env.VITE_API_URL
		? import.meta.env.VITE_API_URL
		: "localhost";
	try {
		const response = await fetch(localurl + "/api/logout", {
			method: "POST",
			headers: {
				"Content-Type": "application/json; charset=UTF-8",
			},
			body: JSON.stringify({ refreshToken }),
		});

		if (!response.ok) {
			throw new Error("Logout failed");
		}

		return { success: true, message: "Logout successful" };
	} catch (error) {
		console.error("Logout API Error:", error);
		return { success: false, message: "An error occurred during logout" };
	}
}
