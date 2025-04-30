import React from 'react';
import { twMerge } from 'tailwind-merge';

interface DropdownProps {
    /* List of options to display in the dropdown */
    options: string[];

    /* The selected option */
    selected: string;

    /* State function to set the selected option */
    setSelected: React.Dispatch<React.SetStateAction<string>> | ((selected: string) => void);

    /* OnChange function */
    onChange?: (() => void) | ((track: string) => void);

    /* Large text */
    large?: boolean;

    /* Classname styling */
    className?: string;
}

const Dropdown = (props: DropdownProps) => {
    return (
        <select
            className={twMerge(
                'border-[2.5px] pl-2 border-solid border-lightest outline-none focus:outline-none rounded-md bg-background cursor-pointer hover:bg-backgroundDark duration-150 text-light',
                props.large && 'text-xl',
                props.className
            )}
            onChange={(e) => {
                props.setSelected(e.target.value);
                if (props.onChange) props.onChange(e.target.value);
            }}
            value={props.selected}
        >
            {props.options.map((option) => (
                <option key={option} value={option}>
                    {option}
                </option>
            ))}
        </select>
    );
};

export default Dropdown;
