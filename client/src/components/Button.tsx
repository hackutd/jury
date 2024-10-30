import React from 'react';
import { twMerge } from 'tailwind-merge';

interface ButtonProps {
    /* Button internal content */
    children?: React.ReactNode;

    /* Should be defined if the button is a redirect */
    href?: string;

    /* If href is not defined, function should be written to handle button click */
    onClick?: (e: React.MouseEvent<Element>) => void;

    /* Type of the button */
    type: 'primary' | 'outline' | 'text' | 'error' | 'outline-primary' | 'gold';

    /* If true, sets button as a disabled button */
    disabled?: boolean;

    /* Show tooltip on hover */
    tooltip?: string;

    /* Bold button */
    bold?: boolean;

    /* Full width */
    full?: boolean;

    /* Small width */
    small?: boolean;

    /* Flat and long */
    flat?: boolean;

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
        'py-3 text-center text-2xl no-underline outline-none border-solid border-transparent border-[2.5px] bg-transparent text-primary hover:text-primaryDark rounded-xl';

    // Format borders
    const borderFormat =
        props.type.indexOf('outline') !== -1 &&
        'border-lightest hover:bg-backgroundDark text-light hover:text-light';
    const primaryBorderFormat =
        props.type === 'outline-primary' && 'border-primary text-primary hover:text-primary hover:bg-primary/20';

    // Format text
    const typeFormat =
        props.type === 'primary'
            ? 'bg-primary text-background hover:bg-primaryDark hover:text-background'
            : props.type === 'gold' ? 'bg-gold text-black hover:bg-goldDark hover:text-black'
            : props.type === 'error' &&
              'bg-error text-background hover:bg-errorDark hover:text-background';
    const varFormat = !props.disabled
        ? typeFormat + ' cursor-pointer duration-200'
        : 'cursor-auto text-lighter bg-backgroundDark hover:text-lighter';

    // Format bold and width
    const boldFormat = props.bold && 'font-bold text-xl';
    const widthFormat = props.full
        ? 'w-full'
        : props.small || props.flat
        ? 'w-auto px-6 py-1 rounded-md'
        : 'w-3/4 md:w-2/3';

    const flatFormat = props.flat && 'px-16';

    // Combine all formats
    const formatting = twMerge(
        defaultFormat,
        borderFormat,
        primaryBorderFormat,
        varFormat,
        boldFormat,
        widthFormat,
        flatFormat,
        props.className
    );

    // Disable button bc links cannot be disabled
    return props.disabled ? (
        <button disabled className={formatting} title={props.tooltip}>
            {props.children}
        </button>
    ) : props.href ? (
        <a href={props.href || ''} className={formatting + ' block'} title={props.tooltip}>
            {props.children}
        </a>
    ) : (
        <button className={formatting} onClick={props.onClick} title={props.tooltip}>
            {props.children}
        </button>
    );
};

export default Button;
