import type { UseFormRegister } from 'react-hook-form';

interface TextInputProps {
    /* Name of the field */
    name: string;

    /* Placeholder of the field */
    placeholder: string;

    /* Register function from react-hook-form */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    register: UseFormRegister<any>;

    /* Required field */
    required?: boolean;
}

const TextInput = (props: TextInputProps) => {
    return (
        <input
            className="w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
            placeholder={props.placeholder}
            {...props.register(props.name, { required: props.required || false })}
        />
    );
};

export default TextInput;
