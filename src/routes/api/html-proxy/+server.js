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
	const providedSecret = secretFromQuery || secretFromHeader;
	
	// Check if request is same-origin (from the frontend itself)
	const origin = request.headers.get('origin');
	const referer = request.headers.get('referer');
	const host = request.headers.get('host');
	
	let isSameOrigin = false;
	if (origin && host) {
		try {
			const originUrl = new URL(origin);
			// Check if origin matches the current host
			isSameOrigin = originUrl.host === host || originUrl.hostname === host.split(':')[0];
		} catch {
			// Invalid origin URL
		}
	}
	
	// If no origin header, check referer as fallback
	if (!isSameOrigin && referer && host) {
		try {
			const refererUrl = new URL(referer);
			isSameOrigin = refererUrl.host === host || refererUrl.hostname === host.split(':')[0];
		} catch {
			// Invalid referer URL
		}
	}
	
	// For same-origin requests (from frontend), automatically use server-side secret
	// For external requests, require secret to be provided
	if (isSameOrigin) {
		// Same-origin request - automatically validated using server-side secret
		// No need to check providedSecret, as we use the server-side apiSecret
	} else if (providedSecret !== apiSecret) {
		// External request - must provide valid secret
		throw error(403, 'Invalid or missing API secret key');
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

