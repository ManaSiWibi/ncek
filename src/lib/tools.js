/**
 * @typedef {Object} ToolConfig
 * @property {string} title
 * @property {string} description
 * @property {string} placeholder
 * @property {string} formField
 * @property {(value: string) => Record<string, any>} apiParams
 * @property {string} resultCard
 * @property {(data: any) => any} transform
 * @property {string} [path]
 * @property {boolean} [clientSide]
 */

/**
 * @type {Record<string, ToolConfig>}
 */
export const tools = {
	// key becomes the route path: /[key]
	comprehensive: {
		title: 'Comprehensive Check',
		description: 'Run all checks at once - SSL, DNS, HTTP/3, web settings, and more.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'comprehensive', domain: value }),
		resultCard: 'raw',
		transform: (data) => data
	},
	ssl: {
		title: 'SSL Certificate Check',
		description: 'Check SSL certificate information for any domain',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'ssl', domain: value }),
		resultCard: 'ssl',
		transform: (data) => data
	},
	http3: {
		title: 'HTTP/3 Test',
		description: 'Check if a website supports HTTP/3 and QUIC protocol.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'http3', domain: value }),
		resultCard: 'http3',
		transform: (data) => data
	},
	dns: {
		title: 'DNS Lookup',
		description: 'Lookup DNS records including A, AAAA, CNAME, MX, TXT, and NS records.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'dns', domain: value }),
		resultCard: 'dns',
		transform: (data) => data
	},
	ip: {
		title: 'IP Address Lookup',
		description: 'Get information about an IP address or resolve a domain to IP addresses.',
		placeholder: 'e.g., 8.8.8.8 or example.com',
		formField: 'ip',
		apiParams: (value) => ({ type: 'ip', ip: value }),
		resultCard: 'ip',
		transform: (data) => data
	},
	"web-settings": {
		title: 'Web Settings',
		description: 'Check web server settings, headers, and response information.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'web-settings', domain: value }),
		resultCard: 'web_settings',
		transform: (data) => data
	},
	"email-config": {
		title: 'Email Config',
		description: 'Check SPF, DKIM, DMARC, and BIMI.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'email-config', domain: value }),
		resultCard: 'email_config',
		transform: (data) => data
	},
	blocklist: {
		title: 'Blocklist Check',
		description: 'Check if a domain is blocked by various DNS servers and blocklist services.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'blocklist', domain: value }),
		resultCard: 'raw',
		transform: (data) => data
	},
	hsts: {
		title: 'HSTS Check',
		description: 'Check HTTP Strict Transport Security (HSTS) configuration and policy.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'hsts', domain: value }),
		resultCard: 'raw',
		transform: (data) => data
	},
	robots: {
		title: 'Robots.txt Check',
		description: 'Check and analyze robots.txt file to see what search engines can crawl.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		path: '/robots-txt',
		apiParams: (value) => ({ type: 'robots-txt', domain: value }),
		resultCard: 'raw',
		transform: (data) => data
	},
	sitemap: {
		title: 'Sitemap Check',
		description: 'Check and analyze sitemap.xml file to see what pages are indexed.',
		placeholder: 'e.g., example.com',
		formField: 'description',
		apiParams: (value) => ({ type: 'sitemap', domain: value }),
		resultCard: 'raw',
		transform: (data) => data
	}
};

export function getToolConfig(slug) {
	return tools[slug] ?? null;
}
