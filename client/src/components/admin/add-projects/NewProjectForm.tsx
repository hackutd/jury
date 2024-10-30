import { SubmitHandler, useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import { useState } from 'react';
import { postRequest } from '../../../api';
import { errorAlert } from '../../../util';

interface NewProjectData {
    name: string;
    description: string;
    url: string;
    link: string;
    video: string;
    challenge_list: string;
}

const NewProjectForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { register, handleSubmit, reset } = useForm<NewProjectData>();

    const onSubmit: SubmitHandler<NewProjectData> = async (data) => {
        // Upload project
        setIsSubmitting(true);

        const res = await postRequest('/project/new', 'admin', data);
        if (res.status === 400) {
            alert(`Error with form: ${res.error}`);
            setIsSubmitting(false);
            return;
        }
        if (res.status !== 200) {
            errorAlert(res);
            setIsSubmitting(false);
            return;
        }

        alert('Project added successfully!');
        reset();
        setIsSubmitting(false);
    };

    return (
        <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Project</h1>
                <form className="flex flex-col w-full space-y-4" onSubmit={handleSubmit(onSubmit)}>
                    <TextInput name="name" placeholder="Name" register={register} />
                    <TextArea
                        name="description"
                        placeholder="Description"
                        register={register}
                    />
                    <TextInput name="url" placeholder="URL" register={register} />
                    <div className="flex flex-row w-full mt-4 space-x-6">
                        <TextInput
                            name="link"
                            placeholder='"Try It" Link (optional)'
                            register={register}
                        />
                        <TextInput
                            name="video"
                            placeholder="Video Link (optional)"
                            register={register}
                        />
                    </div>
                    <TextInput
                        name="challenge_list"
                        placeholder="Challenge List (comma separated, optional)"
                        register={register}
                    />
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

export default NewProjectForm;
