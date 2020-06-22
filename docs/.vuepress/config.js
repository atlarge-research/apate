// .vuepress/config.js
module.exports = {
	title: "Apate",
	themeConfig: {
		sidebar: 'auto',
		search: true,
		lastUpdated: true,
		nav: [
			{ text: 'Home', link: '/' },
			{ text: 'Usage', link: '/usage/' },
			{ text: 'CRD Configuration', link: '/configuration/' },
			{ text: 'Metrics', link: '/metrics/' },
			{ text: 'Examples', link: '/examples/' },
			{ 
				text: 'Development',
				items: [
					{ text: 'Build', link: '/build/' },
					{ text: 'Environment variables', link: '/env/' },
					{ text: 'Design & Implementation Details', link: '/ApateDesignImplementation.pdf', target:'_blank' },
				]
			}
		],
		// Edit links
		repo: 'atlarge-research/apate',
		docsDir: 'docs',
		editLinks: true,
	}
}
