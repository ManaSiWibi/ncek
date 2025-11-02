/** @type {import('./$types').LayoutServerLoad} */
export function load() {
	const siteName = 'NCEK - Network & Website Check Tools';
	const siteDescription = 'Free comprehensive tools to analyze SSL certificates, DNS records, HTTP/3 support, email configuration, OG images, and more. All network checking tools in one place.';
	const siteUrl = 'https://tools.kenadera.org';
	
	return {
		seo: {
			siteName,
			siteDescription,
			siteUrl
		}
	};
}

