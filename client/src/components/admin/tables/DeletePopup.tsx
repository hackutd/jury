import { deleteRequest } from '../../../api';
import { useAdminStore } from '../../../store';
import { errorAlert } from '../../../util';

type DeleteElement = Project | Judge;

interface DeletePopupProps {
    /* Element to delete */
    element: DeleteElement;

    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;
}

function isProject(e: DeleteElement): e is Project {
    return 'mu' in e;
}

const DeletePopup = ({ element, close }: DeletePopupProps) => {
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);

    const deleteElement = async () => {
        const resource = isProject(element) ? 'project' : 'judge';
        const res = await deleteRequest(`/${resource}/${element.id}`, 'admin');
        if (res.status === 200) {
            fetchStats();
            isProject(element) ? fetchProjects() : fetchJudges();
            alert(`${resource} deleted successfully!`);
        } else {
            errorAlert(res);
        }

        close(false);
    };
    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] py-6 px-10 w-1/3">
                <h1 className="text-5xl font-bold mb-2 text-center">Heads Up!</h1>
                <p className="text-xl">
                    Are you sure you want to delete{' '}
                    <span className="text-primary">{element.name}</span>? This action is permanent
                    and cannot be undone.
                </p>
                <div className="flex flex-row justify-around">
                    <button
                        className=" border-lightest border-2 rounded-full px-6 py-1 mt-4 w-2/5 font-bold text-2xl text-lighter hover:bg-lighter/30 duration-200"
                        onClick={() => close(false)}
                    >
                        Cancel
                    </button>
                    <button
                        className="bg-error text-white rounded-full px-4 py-2 mt-4 w-2/5 font-bold text-2xl hover:brightness-110 duration-200"
                        onClick={deleteElement}
                    >
                        Delete
                    </button>
                </div>
            </div>
        </>
    );
};

export default DeletePopup;
