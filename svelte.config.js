import adapter from '@sveltejs/adapter-node';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		// Using adapter-node for Docker deployment
		adapter: adapter(),
		csrf: {
			checkOrigin: false
		}
	}
};

export default config;
