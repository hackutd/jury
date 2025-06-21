import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
    title: 'Jury',
    tagline: 'A modern hackathon judging platform',
    favicon: 'img/favicon.ico',

    // Set the production url of your site here
    url: 'https://jury.mikz.dev',
    // Set the /<baseUrl>/ pathname under which your site is served
    // For GitHub pages deployment, it is often '/<projectName>/'
    baseUrl: '/',

    // GitHub pages deployment config.
    // If you aren't using GitHub pages, you don't need these.
    organizationName: 'hackutd', // Usually your GitHub org/user name.
    projectName: 'jury', // Usually your repo name.

    onBrokenLinks: 'throw',
    onBrokenMarkdownLinks: 'warn',

    // Even if you don't use internationalization, you can use this field to set
    // useful metadata like html lang. For example, if your site is Chinese, you
    // may want to replace "en" with "zh-Hans".
    i18n: {
        defaultLocale: 'en',
        locales: ['en'],
    },

    presets: [
        [
            'classic',
            {
                docs: {
                    sidebarPath: './sidebars.ts',
                    editUrl: 'https://github.com/hackutd/jury/tree/master',
                },
                blog: false,
                theme: {
                    customCss: './src/css/custom.css',
                },
            } satisfies Preset.Options,
        ],
    ],

    themeConfig: {
        image: 'img/social-card.jpg',
        navbar: {
            title: 'Jury',
            logo: {
                alt: 'Jury Logo',
                src: 'img/logo.svg',
            },
            items: [
                {
                    type: 'doc',
                    position: 'right',
                    docId: 'usage/deploy',
                    label: 'Get Started',
                },
                {
                    type: 'doc',
                    position: 'right',
                    docId: '/category/reference',
                    label: 'Reference',
                },
                {
                    href: 'https://github.com/hackutd/jury',
                    label: 'GitHub',
                    position: 'right',
                },
            ],
        },
        footer: {
            style: 'dark',
            links: [
                {
                    title: 'Docs',
                    items: [
                        {
                            label: 'Deploy Jury',
                            to: '/docs/usage/deploy',
                        },
                        {
                            label: 'Contributing',
                            to: '/docs/contributing',
                        },
                        {
                            label: 'Public API',
                            to: '/docs/reference/public-api',
                        },
                        {
                            label: 'Internal API Reference',
                            to: '/docs/reference/internal-api',
                        },
                    ],
                },
                {
                    title: 'More',
                    items: [
                        {
                            label: 'GitHub',
                            href: 'https://github.com/hackutd/jury',
                        },
                    ],
                },
            ],
            copyright: `Copyright © ${new Date().getFullYear()} HackUTD. Built with Docusaurus.`,
        },
        prism: {
            theme: prismThemes.github,
            darkTheme: prismThemes.dracula,
        },
    } satisfies Preset.ThemeConfig,
};

export default config;
