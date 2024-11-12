import React from 'react';
import { twMerge } from 'tailwind-merge';

interface StarProps {
    /* Star active state */
    active: boolean;

    /* Setter function for the star */
    setActive: React.Dispatch<React.SetStateAction<boolean>>;

    /* True if the star is disabled */
    disabled?: boolean;

    /* Smaller star */
    small?: boolean;

    /* Function to run on click */
    onClick?: () => void;

    /* Additional classes */
    className?: string;
}

const Star = (props: StarProps) => {
    return (
        <button
            className={props.className}
            onClick={() => {
                if (props.disabled) return;

                props.setActive(!props.active);
                if (props.onClick) props.onClick();
            }}
        >
            <svg
                width="36"
                height="35"
                viewBox="0 0 36 35"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                className={twMerge('w-5 h-5', props.small && 'w-4 h-4', props.disabled && 'cursor-default')}
            >
                <path
                    d="M9.95714 29.4L17.8071 24.7L25.6571 29.45L23.5571 20.55L30.4571 14.55L21.3571 13.75L17.8071 5.35L14.2571 13.7L5.15714 14.5L12.0571 20.5L9.95714 29.4ZM17.8071 28.2L8.50714 33.8C8.24048 33.9667 7.95714 34.0417 7.65714 34.025C7.35714 34.0083 7.09048 33.9167 6.85714 33.75C6.62381 33.5833 6.44881 33.3667 6.33214 33.1C6.21548 32.8333 6.19048 32.5333 6.25714 32.2L8.70714 21.6L0.507143 14.45C0.240476 14.2167 0.0821428 13.9583 0.0321429 13.675C-0.0178571 13.3917 -0.00952381 13.1167 0.0571429 12.85C0.12381 12.5833 0.27381 12.3583 0.507143 12.175C0.740476 11.9917 1.02381 11.8833 1.35714 11.85L12.2071 10.9L16.4071 0.9C16.5405 0.6 16.7405 0.375 17.0071 0.225C17.2738 0.075 17.5405 0 17.8071 0C18.0738 0 18.3405 0.075 18.6071 0.225C18.8738 0.375 19.0738 0.6 19.2071 0.9L23.4071 10.9L34.2571 11.85C34.5905 11.8833 34.8738 11.9917 35.1071 12.175C35.3405 12.3583 35.4905 12.5833 35.5571 12.85C35.6238 13.1167 35.6321 13.3917 35.5821 13.675C35.5321 13.9583 35.3738 14.2167 35.1071 14.45L26.9071 21.6L29.3571 32.2C29.4238 32.5333 29.3988 32.8333 29.2821 33.1C29.1655 33.3667 28.9905 33.5833 28.7571 33.75C28.5238 33.9167 28.2571 34.0083 27.9571 34.025C27.6571 34.0417 27.3738 33.9667 27.1071 33.8L17.8071 28.2Z"
                    className={
                        props.disabled
                            ? 'fill-light'
                            : props.active
                            ? 'fill-goldDark'
                            : 'fill-lightest'
                    }
                />
                <path
                    d="M17.8071 24.7L9.9571 29.4L12.0571 20.5L5.1571 14.5L14.2571 13.7L17.8071 5.34998L21.3571 13.75L30.4571 14.55L23.5571 20.55L25.6571 29.45L17.8071 24.7Z"
                    className={
                        props.disabled
                            ? 'fill-lightest'
                            : props.active
                            ? 'fill-gold'
                            : 'fill-background'
                    }
                />
            </svg>
        </button>
    );
};

export default Star;
