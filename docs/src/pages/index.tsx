import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
    const { siteConfig } = useDocusaurusContext();
    return (
        <header className={styles.heroBanner}>
            <div className="container">
                <Heading as="h1" className={styles.heroTitle}>
                    Jury
                </Heading>
                <p className={styles.heroTagline}>{siteConfig.tagline}</p>
                <div className={styles.buttons}>
                    <Link className="button button--secondary button--lg" to="/docs/usage">
                        Get Started 🚀
                    </Link>
                </div>
                <div className={styles.description}>
                    <h1>Features</h1>
                    <ul className={styles.list}>
                        <li>Automated judging application with one-click deployment</li>
                        <li>User-friendly judging interface with randomized project assignments</li>
                        <li>
                            Full-featured admin dashboard with a variety of configurations and
                            export options
                        </li>
                        <li>
                            Accurate rank aggregation algorithm with aggregated sub-category scores
                        </li>
                        <li>Built with modern technologies to run blazingly fast</li>
                        <li>Quick and easy import from Devpost</li>
                        <li>Self-hosting on DigitalOcean and MongoDB Atlas for almost no cost</li>
                    </ul>
                    <h1 style={{ marginTop: '2rem' }}>Don't Believe Us?</h1>
                    <Link
                        className="button button--primary button--lg"
                        to="https://jury-dev.mikz.dev"
                    >
                        Try it Out Online, right now!
                    </Link>
                    <p>^ Admin password is "admin"</p>
                </div>
            </div>
        </header>
    );
}

export default function Home(): JSX.Element {
    const { siteConfig } = useDocusaurusContext();
    return (
        <Layout description="An all-in-one hackathon judging platform, created with modern technologies aimed at optimizing the user experience of hackers, judges, and organizers.">
            <HomepageHeader />
            <main></main>
        </Layout>
    );
}
