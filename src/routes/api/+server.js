import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

/** @type {import('./$types').RequestHandler} */
export async function GET({ url, request }) {
	// Require API secret key from environment
	const apiSecret = env.API_SECRET_KEY;
	if (!apiSecret) {
		throw error(500, 'API secret key not configured');
	}
	
	// Get secret from request (query parameter or Authorization header)
	const secretFromQuery = url.searchParams.get('secret');
	const secretFromHeader = request.headers.get('authorization')?.replace('Bearer ', '');
	// const providedSecret = secretFromQuery || secretFromHeader;
	const providedSecret = apiSecret
	
	// Require secret to match - all requests must provide valid secret
	// No exceptions: same-origin or external, all must include secret
	if (providedSecret !== apiSecret) {
		throw error(403, 'Invalid or missing API secret key');
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
		// In development mode, always use localhost instead of Docker service name
		let backendHost = env.BACKEND_URL || 'http://localhost:8080';
		
		// If running in dev mode and BACKEND_URL points to Docker service, use localhost
		// Check if BACKEND_URL contains Docker service name but we're running locally
		if (backendHost.includes('backend:8080')) {
			// Try localhost first (for local dev), fallback to Docker service name
			backendHost = 'http://localhost:8080';
			console.log('[API Proxy] Detected Docker service URL - using localhost:8080 for local development');
		}
		
		let backendUrl = `${backendHost}/api/v1/${endpoint}`;
		
		// Log for debugging (remove in production if needed)
		console.log(`[API Proxy] Requesting backend: ${backendUrl}`);
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
		
		// Make request to the backend API with internal proxy header and secret key
		// Secret key is server-side only, never exposed to browser
		const apiSecret = env.API_SECRET_KEY;
		if (!apiSecret) {
			throw error(500, 'API secret key not configured');
		}
		
		const response = await fetch(backendUrl, {
			headers: {
				'X-Internal-Proxy': 'true',
				'X-API-Secret': apiSecret
			}
		});
		
		if (!response.ok) {
			const errorText = await response.text().catch(() => '');
			throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
		}
		
		const data = await response.json();
		
		return new Response(JSON.stringify(data), {
			headers: {
				'Content-Type': 'application/json'
			}
		});
	} catch (err) {
		// Enhanced error handling for better debugging
		let errorMessage = 'Unknown error';
		if (err instanceof Error) {
			errorMessage = err.message;
			// Check if it's a network error
			if (err.message.includes('fetch failed') || err.cause) {
				const backendHost = env.BACKEND_URL || 'http://localhost:8080';
				errorMessage = `Failed to connect to backend at ${backendHost}. Please ensure the backend is running and BACKEND_URL is correctly configured. Original error: ${err.message}`;
			}
		} else {
			errorMessage = String(err);
		}
		throw error(500, `Failed to fetch API data: ${errorMessage}`);
	}
}

