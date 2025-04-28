import React, { useEffect } from 'react';
import { twMerge } from 'tailwind-merge';

interface PopupProps {
    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* React children, corresponds to the body content */
    children?: React.ReactNode;

    /* Optional classname style to apply to outer div of popup */
    className?: string;
}

/**
 * Generic popup component with a backdrop and popup modal. Create a state variable and pass in the variable and setter function.
 * Clicking on the backdrop will close the popup.
 */
const Popup = (props: PopupProps) => {
    useEffect(() => {
        // handle esc key press
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === 'Escape') {
                props.setEnabled(false);
            }
        };

        if (props.enabled) {
            window.addEventListener('keydown', handleKeyDown);
        }

        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [props.enabled, props.setEnabled]);

    if (!props.enabled) {
        return null;
    }

    return (
        <>
            <div
                className="fixed left-0 top-0 mt-0 mb-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => props.setEnabled(false)}
            ></div>
            <div
                className={twMerge(
                    'fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] px-6 py-6 md:px-10 md:w-1/3 w-11/12 rounded-md border-2 border-lighter border-solid flex flex-col bg-background',
                    props.className
                )}
            >
                {/* x to close the popup */}
                <button
                    onClick={() => props.setEnabled(false)}
                    className="absolute top-1 left-3 text-3xl text-lightest hover:text-black duration-150"
                    aria-label="Close popup"
                >
                    &times;
                </button>
                {props.children}
            </div>
        </>
    );
};

export default Popup;
