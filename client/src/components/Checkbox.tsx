interface CheckboxProps {
    /* Label for the checkbox */
    children: React.ReactNode;

    /* Whether the checkbox is checked */
    checked: boolean;

    /* Function to set state variable */
    onChange: (checked: boolean) => void;
}

const Checkbox = (props: CheckboxProps) => {
    return (
        <div
            className="p-2 cursor-pointer text-light hover:text-black duration-200"
            onClick={() => props.onChange(!props.checked)}
        >
            <input
                type="checkbox"
                checked={props.checked}
                onChange={(e) => {
                    props.onChange(e.target.checked);
                }}
                className="mr-4 rounded-sm bg-white border-primary border-2 text-primary focus:ring-0"
            />
            <p className="inline">{props.children}</p>
        </div>
    );
};

export default Checkbox;
