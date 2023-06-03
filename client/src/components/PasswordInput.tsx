import { useState } from 'react';
import { twMerge } from 'tailwind-merge';

interface PasswordInputProps {
    /* Set true if text input has errored */
    error?: boolean;

    /* Setter for error state; will reset if no longer errored */
    setError?: React.Dispatch<React.SetStateAction<boolean>>;

    /* Error message to show if text incorrect */
    errorMessage?: string;

    /* Max Length of text field; null for no max */
    maxLength?: number;

    /* Placeholder of text field */
    placeholder?: string;

    /* Label under the field */
    label: string;

    /* Classname styling */
    className?: string;

    /* Handler for text input onChange */
    onChange?: React.ChangeEventHandler<HTMLInputElement>;

    /* Handler for text input onKeyPress */
    onKeyPress?: React.KeyboardEventHandler<HTMLInputElement>;

    /* Set default value of input field */
    value?: string;

    /* True if it's a password field */
    isHidden?: boolean;
}

const PasswordInput = (props: PasswordInputProps) => {
    const [focused, setFocused] = useState(false);

    const handleChange = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (props.onKeyPress) props.onKeyPress(e);
        if (props.setError) props.setError(false);
    };

    return (
        <>
            <input
                type={props.isHidden ? 'password' : 'text'}
                maxLength={props.maxLength}
                id="text-input"
                placeholder={props.placeholder}
                className={twMerge(
                    'block text-4xl font-mono font-normal text-black bg-transparent outline-none border-0 border-b-[3px] border-solid p-0 duration-200',
                    props.error ? 'border-error ' : 'border-light focus:border-primary ',
                    props.className
                )}
                onFocus={() => {
                    setFocused(true);
                }}
                onBlur={() => {
                    setFocused(false);
                }}
                onChange={props.onChange}
                onKeyDown={handleChange}
            />
            <label
                htmlFor="text-input"
                id="text-input-label"
                className={twMerge(
                    'mt-4 text-2xl duration-200',
                    props.error ? 'text-error' : focused ? 'text-primary' : ''
                )}
            >
                {props.error ? (props.errorMessage || props.label) : props.label}
            </label>
        </>
    );
};

export default PasswordInput;
