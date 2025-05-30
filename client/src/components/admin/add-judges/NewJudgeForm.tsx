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
import Card from '../../Card';

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
        <Card>
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Judge</h1>
                <div className="flex flex-col w-full gap-2">
                    <div className="flex flex-col md:flex-row w-full mt-4 gap-2 md:gap-6">
                        <TextInput label="Name" text={name} setText={setName} full />
                        <TextInput label="Email" text={email} setText={setEmail} full />
                    </div>
                    <div className="w-full">
                        <p className="text-light text-sm mb-1">Track</p>
                        <Dropdown
                            options={['', ...options.tracks]}
                            selected={track}
                            setSelected={setTrack}
                            className="bg-white rounded-sm hover:bg-background w-full"
                        />
                    </div>
                    <TextArea
                        value={notes}
                        setValue={setNotes}
                        label="Notes (optional)"
                        className="text-xl"
                    />
                    <Checkbox checked={noSend} onChange={setNoSend}>
                        Do not send an email
                    </Checkbox>
                    <Button type="primary" onClick={submit} full className="py-1 rounded-md">
                        Add
                    </Button>
                </div>
            </div>
            <Loading disabled={!isSubmitting} />
        </Card>
    );
};

export default NewJudgeForm;
