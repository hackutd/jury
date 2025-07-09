import { twMerge } from 'tailwind-merge';

interface CheckboxProps {
    /* Label for the checkbox */
    children: React.ReactNode;

    /* Whether the checkbox is checked */
    checked: boolean;

    /* Function to set state variable */
    onChange: (checked: boolean) => void;

    /* Disable the checkbox */
    disabled?: boolean;

    /* Optional class name */
    className?: string;
}

const Checkbox = (props: CheckboxProps) => {
    return (
        <div
            className={twMerge(
                'p-2 cursor-pointer text-light hover:text-black duration-200',
                props.disabled && 'cursor-not-allowed hover:text-light',
                props.className
            )}
            onClick={() => props.disabled || props.onChange(!props.checked)}
        >
            <input
                type="checkbox"
                checked={props.checked}
                disabled={props.disabled}
                onChange={(e) => {
                    props.disabled || props.onChange(e.target.checked);
                }}
                className={twMerge(
                    'mr-4 rounded-sm bg-white border-primary border-2 text-primary focus:ring-0',
                    props.disabled && 'text-lighter border-lighter'
                )}
            />
            <p className="inline">{props.children}</p>
        </div>
    );
};

export default Checkbox;
