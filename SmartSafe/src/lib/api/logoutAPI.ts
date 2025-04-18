export default async function logoutAPI({
	ResetToken,
}: {
	ResetToken: string;
}) {
	try {
		const response = await fetch("/api/logout", {
			method: "POST",
			headers: {
				"Content-Type": "application/json; charset=UTF-8",
			},
			body: JSON.stringify({ ResetToken }),
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
