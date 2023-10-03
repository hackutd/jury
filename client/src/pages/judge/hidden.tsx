import Back from '../../components/Back';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';

const Hidden = () => {
    return (
        <>
            <JuryHeader withLogout />
            <Container className="px-2 justify-start">
                <Back location="/judge" className="self-start" />
                <h1>You are currently disabled</h1>
                <p>Please contact an organizer to enable your account to continue judging.</p>
            </Container>
        </>
    );
};

export default Hidden;
