import { twMerge } from 'tailwind-merge';

interface SelectionButtonProps {
    /* Button internal content */
    children?: React.ReactNode;

    /* Should be defined if the button is a redirect */
    href?: string;

    /* Function to run when button is clicked */
    onClick?: (newSelection: string) => void;

    /* State variable of selection */
    selected: string;

    /* Function to set state variable of selection */
    setSelected: React.Dispatch<React.SetStateAction<string>>;

    /* List of options for selection */
    options: string[];

    /* Disabled button */
    disabled?: boolean;

    /* Classname styling */
    className?: string;
}

const SelectionButton = (props: SelectionButtonProps) => {
    return (
        <div className={twMerge('flex flex-row space-x-[2px] bg-lightest border-solid border-2 border-lightest rounded-md overflow-hidden text-xl w-max', props.className)}>
            {props.options.map((option, index) => (
                <button
                    disabled={props.disabled}
                    key={index}
                    className={twMerge(
                        'py-1 px-4 bg-background hover:bg-backgroundDark',
                        props.selected === option &&
                            'bg-primary text-background hover:bg-primaryDark',
                        props.disabled && 'cursor-auto text-lighter bg-backgroundDark hover:text-lighter hover:bg-backgroundDark'
                    )}
                    onClick={() => {
                        props.setSelected(option);
                        props.onClick && props.onClick(option);
                    }}
                >
                    {option}
                </button>
            ))}
        </div>
    );
};

export default SelectionButton;
