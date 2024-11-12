import { useState } from 'react';
import { putRequest } from '../../../api';
import { errorAlert } from '../../../util';
import { useAdminStore } from '../../../store';
import ConfirmPopup from '../../ConfirmPopup';
import RawTextInput from '../../RawTextInput';
import RawTextArea from '../../RawTextArea';

interface EditJudgePopupProps {
    /* Judge to edit */
    judge: Judge;

    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
}

const EditJudgePopup = (props: EditJudgePopupProps) => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const [name, setName] = useState(props.judge.name);
    const [email, setEmail] = useState(props.judge.email);
    const [notes, setNotes] = useState(props.judge.notes);

    const onSubmit = async () => {
        setIsSubmitting(true);

        const data = { name, email, notes };
        const res = await putRequest(`/judge/${props.judge.id}`, 'admin', data);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert('Judge updated successfully!');
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
            onSubmit={onSubmit}
            disabledSubmit={isSubmitting}
        >
            <h2 className="text-2xl font-bold mb-2 text-center text-primary">{props.judge.name}</h2>
            <div className="flex flex-row w-full my-3 space-x-3">
                <RawTextInput placeholder="Name" text={name} setText={setName} full />
                <RawTextInput placeholder="Email" text={email} setText={setEmail} full />
            </div>
            <RawTextArea placeholder="Notes (optional)" value={notes} setValue={setNotes} />
        </ConfirmPopup>
    );
};

export default EditJudgePopup;
