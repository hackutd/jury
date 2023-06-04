import { useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';

const NewProjectForm = () => {
    const { register, handleSubmit } = useForm();

    const onSubmit = (data: any) => {
        console.log(data);
    };

    return (
        <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Project</h1>
                <form className="flex flex-col w-full space-y-4" onSubmit={handleSubmit(onSubmit)}>
                    <TextInput name="name" placeholder="Name" register={register} />
                    <TextArea name="description" placeholder="Description (optional)" register={register} />
                    <div className="flex flex-row w-full mt-4 space-x-6">
                        <TextInput name="link" placeholder='"Try It" Link (optional)' register={register} />
                        <TextInput name="video" placeholder="Video Link (optional)" register={register} />
                    </div>
                    <TextInput
                        name="challenge-list"
                        placeholder="Challenge List (comma separated, optional)"
                        register={register}
                    />
                    <button className="w-full h-11 px-4 text-2xl text-white bg-primary rounded-full">
                        Add
                    </button>
                </form>
            </div>
        </div>
    );
};

export default NewProjectForm;
