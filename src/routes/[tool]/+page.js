import { getToolConfig, tools } from '$lib/tools.js';

/** @type {import('./$types').PageLoad} */
export function load({ params }) {
	const { tool } = params;
	
	// Try direct lookup first
	let config = getToolConfig(tool);
	
	// If not found, search by path
	if (!config) {
		for (const [key, cfg] of Object.entries(tools)) {
			// @ts-ignore
			if (cfg.path === `/${tool}`) {
				config = cfg;
				break;
			}
		}
	}
	
	if (!config) {
		return {
			status: 404,
			error: new Error('Tool not found')
		};
	}
	return { config };
}
