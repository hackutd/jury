import { useAdminStore } from '../../../store';
import InfoPopup from '../../InfoPopup';

interface JudgeNote {
    judgeName: string;
    notes: string;
}

interface ProjectNotesPopupProps {
    enabled: boolean;
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
    project: Project;
}

const ProjectNotesPopup = (props: ProjectNotesPopupProps) => {
    const judges = useAdminStore((state) => state.judges);

    const judgeNotes = judges
        .map((judge) => {
            const seenProject = judge.seen_projects.find(
                (sp) => sp.project_id === props.project.id
            );
            if (seenProject && seenProject.notes && seenProject.notes.trim()) {
                return { judgeName: judge.name, notes: seenProject.notes };
            }
            return null;
        })
        .filter((entry): entry is JudgeNote => !!entry);

    return (
        <InfoPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            title={`Notes for ${props.project.name}`}
            submitText="Close"
        >
            <div className="flex flex-col px-4 items-center my-4">
                {judgeNotes.length === 0 ? (
                    <p className="text-lighter">
                        No judges have left notes for this project.
                    </p>
                ) : (
                    judgeNotes.map((entry, idx) => (
                        <div key={idx} className="w-full mb-4">
                            <h2 className="text-primary font-bold">{entry.judgeName}</h2>
                            <p className="text-light">{entry.notes}</p>
                        </div>
                    ))
                )}
            </div>
        </InfoPopup>
    );
};

export default ProjectNotesPopup;
