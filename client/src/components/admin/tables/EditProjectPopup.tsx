import ConfirmPopup from '../../ConfirmPopup';

interface EditProjectPopupProps {
    /* Project to edit */
    project: Project;

    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
}

const EditProjectPopup = (props: EditProjectPopupProps) => {
    return (
        <ConfirmPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            title="Edit Project"
            submitText="Save"
            onSubmit={() => {
                props.setEnabled(false);
            }}
        >
            <p>EDIT WIP</p>
        </ConfirmPopup>
    );
};

export default EditProjectPopup;
