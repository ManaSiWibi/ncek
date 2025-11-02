import adapter from '@sveltejs/adapter-node';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		// Using adapter-node for Docker deployment
		adapter: adapter(),
		csrf: {
			trustedOrigins: ['https://tools.kenadera.org', 'http://localhost:3001'],
			// checkOrigin: false
		}
	}
};

export default config;
