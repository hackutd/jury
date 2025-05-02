import { useState } from 'react';
import ConfirmPopup from '../../ConfirmPopup';
import { useAdminStore } from '../../../store';
import TextInput from '../../TextInput';
import TextArea from '../../TextArea';
import { putRequest } from '../../../api';
import { errorAlert } from '../../../util';

interface EditProjectPopupProps {
    /* Project to edit */
    project: Project;

    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
}

const EditProjectPopup = (props: EditProjectPopupProps) => {
    const [isSubmitting, setIsSubmitting] = useState(false);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const [name, setName] = useState(props.project.name);
    const [description, setDescription] = useState(props.project.description);
    const [challengeList, setChallengeList] = useState(props.project.challenge_list.join(', '));
    const [url, setUrl] = useState(props.project.url);
    const [tryLink, setTryLink] = useState(props.project.try_link);
    const [videoLink, setVideoLink] = useState(props.project.video_link);

    const onSubmit = async () => {
        setIsSubmitting(true);

        const data = {
            name,
            description,
            challenge_list: challengeList,
            url,
            try_link: tryLink,
            video_link: videoLink,
        };
        const res = await putRequest(`/project/${props.project.id}`, 'admin', data);
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert('Project updated successfully!');
        setIsSubmitting(false);
        props.setEnabled(false);
        fetchProjects();
    };

    return (
        <ConfirmPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            title="Edit Project"
            submitText="Save"
            onSubmit={onSubmit}
        >
            <h2 className="text-2xl font-bold mb-2 text-center text-primary">
                {props.project.name}
            </h2>
            <div className="flex flex-col gap-3">
                <TextInput label="Name" text={name} setText={setName} full />
                <TextInput
                    label="Challenge List"
                    text={challengeList}
                    setText={setChallengeList}
                    full
                />
                <TextInput label="URL" text={url} setText={setUrl} full />
                <div className="flex flex-row w-full space-x-3">
                    <TextInput label="Try Link" text={tryLink} setText={setTryLink} full />
                    <TextInput
                        label="Video Link"
                        text={videoLink}
                        setText={setVideoLink}
                        full
                    />
                </div>
                <TextArea label="Description" value={description} setValue={setDescription} />
            </div>
        </ConfirmPopup>
    );
};

export default EditProjectPopup;
