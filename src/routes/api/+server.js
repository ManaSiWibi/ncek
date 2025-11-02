import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

/** @type {import('./$types').RequestHandler} */
export async function GET({ url, request }) {
	// Validate origin to prevent direct API access
	const allowedOrigins = env.ALLOWED_ORIGIN ? env.ALLOWED_ORIGIN.split(',').map(o => o.trim()) : ['http://localhost:3001'];
	const origin = request.headers.get('origin');
	const referer = request.headers.get('referer');
	
	// Check if request is from allowed origin
	let isAllowed = false;
	
	// Priority 1: Check Origin header (most reliable for browser requests)
	if (origin) {
		try {
			const originUrl = new URL(origin);
			isAllowed = allowedOrigins.some(allowed => {
				try {
					const allowedUrl = new URL(allowed);
					return allowedUrl.origin === originUrl.origin;
				} catch {
					return allowed === origin || allowed === originUrl.origin;
				}
			});
		} catch {
			// Invalid origin URL
		}
	}
	
	// Priority 2: Fallback to Referer header if Origin is not present
	if (!isAllowed && referer) {
		try {
			const refererUrl = new URL(referer);
			isAllowed = allowedOrigins.some(allowed => {
				try {
					const allowedUrl = new URL(allowed);
					return allowedUrl.origin === refererUrl.origin;
				} catch {
					return allowed === referer || allowed === refererUrl.origin;
				}
			});
		} catch {
			// Invalid referer URL
		}
	}
	
	// Reject if no valid Origin or Referer header found
	// This prevents direct API access via curl, Postman, etc.
	if (!isAllowed) {
		throw error(403, 'Direct API access not allowed. Please use the frontend application.');
	}
	try {
		const type = url.searchParams.get('type');
		const domain = url.searchParams.get('domain');
		const ip = url.searchParams.get('ip');
		const urlParam = url.searchParams.get('url');
		
		// Map client-friendly types to backend API endpoints
		/** @type {Record<string, string>} */
		const endpointMap = {
			'ssl': 'ssl',
			'http3': 'http3',
			'dns': 'dns',
			'ip': 'ip',
			'web-settings': 'web-settings',
			'web_settings': 'web-settings', // Support both formats
			'email-config': 'email-config',
			'email_config': 'email-config',
		'blocklist': 'blocklist',
		'hsts': 'hsts',
		'robots-txt': 'robots-txt',
		'sitemap': 'sitemap',
		'html-proxy': 'html-proxy',
		'comprehensive': 'comprehensive',
		'full': 'comprehensive'
		};
		
		// If no type specified, default to comprehensive
		const endpoint = (type && endpointMap[type]) || 'comprehensive';
		
		// Build the backend URL using environment variable or default
		const backendHost = env.BACKEND_URL || 'http://localhost:8080';
		let backendUrl = `${backendHost}/api/v1/${endpoint}`;
		const params = new URLSearchParams();
		
		if (domain) {
			params.append('domain', domain);
		}
		if (ip) {
			params.append('ip', ip);
		}
		if (urlParam) {
			params.append('url', urlParam);
		}
		
		if (params.toString()) {
			backendUrl += `?${params.toString()}`;
		}
		
		// Make request to the backend API with internal proxy header
		const response = await fetch(backendUrl, {
			headers: {
				'X-Internal-Proxy': 'true'
			}
		});
		
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		
		const data = await response.json();
		
		return new Response(JSON.stringify(data), {
			headers: {
				'Content-Type': 'application/json'
			}
		});
	} catch (err) {
		const errorMessage = err instanceof Error ? err.message : String(err);
		throw error(500, `Failed to fetch API data: ${errorMessage}`);
	}
}

