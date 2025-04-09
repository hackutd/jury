import { useEffect, useRef } from 'react';
import { twMerge } from 'tailwind-merge';

interface ActionsDropdownProps {
    /* State variable to open or close the dropdown */
    open: boolean;

    /* Function to close the dropdown */
    setOpen: React.Dispatch<React.SetStateAction<boolean>>;

    /* List of actions to be displayed in the dropdown */
    actions: string[];

    /* List of functions for the action in the list */
    actionFunctions: (() => void)[];

    /* Indices to make red */
    redIndices?: number[];

    /* Large text */
    large?: boolean;

    /* Additional class names */
    className?: string;
}

const ActionsDropdown = (props: ActionsDropdownProps) => {
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        function closeClick(event: MouseEvent) {
            if (ref && ref.current && !ref.current.contains(event.target as Node)) {
                props.setOpen(false);
            }
        }

        // Bind the event listener
        document.addEventListener('mousedown', closeClick);
        return () => {
            // Unbind the event listener on clean up
            document.removeEventListener('mousedown', closeClick);
        };
    }, [ref]);

    const handleClick = (index: number) => {
        props.actionFunctions[index]();
        props.setOpen(false);
    };

    if (!props.open) {
        return null;
    }

    return (
        <div
            className={twMerge(
                'absolute flex flex-col bg-background rounded-md border-lightest border-2 font-normal text-sm',
                props.large && 'text-lg',
                props.className
            )}
            ref={ref}
        >
            {props.actions.map((action, index) => (
                <div
                    key={index}
                    className={twMerge(
                        'py-1 px-2 cursor-pointer hover:bg-primary/20 duration-150',
                        props.redIndices?.includes(index) && 'text-error'
                    )}
                    onClick={handleClick.bind(null, index)}
                >
                    {action}
                </div>
            ))}
        </div>
    );
};

export default ActionsDropdown;
