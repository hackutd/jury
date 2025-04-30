import { useState } from 'react';
import { postRequest } from '../../../api';
import { errorAlert } from '../../../util';
import { useAdminStore } from '../../../store';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import Button from '../../Button';
import Card from '../../Card';

const NewProjectForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const fetchProjectStats = useAdminStore((state) => state.fetchProjectStats);
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [url, setUrl] = useState('');
    const [tryLink, setTryLink] = useState('');
    const [videoLink, setVideoLink] = useState('');
    const [challengeList, setChallengeList] = useState('');

    const submit = async () => {
        if (isSubmitting) return;
        setIsSubmitting(true);

        const data = {
            name,
            description,
            url,
            try_link: tryLink,
            video_link: videoLink,
            challenge_list: challengeList,
        };

        const res = await postRequest('/project/new', 'admin', data);
        if (res.status !== 200) {
            errorAlert(res);
            setIsSubmitting(false);
            return;
        }

        alert('Project added successfully!');
        fetchProjectStats();
        setIsSubmitting(false);
    };

    return (
        <Card>
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Project</h1>
                <div className="flex flex-col w-full space-y-3">
                    <TextInput label="Name" text={name} setText={setName} />
                    <TextArea label="Description" value={description} setValue={setDescription} />
                    <TextInput label="URL" text={url} setText={setUrl} />
                    <div className="flex flex-col md:flex-row w-full mt-4 md:gap-6">
                        <TextInput
                            label='"Try It" Link (optional)'
                            text={tryLink}
                            setText={setTryLink}
                            full
                        />
                        <TextInput
                            label="Video Link (optional)"
                            text={videoLink}
                            setText={setVideoLink}
                            full
                        />
                    </div>
                    <TextInput
                        label="Challenge List (comma separated, optional)"
                        text={challengeList}
                        setText={setChallengeList}
                    />
                    <Button type="primary" onClick={submit} full>
                        Add
                    </Button>
                </div>
            </div>
        </Card>
    );
};

export default NewProjectForm;
