import { twMerge } from 'tailwind-merge';

interface RawTextInputProps {
    /* Name of the field */
    name: string;

    /* Placeholder of the field */
    placeholder: string;

    /* State variable for input */
    text: string;

    /* State setting function for input */
    setText: React.Dispatch<React.SetStateAction<string>>;

    /* Default value of the field */
    defaultValue?: string;

    /* Custom styling */
    className?: string;

    /* Full Width */
    full?: boolean;

    /* Large text */
    large?: boolean;
}

const RawTextInput = (props: RawTextInputProps) => {
    return (
        <input
            className={twMerge(
                'w-auto h-10 px-2 text-md rounded-sm border-lightest border-2 focus:border-primary outline-none',
                props.full && 'w-full',
                props.large && 'text-xl h-12',
                props.className
            )}
            placeholder={props.placeholder}
            defaultValue={props.defaultValue}
            value={props.text}
            onChange={(e) => props.setText(e.target.value)}
        />
    );
};

export default RawTextInput;
