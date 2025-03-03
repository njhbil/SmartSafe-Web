import type { Interface } from "./interface";

export default async function registerAPI({ username, email, password }: Interface) {
	try {

		const localurl = (import.meta.env.VITE_API_URL ? import.meta.env.VITE_API_URL : "localhost")
		const response = await fetch(localurl + "/api/register", {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json; charset=UTF-8',
			},
			body: JSON.stringify({ username, email, password }),
		});

		const json = await response.json();

		if (!response.ok) {
			return { success: false, message: json.message || "Failed to authenticate" };
		}

		return { success: true, ...json };
	} catch (error) {
		console.error("postFetch error:", error);
		throw error;
	}
}
