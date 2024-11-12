import { useState } from 'react';
import { postRequest } from '../../../api';
import { errorAlert } from '../../../util';
import { useAdminStore } from '../../../store';
import RawTextInput from '../../RawTextInput';
import RawTextArea from '../../RawTextArea';

const NewProjectForm = () => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const fetchProjectStats = useAdminStore((state) => state.fetchProjectStats);
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [url, setUrl] = useState('');
    const [tryLink, setTryLink] = useState('');
    const [videoLink, setVideoLink] = useState('');
    const [challengeList, setChallengeList] = useState('');

    const handleSubmit: React.FormEventHandler<HTMLFormElement> = async (e) => {
        e.preventDefault();
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
        <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Project</h1>
                <form className="flex flex-col w-full space-y-4" onSubmit={handleSubmit}>
                    <RawTextInput placeholder="Name" text={name} setText={setName} />
                    <RawTextArea
                        placeholder="Description"
                        value={description}
                        setValue={setDescription}
                    />
                    <RawTextInput placeholder="URL" text={url} setText={setUrl} />
                    <div className="flex flex-row w-full mt-4 space-x-6">
                        <RawTextInput
                            placeholder='"Try It" Link (optional)'
                            text={tryLink}
                            setText={setTryLink}
                            full
                        />
                        <RawTextInput
                            placeholder="Video Link (optional)"
                            text={videoLink}
                            setText={setVideoLink}
                            full
                        />
                    </div>
                    <RawTextInput
                        placeholder="Challenge List (comma separated, optional)"
                        text={challengeList}
                        setText={setChallengeList}
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
