import Back from '../../components/Back';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';

const NotStarted = () => {
    return (
        <>
            <JuryHeader withLogout />
            <Container className="px-2 justify-start">
                <Back location="/judge" className="self-start" />
                <h1 className="text-xl font-bold mb-4">Judging Paused</h1>
                <p>
                    Judging either has not started or is paused. If you think this is an error,
                    please let an event organizer know.
                </p>
            </Container>
        </>
    );
};

export default NotStarted;
