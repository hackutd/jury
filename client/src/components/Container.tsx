import React from 'react';
import { twMerge } from 'tailwind-merge';

interface ContainerProps {
    /* Element children */
    children?: React.ReactNode;

    /* Tailwind classes, will override defaults */
    className?: string;

    /* If true, don't justify center the entire page */
    noCenter?: boolean;
}

/**
 * This is the main container for the site, which is optimized for mobile.
 * On desktop devices, a set width will be imposed to mimic mobile displays.
 */
const Container = (props: ContainerProps) => {
    return (
        <div
            className={twMerge(
                'flex flex-col grow shrink-0 my-0 md:mx-auto md:w-[30rem] h-full',
                props.noCenter? '' : 'justify-center items-center',
                props.className
            )}
        >
            {props.children}
        </div>
    );
};

export default Container;
