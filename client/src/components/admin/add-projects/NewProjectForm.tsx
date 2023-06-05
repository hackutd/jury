import { SubmitHandler, useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import { useState } from 'react';

interface NewProjectData {
    name: string;
    description: string;
    link: string;
    video: string;
    challenge_list: string;
}

interface ProjectUpload {
    name: string;
    description: string;
    link: string;
    video: string;
    challenge_list: string[];
}

const NewProjectForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { register, handleSubmit, reset } = useForm<NewProjectData>();

    const onSubmit: SubmitHandler<NewProjectData> = async (data) => {
        const project: ProjectUpload = {
            name: data.name,
            description: data.description,
            link: data.link,
            video: data.video,
            challenge_list: data.challenge_list.split(',').map((x) => x.trim()),
        };

        // Upload project
        setIsSubmitting(true);
        try {
            const response = await fetch(`${process.env.REACT_APP_JURY_URL}/project/new`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(project),
                credentials: 'include',
            });

            // Throw error if response is not ok
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            // Reset the form and show success message
            reset();
            alert(`Added project ${data.name} successfully!`);
        } catch (err) {
            console.error(err);
            alert('Error adding project: ' + err);
        }
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
                        placeholder="Description (optional)"
                        register={register}
                    />
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
