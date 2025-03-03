export default async function accountsAPI() {
	const localurl = (import.meta.env.VITE_API_URL ? import.meta.env.VITE_API_URL : "localhost")
	const fetched = await fetch(localurl + "/api/accounts");
	const data = await fetched.json();
	return data;
}
