import React from 'react';

interface ButtonProps {
    /* Button internal content */
    children?: React.ReactNode;

    /* Should be defined if the button is a redirect */
    href?: string;

    /* If href is not defined, function should be written to handle button click */
    onClick?: (e: React.MouseEvent<any>) => void;

    /* Type of the button */
    type: 'primary' | 'text';

    /* If true, sets button as a disabled button */
    disabled?: boolean;

    /* Classname styling */
    className?: string;
}

/**
 * This is a button component that displays a clickable button.
 * The button can be disabled, which will grey out the button and make it unclickable.
 * This component can be used as a link (with the href prop) or a button (with the onClick prop).
 */
const Button = (props: ButtonProps) => {
    // Define formatting
    const defaultFormat =
        'rounded-full py-4 w-3/4 text-center text-2xl no-underline outline-none border-none md:w-2/3 ';
    const typeFormat =
        props.type === 'primary' ? 'bg-primary text-background' : 'bg-transparent text-primary';
    const varFormat = !props.disabled
        ? typeFormat + ' cursor-pointer duration-200 hover:scale-110 focus:scale-110'
        : 'cursor-auto text-lighter bg-backgroundDark';
    const formatting = defaultFormat + ' ' + varFormat;

    // Return a disabled button or link if not disabled
    return props.disabled ? (
        <button disabled className={formatting}>
            {props.children}
        </button>
    ) : props.href ? (
        <a href={props.href || ''} className={formatting}>
            {props.children}
        </a>
    ) : (
        <button className={formatting} onClick={props.onClick}>
            {props.children}
        </button>
    );
};

export default Button;
