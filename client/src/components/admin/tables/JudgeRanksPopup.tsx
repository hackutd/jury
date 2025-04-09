import { useAdminStore } from '../../../store';
import InfoPopup from '../../InfoPopup';

interface JudgeRanksPopupProps {
    enabled: boolean;
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;
    judge: Judge;
}

const JudgeRanksPopup = (props: JudgeRanksPopupProps) => {
    const projects = useAdminStore((state) => state.projects);

    const projEntryFromId = (id: string) => {
        const proj = projects.find((p) => p.id === id);
        return proj ? `${proj.name} [${proj.location}]` : `Unknown (${id})`;
    };

    return (
        <InfoPopup
            enabled={props.enabled}
            setEnabled={props.setEnabled}
            title={`Scores for ${props.judge.name}`}
            submitText="Close"
        >
            <div className="flex flex-col px-4 items-center my-4">
                <h1 className="text-primary text-2xl font-bold">Rankings</h1>
                {props.judge.rankings.map((id, idx) => (
                    <p key={idx} className="text-light">
                        {idx + 1}. {projEntryFromId(id)}
                    </p>
                ))}
                <h1 className="text-primary text-2xl font-bold mt-4">Stars</h1>
                {props.judge.seen_projects
                    .filter((p) => p.starred)
                    .map((proj, idx) => (
                        <p key={idx} className="text-light">
                            {projEntryFromId(proj.project_id)}
                        </p>
                    ))}
            </div>
        </InfoPopup>
    );
};

export default JudgeRanksPopup;
