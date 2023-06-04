import type { UseFormRegister } from 'react-hook-form';

interface TextAreaProps {
    /* Name of the field */
    name: string;

    /* Placeholder of the field */
    placeholder: string;

    /* Register function from react-hook-form */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    register: UseFormRegister<any>;
}

const TextArea = (props: TextAreaProps) => {
    return (
        <textarea
            className="w-full h-36 px-4 py-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
            placeholder={props.placeholder}
            {...props.register(props.name)}
        />
    );
};

export default TextArea;
