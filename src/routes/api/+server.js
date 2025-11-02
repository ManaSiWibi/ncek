import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

/** @type {import('./$types').RequestHandler} */
export async function GET({ url, request }) {
	// Require API secret key from environment (server-side only, never exposed to client)
	const apiSecret = env.API_SECRET_KEY;
	if (!apiSecret) {
		throw error(500, 'API secret key not configured');
	}

	// Note: This route is server-side only. The secret is never exposed to the browser.
	// The backend API will validate the secret when we forward the request.
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
			'og-image': 'og-image',
			'og_image': 'og-image',
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
		
		// For og-image endpoint, prefer url parameter, fallback to domain
		if (endpoint === 'og-image') {
			if (urlParam) {
				params.append('url', urlParam);
			} else if (domain) {
				params.append('url', domain);
			}
		} else {
			if (domain) {
				params.append('domain', domain);
			}
			if (urlParam) {
				params.append('url', urlParam);
			}
		}
		if (ip) {
			params.append('ip', ip);
		}
		
		if (params.toString()) {
			backendUrl += `?${params.toString()}`;
		}
		
		// Make request to the backend API with internal proxy header and secret key
		// Secret key is server-side only, never exposed to browser
		// Forward the real client IP for rate limiting
		const clientIP = request.headers.get('x-forwarded-for')?.split(',')[0]?.trim() ||
		                 request.headers.get('x-real-ip') ||
		                 request.headers.get('cf-connecting-ip') || // Cloudflare
		                 null;
		
		/** @type {Record<string, string>} */
		const headers = {
			'X-Internal-Proxy': 'true',
			'X-API-Secret': apiSecret
		};
		
		// Forward client IP if available
		if (clientIP) {
			headers['X-Forwarded-For'] = clientIP;
		}
		
		const response = await fetch(backendUrl, { headers });
		
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

