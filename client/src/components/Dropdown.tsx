import React from 'react';
import { twMerge } from 'tailwind-merge';

interface DropdownProps {
    /* List of options to display in the dropdown */
    options: string[];

    /* The selected option */
    selected: string;

    /* State function to set the selected option */
    setSelected: React.Dispatch<React.SetStateAction<string>> | ((selected: string) => void);

    /* Classname styling */
    className?: string;
}

const Dropdown = (props: DropdownProps) => {
    return (
        <select
            className={twMerge(
                'border-[2.5px] border-solid border-lightest outline-none focus:outline-none rounded-md bg-background cursor-pointer hover:bg-backgroundDark duration-150 text-xl text-center text-light',
                props.className
            )}
            onChange={(e) => props.setSelected(e.target.value)}
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
