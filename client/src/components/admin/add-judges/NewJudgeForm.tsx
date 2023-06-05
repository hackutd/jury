import { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';

interface NewJudgeData {
    name: string;
    email: string;
    notes: string;
}

const NewJudgeForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { register, handleSubmit, reset } = useForm<NewJudgeData>();

    const onSubmit: SubmitHandler<NewJudgeData> = async (data) => {
        setIsSubmitting(true);

        // Upload judge
        try {
            const response = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/new`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
                credentials: 'include',
            });

            // Throw error if response is not ok
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            // Reset the form and show success message
            reset();
            alert(`Added judge ${data.name} successfully!`);
        } catch (err) {
            console.error(err);
            alert('Error adding judge: ' + err);
        }
        setIsSubmitting(false);
    };

    return (
        <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Judge</h1>
                <form className="flex flex-col w-full space-y-4" onSubmit={handleSubmit(onSubmit)}>
                    <div className="flex flex-row w-full mt-4 space-x-6">
                        <TextInput name="name" placeholder="Name" register={register} required />
                        <TextInput name="email" placeholder="Email" register={register} required />
                    </div>
                    <TextArea name="notes" placeholder="Notes (optional)" register={register} />
                    <button
                        className="w-full h-11 px-4 text-2xl text-white bg-primary rounded-full"
                        disabled={isSubmitting}
                    >
                        Add
                    </button>
                </form>
            </div>
        </div>
    );
};

export default NewJudgeForm;
