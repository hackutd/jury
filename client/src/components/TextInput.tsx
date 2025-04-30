import { twMerge } from 'tailwind-merge';

interface TextInputProps {
    /* Label for the field */
    label: string;

    /* Placeholder for the field */
    placeholder?: string;

    /* State variable for input */
    text: string | number;

    /* State setting function for input */
    setText:
        | React.Dispatch<React.SetStateAction<string>>
        | React.Dispatch<React.SetStateAction<number>>;

    /* Default value of the field */
    defaultValue?: string;

    /* Custom styling */
    className?: string;

    /* Full Width */
    full?: boolean;

    /* Large text */
    large?: boolean;

    /* Number input */
    number?: boolean;
}

const TextInput = (props: TextInputProps) => {
    return (
        <div
            className={twMerge(
                'flex flex-col w-auto',
                props.full && 'w-full',
                props.className
            )}
        >
            <label htmlFor={props.label} className="text-sm text-light mb-[2px]">
                {props.label}
            </label>
            <input
                id={props.label}
                type={props.number ? 'number' : 'text'}
                className={twMerge(
                    'w-full h-10 px-2 text-md rounded-sm border-lightest border-2 focus:border-primary outline-none',
                    props.large && 'text-xl h-12'
                )}
                placeholder={props.placeholder || props.label}
                defaultValue={props.defaultValue}
                value={props.text}
                onChange={(e) =>
                    props.setText(props.number ? Number(e.target.value) : (e.target.value as any))
                }
            />
        </div>
    );
};

export default TextInput;
