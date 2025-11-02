/** @type {import('./$types').RequestHandler} */
export async function GET({ url }) {
	// Get base URL from the request origin
	// You can set PUBLIC_BASE_URL in .env if you need to override this
	const baseURL = `${url.protocol}//${url.host}`;

	// Define all public routes with their priorities and change frequencies
	const routes = [
		{ loc: '/', changefreq: 'daily', priority: '1.0' }, // Home/Comprehensive
		{ loc: '/ssl', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/http3', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/dns', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/ip', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/web-settings', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/email-config', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/blocklist', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/hsts', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/robots-txt', changefreq: 'monthly', priority: '0.9' },
		{ loc: '/sitemap', changefreq: 'monthly', priority: '0.8' },
		{ loc: '/og-image', changefreq: 'monthly', priority: '0.9' },
	];

	// Get current date in ISO format
	const lastmod = new Date().toISOString().split('T')[0];

	// Generate XML sitemap
	const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes
	.map(
		(route) => `  <url>
    <loc>${baseURL}${route.loc}</loc>
    <lastmod>${lastmod}</lastmod>
    <changefreq>${route.changefreq}</changefreq>
    <priority>${route.priority}</priority>
  </url>`
	)
	.join('\n')}
</urlset>`;

	return new Response(sitemap, {
		headers: {
			'Content-Type': 'application/xml; charset=utf-8',
			'Cache-Control': 'public, max-age=3600' // Cache for 1 hour
		}
	});
}

