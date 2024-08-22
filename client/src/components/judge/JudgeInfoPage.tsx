import Back from '../Back';
import Container from '../Container';
import JuryHeader from '../JuryHeader';

interface JudgeInfoPageProps {
    /* Title to show */
    title: string;
    
    /* Description to show */
    description: string;
}

/**
 * Page to display when judging cannot happen. This could be because
 * of multiple reasons such as a paused session or no more projects left.
 */
const JudgeInfoPage = (props: JudgeInfoPageProps) => {
    return (
        <>
            <JuryHeader withLogout />
            <Container className="px-2 justify-start">
                <Back location="/judge" className="self-start" />
                <h1 className="text-xl font-bold mb-4">{props.title}</h1>
                <p>{props.description}</p>
            </Container>
        </>
    );
};

export default JudgeInfoPage;
