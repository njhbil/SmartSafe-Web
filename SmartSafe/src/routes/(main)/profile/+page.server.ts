import { fail } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit';
import logoutAPI from '$lib/api/logoutAPI';


export const actions = {
	logout: async ({ cookies }: RequestEvent) => {
		try {
			const refreshToken = cookies.get('refreshToken') as string;
			const response = await logoutAPI({ refreshToken });
			if (response.success) {
				cookies.delete('refreshToken', {
					path: '/',
					httpOnly: true,
					secure: true,
					sameSite: 'lax',
				});
				cookies.delete('loginToken', {
					path: '/',
					httpOnly: true,
					secure: true,
					sameSite: 'lax',
				});

				return {
					status: 302,
					headers: { Location: '/' },
				};
			} else {
				return fail(400, { error: response.message || 'Logout failed' });
			}
		} catch (error) {
			console.error('Logout Error:', error);
			return fail(500, { error: 'An error occurred during logout. Please try again later.' });
		}
	},
};

