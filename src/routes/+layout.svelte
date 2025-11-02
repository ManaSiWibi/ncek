<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import * as NavigationMenu from "$lib/components/ui/navigation-menu/index.js";
	import { cn } from "$lib/utils.js";
	import { navigationMenuTriggerStyle } from "$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte";
	import type { HTMLAttributes } from "svelte/elements";
	import { page } from '$app/stores';
 
	let { data, children } = $props();
	const seo = data?.seo || {};
 
	const navigationItems: { title: string; href: string; description: string }[] = [
		{
		title: "Comprehensive Check",
		href: "/",
		description:
			"Run all checks at once - SSL, DNS, HTTP/3, web settings, and more."
		},
		{
		title: "SSL Certificate",
		href: "/ssl",
		description:
			"Check SSL certificate information, validity, and expiry dates."
		},
		{
		title: "HTTP/3 Test",
		href: "/http3",
		description:
			"Check if a website supports HTTP/3 and QUIC protocol."
		},
		{
		title: "DNS Lookup",
		href: "/dns",
		description:
			"Lookup DNS records including A, AAAA, CNAME, MX, TXT, and NS records."
		},
	{
		title: "IP Lookup",
		href: "/ip",
		description:
			"Get information about an IP address or resolve a domain to IP addresses."
	},
	{
		title: "My IP Address",
		href: "/my-ip",
		description:
			"Check your current IP address and location information."
	},
	{
		title: "Web Settings",
		href: "/web-settings",
		description:
			"Check web server settings, headers, and response information."
	},
		{
		title: "Email Config",
		href: "/email-config",
		description:
			"Check email authentication configuration including SPF, DKIM, DMARC, and BIMI."
		},
		{
		title: "Blocklist Check",
		href: "/blocklist",
		description:
			"Check if a domain is blocked by various DNS servers and blocklist services."
		},
		{
		title: "HSTS Check",
		href: "/hsts",
		description:
			"Check HTTP Strict Transport Security (HSTS) configuration and policy."
		},
		{
		title: "Robots.txt Check",
		href: "/robots-txt",
		description:
			"Check and analyze robots.txt file to see what search engines can crawl."
		},
		{
		title: "Sitemap Check",
		href: "/sitemap",
		description:
			"Check and analyze sitemap.xml file to see what pages are indexed."
		},
		{
		title: "OG Image Check",
		href: "/og-image",
		description:
			"Check and validate Open Graph image tags for social media sharing."
		}
	];
	
	type ListItemProps = HTMLAttributes<HTMLAnchorElement> & {
		title: string;
		href: string;
		content: string;
	};
</script>


{#snippet ListItem({
	title,
	content,
	href,
	class: className,
	...restProps
  }: ListItemProps)}
	<li>
	  <NavigationMenu.Link
		{href}
		aria-label={`${title} - ${content}`}
		class={cn(
		  "hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2",
		  className
		)}
		{...restProps}
	  >
		<div class="text-sm font-medium leading-none">{title}</div>
		<p class="text-muted-foreground line-clamp-2 text-sm leading-snug">
		  {content}
		</p>
	  </NavigationMenu.Link>
	</li>
  {/snippet}


<svelte:head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta name="description" content={seo.siteDescription || "Free comprehensive network and website checking tools"} />
	<meta name="keywords" content="SSL checker, DNS lookup, HTTP/3 test, email config checker, network tools, website analysis, OG image checker" />
	<meta name="author" content="NCEK" />
	<meta name="robots" content="index, follow" />
	<meta name="theme-color" content="#2563eb" />
	
	<!-- Open Graph / Facebook -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content={seo.siteUrl || "https://tools.kenadera.org"} />
	<meta property="og:title" content={seo.siteName || "NCEK - Network & Website Check Tools"} />
	<meta property="og:description" content={seo.siteDescription || "Free comprehensive tools to analyze SSL certificates, DNS records, HTTP/3 support, email configuration, and more."} />
	<meta property="og:site_name" content="NCEK" />
	
	<!-- Twitter -->
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content={seo.siteName || "NCEK - Network & Website Check Tools"} />
	<meta name="twitter:description" content={seo.siteDescription || "Free comprehensive network checking tools"} />
	
	<link rel="icon" href={favicon} />
	<link rel="canonical" href={`${seo.siteUrl || "https://tools.kenadera.org"}${$page.url.pathname}`} />
	<script defer src="https://umami.kenadera.org/script.js" data-website-id="0eb8a86a-8f82-4f78-9ca5-9bd205c0800c"></script>
</svelte:head>

<div class="min-h-screen flex flex-col">
	<!-- Skip to main content link for accessibility -->
	<a 
		href="#main-content" 
		class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-blue-600 focus:text-white focus:rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
	>
		Skip to main content
	</a>
	
	<!-- Header -->
	<header>
		<NavigationMenu.Root viewport={false} class="relative z-10 flex w-full self-center border-b border-gray-200 max-w-full
">
		
		<NavigationMenu.List class="group flex list-none p-1">
						
			<NavigationMenu.Item class="">
				<NavigationMenu.Link href="/" class={navigationMenuTriggerStyle()} aria-label="Home - NCEK Network Check Tools">
					Home
				</NavigationMenu.Link>
			</NavigationMenu.Item>

			  <NavigationMenu.Item class="">
				<NavigationMenu.Trigger class={navigationMenuTriggerStyle()} aria-label="Navigation menu - Available tools">Tools</NavigationMenu.Trigger>
				<NavigationMenu.Content class="">
				  <ul
					class="grid w-[400px] gap-2 p-2 md:w-[500px] md:grid-cols-2 lg:w-[600px]"
				  >
					{#each navigationItems as component, i (i)}
					  {@render ListItem({
						href: component.href,
						title: component.title,
						content: component.description
					  })}
					{/each}
				  </ul>
				</NavigationMenu.Content>
			  </NavigationMenu.Item>

		</NavigationMenu.List>
	</NavigationMenu.Root>
	</header>

	<!-- Main Content -->
	<main id="main-content" class="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		{@render children?.()}
	</main>

	<!-- Footer -->
	<footer class="bg-gray-50 border-t border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="grid grid-cols-1 md:grid-cols-3 gap-8">
				<!-- Brand Section -->
				<div>
					<h2 class="text-lg font-semibold text-gray-900 mb-4">NCEK</h2>
					<p class="text-gray-600 text-sm">
						Free Website and Network analysis, checking and troubleshooting tools.
					</p>
				</div>
				
				<!-- Links Section -->
				<nav aria-label="Quick links">
					<h2 class="text-sm font-semibold text-gray-900 mb-4">Quick Links</h2>
					<ul class="space-y-2" role="list">
						<li>
							<a href="/" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								Comprehensive Check
							</a>
						</li>
						<li>
							<a href="/ssl" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								SSL Certificate
							</a>
						</li>
						<li>
							<a href="/http3" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								HTTP/3 Test
							</a>
						</li>
						<li>
							<a href="/dns" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								DNS Lookup
							</a>
						</li>
						<li>
							<a href="/ip" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								IP Lookup
							</a>
						</li>
						<li>
							<a href="/web-settings" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								Web Settings
							</a>
						</li>
					<li>
						<a href="/email-config" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
							Email Config
						</a>
					</li>
					<li>
						<a href="/robots-txt" class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
							Robots.txt Check
						</a>
					</li>
				</ul>
				</nav>
				
				<!-- Info Section -->
				<div>
					<h2 class="text-sm font-semibold text-gray-900 mb-4">About</h2>
					<p class="text-gray-600 text-sm">
						A free web application for testing web and network connectivity and performance.
					</p>
				</div>
			</div>
			
			<!-- Bottom Bar -->
			<div class="mt-8 pt-8 border-t border-gray-200">
				<div class="flex flex-col md:flex-row justify-between items-center">
					<p class="text-gray-500 text-sm">
						© 2025 NCEK.
					</p>
					<div class="mt-4 md:mt-0">
						<p class="text-gray-500 text-sm">
							Built with SvelteKit & Tailwind CSS
						</p>
					</div>
				</div>
			</div>
		</div>
	</footer>
</div>
