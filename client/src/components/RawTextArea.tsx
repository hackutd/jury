import { twMerge } from 'tailwind-merge';

interface RawTextAreaProps {
    /* Placeholder of the field */
    placeholder: string;

    /* State variable for the field */
    value: string;

    /* Setter function for the field */
    setValue: React.Dispatch<React.SetStateAction<string>>;

    /* Custom class name */
    className?: string;
}

const RawTextArea = (props: RawTextAreaProps) => {
    return (
        <textarea
            className={twMerge(
                'w-full p-2 text-md border-lightest border-2 rounded-sm focus:border-primary focus:border-2 focus:outline-none',
                props.className
            )}
            value={props.value}
            onChange={(e) => props.setValue(e.target.value)}
            placeholder={props.placeholder}
        />
    );
};

export default RawTextArea;
