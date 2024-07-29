import { twMerge } from 'tailwind-merge';
import Button from './Button';

export interface RadioOption {
    /* Internal value of radio option */
    value: string;

    /* Display text of radio option */
    title: string;

    /* Subtitle display text of radio option */
    subtitle?: string;
}

interface RadioButtonProps {
    /* Radio button content */
    option: RadioOption;

    /* State variable for selection */
    selected: boolean;

    /* Function to run on click */
    onClick: (e: React.MouseEvent<Element>) => void;

    /* Color of the popup -- use colors defined in tailwind config */
    color: string;
}

/**
 * Component used in a RadioSelect component. The design of this button is so that
 * we can have an array of values (hence onClick instead of the state setter function).
 * See the RadioSelect component; this component shouldn't be directly used.
 */
const RadioButton = (props: RadioButtonProps) => {
    return (
        <Button
            type="outline"
            full
            square
            className={twMerge(
                'flex flex-col items-center py-2 px-4 my-2',
                props.selected
                    ? `text-${props.color} bg-${props.color}/10 border-${props.color}`
                    : 'text-black'
            )}
            onClick={props.onClick}
        >
            <h3 className="text-2xl">{props.option.title}</h3>
            {props.option.subtitle && <p className="text-lg text-light">{props.option.subtitle}</p>}
        </Button>
    );
};

export default RadioButton;
