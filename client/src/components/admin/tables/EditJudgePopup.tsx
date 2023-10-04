import { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import { putRequest } from '../../../api';
import { errorAlert } from '../../../util';
import useAdminStore from '../../../store';

interface EditJudgePopupProps {
    /* Judge to edit */
    judge: Judge;

    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;
}

interface UpdateJudgeData {
    name: string;
    email: string;
    notes: string;
}

const EditJudgePopup = ({ judge, close }: EditJudgePopupProps) => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { register, handleSubmit, reset } = useForm<UpdateJudgeData>();
    const fetchJudges = useAdminStore((state) => state.fetchJudges);

    const onSubmit: SubmitHandler<UpdateJudgeData> = async (data) => {
        setIsSubmitting(true);

        const res = await putRequest(`/judge/${judge.id}`, 'admin', data);
        if (res.status !== 200) {
            errorAlert(res.status);
            return;
        }

        alert('Judge updated successfully!');
        reset();
        setIsSubmitting(false);
        close(false);
        fetchJudges();
    };

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] py-6 px-10 w-1/3">
                <form onSubmit={handleSubmit(onSubmit)}>
                    <h1 className="text-5xl font-bold text-center">Edit Judge</h1>
                    <h2 className="text-2xl font-bold mb-2 text-center text-primary">
                        {judge.name}
                    </h2>
                    <div className="flex flex-row w-full my-3 space-x-3">
                        <TextInput
                            name="name"
                            placeholder="Name"
                            register={register}
                            required
                            defaultValue={judge.name}
                        />
                        <TextInput
                            name="email"
                            placeholder="Email"
                            register={register}
                            required
                            defaultValue={judge.email}
                        />
                    </div>
                    <TextArea
                        name="notes"
                        placeholder="Notes (optional)"
                        register={register}
                        defaultValue={judge.notes}
                    />
                    <div className="flex flex-row justify-around">
                        <button
                            className=" border-lightest border-2 rounded-full px-6 py-1 mt-4 w-2/5 font-bold text-2xl text-lighter hover:bg-lighter/30 duration-200"
                            onClick={() => close(false)}
                        >
                            Cancel
                        </button>
                        <button
                            className="bg-primary text-white rounded-full px-4 py-2 mt-4 w-2/5 font-bold text-2xl hover:brightness-110 duration-200"
                            disabled={isSubmitting}
                        >
                            Submit
                        </button>
                    </div>
                </form>
            </div>
        </>
    );
};

export default EditJudgePopup;
