import { deleteRequest } from '../../../api';
import { useAdminStore } from '../../../store';
import { errorAlert } from '../../../util';
import ConfirmPopup from '../../ConfirmPopup';
import Popup from '../../Popup';

type DeleteElement = Project | Judge;

interface DeletePopupProps {
    /* Element to delete */
    element: DeleteElement;

    /* Enabled state variable */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
}

function isProject(e: DeleteElement): e is Project {
    return 'mu' in e;
}

const DeletePopup = (props: DeletePopupProps) => {
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);

    const deleteElement = async () => {
        const resource = isProject(props.element) ? 'project' : 'judge';
        const res = await deleteRequest(`/${resource}/${props.element.id}`, 'admin');
        if (res.status === 200) {
            fetchStats();
            isProject(props.element) ? fetchProjects() : fetchJudges();
            alert(`${resource} deleted successfully!`);
        } else {
            errorAlert(res);
        }

        props.setEnabled(false);
    };
    return (
        <ConfirmPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            title="Heads Up!"
            submitText="Delete"
            onSubmit={deleteElement}
            red
        >
            <p className="text-xl">
                Are you sure you want to delete{' '}
                <span className="text-primary font-bold">{props.element.name}</span>? This action is
                permanent and cannot be undone.
            </p>
        </ConfirmPopup>
    );
};

export default DeletePopup;
