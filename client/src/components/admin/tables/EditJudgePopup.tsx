import { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import { putRequest } from '../../../api';
import { errorAlert } from '../../../util';
import { useAdminStore } from '../../../store';
import ConfirmPopup from '../../ConfirmPopup';
import Loading from '../../Loading';

interface EditJudgePopupProps {
    /* Judge to edit */
    judge: Judge;

    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
}

interface UpdateJudgeData {
    name: string;
    email: string;
    notes: string;
}

const EditJudgePopup = (props: EditJudgePopupProps) => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { register, handleSubmit, reset } = useForm<UpdateJudgeData>();
    const fetchJudges = useAdminStore((state) => state.fetchJudges);

    const onSubmit: SubmitHandler<UpdateJudgeData> = async (data) => {
        setIsSubmitting(true);

        const res = await putRequest(`/judge/${props.judge.id}`, 'admin', data);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert('Judge updated successfully!');
        reset();
        setIsSubmitting(false);
        props.setEnabled(false);
        fetchJudges();
    };

    return (
        <ConfirmPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            title="Edit Judge"
            submitText="Save"
            onSubmit={handleSubmit(onSubmit)}
            disabledSubmit={isSubmitting}
        >
            <h2 className="text-2xl font-bold mb-2 text-center text-primary">{props.judge.name}</h2>
            <div className="flex flex-row w-full my-3 space-x-3">
                <TextInput
                    name="name"
                    placeholder="Name"
                    register={register}
                    required
                    defaultValue={props.judge.name}
                />
                <TextInput
                    name="email"
                    placeholder="Email"
                    register={register}
                    required
                    defaultValue={props.judge.email}
                />
            </div>
            <TextArea
                name="notes"
                placeholder="Notes (optional)"
                register={register}
                defaultValue={props.judge.notes}
            />
        </ConfirmPopup>
    );
};

export default EditJudgePopup;
