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
}

const RawTextInput = (props: RawTextInputProps) => {
    return (
        <input
            className={twMerge(
                'w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none',
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
