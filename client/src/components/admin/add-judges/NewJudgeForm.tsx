import { useState } from 'react';
import { postRequest } from '../../../api';
import { errorAlert } from '../../../util';
import Checkbox from '../../Checkbox';
import Loading from '../../Loading';
import { useAdminStore, useOptionsStore } from '../../../store';
import TextArea from '../../TextArea';
import TextInput from '../../TextInput';
import Dropdown from '../../Dropdown';
import Button from '../../Button';

const NewJudgeForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [track, setTrack] = useState('');
    const [notes, setNotes] = useState('');
    const [noSend, setNoSend] = useState(false);
    const fetchJudgeStats = useAdminStore((state) => state.fetchJudgeStats);
    const options = useOptionsStore((state) => state.options);

    const submit = async () => {
        setIsSubmitting(true);

        // Make sure name and email are defined
        if (!name || !email) {
            alert('Name and email are required');
            setIsSubmitting(false);
            return;
        }

        const newData = { name, email, track, notes, no_send: noSend };

        const res = await postRequest('/judge/new', 'admin', newData);
        if (res.status !== 200) {
            errorAlert(res);
            setIsSubmitting(false);
            return;
        }

        alert('Judge added successfully!');
        resetForm();
        fetchJudgeStats();
        setIsSubmitting(false);
    };

    const resetForm = () => {
        setName('');
        setEmail('');
        setTrack('');
        setNotes('');
        setNoSend(false);
    };

    return (
        <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Judge</h1>
                <div className="flex flex-col w-full">
                    <div className="flex flex-row w-full mt-4 space-x-6">
                        <TextInput placeholder="Name" text={name} setText={setName} large full />
                        <TextInput placeholder="Email" text={email} setText={setEmail} large full />
                    </div>
                    <p className="mt-2 font-bold text-lighter">Track</p>
                    <Dropdown
                        options={['', ...options.tracks]}
                        selected={track}
                        setSelected={setTrack}
                        className="mt-0 mb-4 text-left text-xl bg-white rounded-sm"
                    />
                    <TextArea
                        value={notes}
                        setValue={setNotes}
                        placeholder="Notes (optional)"
                        className="text-xl"
                    />
                    <Checkbox checked={noSend} onChange={setNoSend} className="my-2">
                        Do not send an email
                    </Checkbox>
                    <Button type="primary" onClick={submit} full flat className="py-1 rounded-md">
                        Add
                    </Button>
                </div>
            </div>
            <Loading disabled={!isSubmitting} />
        </div>
    );
};

export default NewJudgeForm;
