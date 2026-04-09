import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
	title: 'opentelemetry-go',
	description: 'Collection of extensions for OpenTelemetry-Go',
	lang: "en-US",
	lastUpdated: true,
	appearance: "dark",
	base: '/opentelemetry-go/',
	sitemap: {
		hostname: 'https://foomo.github.io/opentelemetry-go',
	},
	themeConfig: {
		// https://vitepress.dev/reference/default-theme-config
		logo: '/logo.png',
		outline: [2, 4],
		nav: [
			{ text: 'Guide', link: '/guide/getting-started' },
			{ text: 'API Reference', link: '/api/' },
			{ text: 'Examples', link: '/examples/basic-tracing' },
		],
		sidebar: [
			{
				text: 'Guide',
				items: [
					{ text: 'Getting Started', link: '/guide/getting-started' },
					{ text: 'Exporters', link: '/guide/exporters' },
					{ text: 'Semantic Conventions', link: '/guide/semconv' },
					{ text: 'Testing', link: '/guide/testing' },
				],
			},
			{
				text: 'API Reference',
				items: [
					{ text: 'Overview', link: '/api/' },
					{ text: 'glossytrace', link: '/api/glossytrace' },
					{ text: 'glossymetric', link: '/api/glossymetric' },
					{ text: 'semconv', link: '/api/semconv' },
					{ text: 'oteltesting', link: '/api/testing' },
				],
			},
			{
				text: 'Examples',
				items: [
					{ text: 'Basic Tracing', link: '/examples/basic-tracing' },
					{ text: 'Basic Metrics', link: '/examples/basic-metrics' },
					{ text: 'Testing', link: '/examples/testing-example' },
					{ text: 'Custom Semconv', link: '/examples/custom-semconv' },
				],
			},
			{
				text: 'Contributing',
				collapsed: true,
				items: [
					{ text: 'Guideline', link: '/CONTRIBUTING' },
					{ text: 'Code of Conduct', link: '/CODE_OF_CONDUCT' },
					{ text: 'Security', link: '/SECURITY' },
				],
			},
		],
		socialLinks: [
			{ icon: 'github', link: 'https://github.com/foomo/opentelemetry-go' },
		],
		editLink: {
			pattern: 'https://github.com/foomo/opentelemetry-go/edit/main/docs/:path',
		},
		search: {
			provider: 'local',
		},
		footer: {
			message: 'Made with ♥ <a href="https://www.foomo.org">foomo</a> by <a href="https://www.bestbytes.com">bestbytes</a>',
		},
	},
	markdown: {
		// https://github.com/vuejs/vitepress/discussions/3724
		theme: {
			light: 'catppuccin-latte',
			dark: 'catppuccin-frappe',
		}
	},
	head: [
		['meta', { name: 'theme-color', content: '#ffffff' }],
		['link', { rel: 'icon', href: '/logo.png' }],
		['meta', { name: 'author', content: 'foomo by bestbytes' }],
		// OpenGraph
		['meta', { property: 'og:title', content: 'foomo/opentelemetry-go' }],
		[
			'meta',
			{
				property: 'og:image',
				content: 'https://github.com/foomo/opentelemetry-go/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta',
			{
				property: 'og:description',
				content: 'Stop using `go func`, start using `opentelemetry-go`',
			},
		],
		['meta', { name: 'twitter:card', content: 'summary_large_image' }],
		[
			'meta',
			{
				name: 'twitter:image',
				content: 'https://github.com/foomo/opentelemetry-go/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0, viewport-fit=cover',
			},
		],
	]
})
