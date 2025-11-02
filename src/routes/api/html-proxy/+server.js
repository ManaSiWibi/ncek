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
		const targetUrl = url.searchParams.get('url');
		
		if (!targetUrl) {
			throw error(400, 'URL parameter is required');
		}

		// Normalize URL
		let normalizedUrl = targetUrl;
		if (!normalizedUrl.startsWith('https://') && !normalizedUrl.startsWith('http://')) {
			normalizedUrl = 'https://' + normalizedUrl;
		}

		// Validate URL
		try {
			new URL(normalizedUrl);
		} catch {
			throw error(400, 'Invalid URL format');
		}

		// Fetch HTML content using Node.js fetch
		const controller = new AbortController();
		const timeoutId = setTimeout(() => controller.abort(), 30000); // 30 second timeout

		try {
			const response = await fetch(normalizedUrl, {
				method: 'GET',
				headers: {
					'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
					'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
					'Accept-Language': 'en-US,en;q=0.9',
				},
				signal: controller.signal,
				redirect: 'follow',
			});

			clearTimeout(timeoutId);

			if (!response.ok) {
				// Try HTTP if HTTPS fails
				if (normalizedUrl.startsWith('https://')) {
					const httpUrl = normalizedUrl.replace('https://', 'http://');
					const httpResponse = await fetch(httpUrl, {
						method: 'GET',
						headers: {
							'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
							'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
							'Accept-Language': 'en-US,en;q=0.9',
						},
						signal: controller.signal,
						redirect: 'follow',
					});

					if (httpResponse.ok) {
						const html = await httpResponse.text();
						return new Response(JSON.stringify({
							success: true,
							data: {
								url: httpUrl,
								html: html,
								status: httpResponse.status
							}
						}), {
							headers: {
								'Content-Type': 'application/json'
							}
						});
					}
				}

				throw error(response.status, `Failed to fetch: HTTP ${response.status}`);
			}

			const html = await response.text();

			return new Response(JSON.stringify({
				success: true,
				data: {
					url: normalizedUrl,
					html: html,
					status: response.status
				}
			}), {
				headers: {
					'Content-Type': 'application/json'
				}
			});
		} catch (fetchError) {
			clearTimeout(timeoutId);
			
			if (fetchError instanceof Error && fetchError.name === 'AbortError') {
				throw error(408, 'Request timeout');
			}
			
			const errorMessage = fetchError instanceof Error ? fetchError.message : String(fetchError);
			throw error(500, `Failed to fetch HTML: ${errorMessage}`);
		}
	} catch (err) {
		if (err && typeof err === 'object' && 'status' in err) {
			throw err; // Re-throw SvelteKit errors
		}
		const errorMessage = err instanceof Error ? err.message : String(err);
		throw error(500, `HTML proxy error: ${errorMessage}`);
	}
}

