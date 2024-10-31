import { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import { postRequest } from '../../../api';
import { errorAlert } from '../../../util';
import Checkbox from '../../Checkbox';
import Loading from '../../Loading';
import { useAdminStore } from '../../../store';

interface NewJudgeData {
    name: string;
    email: string;
    track: string;
    notes: string;
}

type NewJudgeDataFull = NewJudgeData & { no_send: boolean };

const NewJudgeForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { register, handleSubmit, reset } = useForm<NewJudgeData>();
    const [noSend, setNoSend] = useState(false);
    const fetchJudgeStats = useAdminStore((state) => state.fetchJudgeStats);

    const onSubmit: SubmitHandler<NewJudgeData> = async (data) => {
        setIsSubmitting(true);

        const newdata: NewJudgeDataFull = { ...data, no_send: noSend };

        const res = await postRequest('/judge/new', 'admin', newdata);
        if (res.status !== 200) {
            errorAlert(res);
            setIsSubmitting(false);
            return;
        }

        alert('Judge added successfully!');
        reset();
        fetchJudgeStats();
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
                    <TextInput
                        name="track"
                        placeholder="Track (optional, general judge if not specified)"
                        register={register}
                    />
                    <TextArea name="notes" placeholder="Notes (optional)" register={register} />
                    <Checkbox checked={noSend} onChange={setNoSend}>
                        Do not send an email
                    </Checkbox>
                    <button
                        className="w-full h-11 px-4 text-2xl text-white bg-primary rounded-full"
                        disabled={isSubmitting}
                    >
                        Add
                    </button>
                </form>
            </div>
            <Loading disabled={!isSubmitting} />
        </div>
    );
};

export default NewJudgeForm;
