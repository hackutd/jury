import { twMerge } from 'tailwind-merge';

interface RawTextInputProps {
    /* Placeholder of the field */
    placeholder: string;

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

const RawTextInput = (props: RawTextInputProps) => {
    return (
        <input
            type={props.number ? 'number' : 'text'}
            className={twMerge(
                'w-auto h-10 px-2 text-md rounded-sm border-lightest border-2 focus:border-primary outline-none',
                props.full && 'w-full',
                props.large && 'text-xl h-12',
                props.className
            )}
            placeholder={props.placeholder}
            defaultValue={props.defaultValue}
            value={props.text}
            onChange={(e) =>
                props.setText(props.number ? Number(e.target.value) : (e.target.value as any))
            }
        />
    );
};

export default RawTextInput;
