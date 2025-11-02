/**
 * API utility functions for making requests to the backend API
 * These functions call server endpoints that proxy to the Go backend API
 * Note: The API secret is handled server-side only and never exposed to the browser
 */

/**
 * Fetch data from a specific API endpoint
 * @param {string} endpoint - API endpoint (e.g., '/api', '/api/html-proxy')
 * @param {object} params - Query parameters
 * @returns {Promise<any>} - API response data
 */
export async function fetchApiData(endpoint, params = {}) {
	try {
		// Build query string from params
		const queryParams = new URLSearchParams();
		Object.entries(params).forEach(([key, value]) => {
			if (value) {
				queryParams.append(key, value);
			}
		});
		
		const queryString = queryParams.toString();
		const url = `${endpoint}${queryString ? `?${queryString}` : ''}`;
		
		// Make request to server-side API route (secret is handled server-side)
		const response = await fetch(url);
		
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		
		const result = await response.json();
		
		// Handle different response structures
		if (result.success && result.data) {
			return result.data;
		} else if (result.error) {
			return { error: result.error };
		} else {
			return result;
		}
	} catch (error) {
		console.error(`Error fetching ${endpoint}:`, error);
		const errorMessage = error instanceof Error ? error.message : String(error);
		return { error: `Failed to fetch data: ${errorMessage}` };
	}
}

/**
 * Check SSL certificate information
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - SSL certificate data
 */
export async function checkSSL(domain) {
	return fetchApiData('/api', { type: 'ssl', domain });
}

/**
 * Check HTTP/3 support
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - HTTP/3 support data
 */
export async function checkHTTP3(domain) {
	return fetchApiData('/api', { type: 'http3', domain });
}

/**
 * Get DNS information
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - DNS data
 */
export async function checkDNS(domain) {
	return fetchApiData('/api', { type: 'dns', domain });
}

/**
 * Check IP address information
 * @param {string} ip - IP address to check
 * @returns {Promise<any>} - IP data
 */
export async function checkIP(ip) {
	return fetchApiData('/api', { type: 'ip', ip });
}

/**
 * Check web server settings
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - Web server settings data
 */
export async function checkWebSettings(domain) {
	return fetchApiData('/api', { type: 'web-settings', domain });
}

/**
 * Check email authentication configuration (SPF, DKIM, DMARC, BIMI)
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - Email authentication data
 */
export async function checkEmailConfig(domain) {
	return fetchApiData('/api', { type: 'email-config', domain });
}

/**
 * Check domain blocklist status
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - Blocklist data
 */
export async function checkBlocklist(domain) {
	return fetchApiData('/api', { type: 'blocklist', domain });
}

/**
 * Check HSTS (HTTP Strict Transport Security) status
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - HSTS data
 */
export async function checkHSTS(domain) {
	return fetchApiData('/api', { type: 'hsts', domain });
}

/**
 * Comprehensive check (all information at once)
 * @param {string} domain - Domain to check
 * @returns {Promise<any>} - Comprehensive data
 */
export async function checkComprehensive(domain) {
	return fetchApiData('/api', { type: 'comprehensive', domain });
}

/**
 * Fetch HTML content via proxy
 * @param {string} url - URL to fetch
 * @returns {Promise<any>} - HTML content
 */
export async function fetchHTMLProxy(url) {
	const response = await fetch(`/api?type=html-proxy&url=${encodeURIComponent(url)}`);
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}
	const result = await response.json();
	return result.success && result.data ? result.data : result;
}

