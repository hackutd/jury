import Back from '../Back';
import Container from '../Container';
import JuryHeader from '../JuryHeader';

interface InfoPageProps {
    title: string;
    description: string;
}

const JudgeInfoPage = (props: InfoPageProps) => {
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
