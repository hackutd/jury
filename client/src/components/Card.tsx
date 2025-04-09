import { twMerge } from 'tailwind-merge';

interface CardProps {
    children: React.ReactNode;
    className?: string;
}

function Card(props: CardProps) {
    return (
        <div
            className={twMerge(
                'bg-white/50 border-2 border-backgroundDark rounded-md p-4 my-4 w-full',
                props.className
            )}
        >
            {props.children}
        </div>
    );
}

export default Card;
