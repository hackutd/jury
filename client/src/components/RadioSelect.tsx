import RadioButton, { RadioOption } from './RadioButton';

interface RadioSelectProps {
    /* List of radio button options */
    options: RadioOption[];

    /* Color of the radio buttons */
    color: string;

    /* State variable for which option is selected, "value" field of the options */
    selected: string;

    /* Function to modify the selected state variable */
    setSelected: React.Dispatch<React.SetStateAction<string>>;
}

/* Shows a list of buttons, where only one can be selected. */
const RadioSelect = (props: RadioSelectProps) => {
    // Function to run when item is selected
    const onClick = (i: number) => {
        props.setSelected(props.options[i].value);
    };

    return (
        <div>
            {props.options.map((v, i) => (
                <RadioButton
                    key={i}
                    color={props.color}
                    onClick={onClick.bind(this, i)}
                    option={v}
                    selected={props.selected === v.value}
                />
            ))}
        </div>
    );
};

export default RadioSelect;
