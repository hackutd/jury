import { twMerge } from 'tailwind-merge';

interface TextAreaProps {
    /* Label of the field */
    label: string;

    /* Placeholder for the field */
    placeholder?: string;

    /* State variable for the field */
    value: string;

    /* Setter function for the field */
    setValue: React.Dispatch<React.SetStateAction<string>>;

    /* Custom class name */
    className?: string;
}

const TextArea = (props: TextAreaProps) => {
    return (
        <div className={twMerge('flex flex-col', props.className)}>
            <label htmlFor={props.label} className="text-sm text-light mb-[2px]">
                {props.label}
            </label>
            <textarea
                className="w-full p-2 text-md border-lightest border-2 rounded-sm focus:border-primary focus:border-2 focus:outline-none"
                value={props.value}
                onChange={(e) => props.setValue(e.target.value)}
                placeholder={props.placeholder || props.label}
            />
        </div>
    );
};

export default TextArea;
