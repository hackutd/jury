import Link from '@docusaurus/Link';

interface CalloutButtonProps {
    label: string;
    to: string;
}

export default function CalloutButton(props: CalloutButtonProps) {
    return (
        <Link
            className="button button--primary button--lg button--block"
            to={props.to}
            style={{ marginBottom: '1em' }}
        >
            {props.label}
        </Link>
    );
}
