import type { UseFormRegister, FieldValues } from 'react-hook-form';

interface TextInputProps {
    /* Name of the field */
    name: string;

    /* Placeholder of the field */
    placeholder: string;

    /* Register function from react-hook-form */
    register: UseFormRegister<FieldValues>;
}

const TextInput = (props: TextInputProps) => {
    return (
        <input
            className="w-full h-14 px-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
            placeholder={props.placeholder}
            {...props.register(props.name)}
        />
    );
};

export default TextInput;
