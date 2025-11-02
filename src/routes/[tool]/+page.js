import { getToolConfig, tools } from '$lib/tools.js';
import { error } from '@sveltejs/kit';

/** @type {import('./$types').PageLoad} */
export function load({ params, depends }) {
	const { tool } = params;
	
	// Ensure this load function re-runs when the tool param changes
	depends(`tool:${tool}`);
	
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
		throw error(404, 'Tool not found');
	}
	return { config };
}
