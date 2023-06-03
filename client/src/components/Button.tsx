import React from 'react';
import { twMerge } from 'tailwind-merge';

interface ButtonProps {
    /* Button internal content */
    children?: React.ReactNode;

    /* Should be defined if the button is a redirect */
    href?: string;

    /* If href is not defined, function should be written to handle button click */
    onClick?: (e: React.MouseEvent<any>) => void;

    /* Type of the button */
    type: 'primary' | 'outline' | 'text' | 'error';

    /* If true, sets button as a disabled button */
    disabled?: boolean;

    /* Square button */
    square?: boolean;

    /* Bold button */
    bold?: boolean;

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
    const defaultFormat = 'py-4 w-3/4 text-center text-2xl no-underline outline-none md:w-2/3 ';
    const borderFormat = props.type === 'outline' ? 'border-lightest border-[3px]' : 'border-none';
    const typeFormat =
        props.type === 'primary'
            ? 'bg-primary text-background'
            : props.type === 'error'
            ? 'bg-error text-background'
            : 'bg-transparent text-primary';
    const squareFormat = props.square ? 'rounded-lg' : 'rounded-full';
    const varFormat = !props.disabled
        ? typeFormat + ' cursor-pointer duration-200 hover:scale-110 focus:scale-110'
        : 'cursor-auto text-lighter bg-backgroundDark';
    const boldFormat = props.bold ? 'font-bold' : 'font-normal';
    const formatting = twMerge(
        defaultFormat,
        borderFormat,
        varFormat,
        squareFormat,
        boldFormat,
        props.className
    );

    // Disable button bc links cannot be disabled
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
