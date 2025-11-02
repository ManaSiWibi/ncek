<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import * as NavigationMenu from "$lib/components/ui/navigation-menu/index.js";
	import { cn } from "$lib/utils.js";
	import { navigationMenuTriggerStyle } from "$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte";
	import type { HTMLAttributes } from "svelte/elements";
	import CircleHelpIcon from "@lucide/svelte/icons/circle-help";
	import CircleIcon from "@lucide/svelte/icons/circle";
	import CircleCheckIcon from "@lucide/svelte/icons/circle-check";
 
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
		}
	];
	
	type ListItemProps = HTMLAttributes<HTMLAnchorElement> & {
		title: string;
		href: string;
		content: string;
	};
	
	let { children } = $props();
</script>


{#snippet ListItem({
	title,
	content,
	href,
	class: className,
	...restProps
  }: ListItemProps)}
	<li>
	  <NavigationMenu.Link class="">
		{#snippet child()}
		  <a
			{href}
			data-sveltekit-reload
			class={cn(
			  "hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors",
			  className
			)}
			{...restProps}
		  >
			<div class="text-sm font-medium leading-none">{title}</div>
			<p class="text-muted-foreground line-clamp-2 text-sm leading-snug">
			  {content}
			</p>
		  </a>
		{/snippet}
	  </NavigationMenu.Link>
	</li>
  {/snippet}


<svelte:head>
	<link rel="icon" href={favicon} />
	<script defer src="https://umami.kenadera.org/script.js" data-website-id="0eb8a86a-8f82-4f78-9ca5-9bd205c0800c"></script>
</svelte:head>



<div class="min-h-screen flex flex-col">
	<!-- Header -->
	<NavigationMenu.Root viewport={false} class="relative z-10 flex w-full self-center border-b border-gray-200 max-w-full
">
		
		<NavigationMenu.List class="group flex list-none p-1">
						
			<NavigationMenu.Item>
				
				<NavigationMenu.Link class={navigationMenuTriggerStyle()}>
					<a href="/" data-sveltekit-reload>Home</a>
				</NavigationMenu.Link>
			</NavigationMenu.Item>

			  <NavigationMenu.Item>
				<NavigationMenu.Trigger>Components</NavigationMenu.Trigger>
				<NavigationMenu.Content>
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

	<!-- Main Content -->
	<main class="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		{@render children?.()}
	</main>

	<!-- Footer -->
	<footer class="bg-gray-50 border-t border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="grid grid-cols-1 md:grid-cols-3 gap-8">
				<!-- Brand Section -->
				<div>
					<h3 class="text-lg font-semibold text-gray-900 mb-4">NCEK</h3>
					<p class="text-gray-600 text-sm">
						Free Website and Network checking tools.
					</p>
				</div>
				
				<!-- Links Section -->
				<div>
					<h4 class="text-sm font-semibold text-gray-900 mb-4">Quick Links</h4>
					<ul class="space-y-2">
						<li>
							<a href="/" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								Comprehensive Check
							</a>
						</li>
						<li>
							<a href="/ssl" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								SSL Certificate
							</a>
						</li>
						<li>
							<a href="/http3" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								HTTP/3 Test
							</a>
						</li>
						<li>
							<a href="/dns" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								DNS Lookup
							</a>
						</li>
						<li>
							<a href="/ip" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								IP Lookup
							</a>
						</li>
						<li>
							<a href="/web-settings" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
								Web Settings
							</a>
						</li>
					<li>
						<a href="/email-config" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
							Email Config
						</a>
					</li>
					<li>
						<a href="/robots-txt" data-sveltekit-reload class="text-gray-600 hover:text-blue-600 text-sm transition-colors">
							Robots.txt Check
						</a>
					</li>
				</ul>
				</div>
				
				<!-- Info Section -->
				<div>
					<h4 class="text-sm font-semibold text-gray-900 mb-4">About</h4>
					<p class="text-gray-600 text-sm">
						A modern web application for testing network connectivity and performance.
					</p>
				</div>
			</div>
			
			<!-- Bottom Bar -->
			<div class="mt-8 pt-8 border-t border-gray-200">
				<div class="flex flex-col md:flex-row justify-between items-center">
					<p class="text-gray-500 text-sm">
						© 2024 NCEK. All rights reserved.
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
