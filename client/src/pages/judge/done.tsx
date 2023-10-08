import Back from '../../components/Back';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';

const Done = () => {
    return (
        <>
            <JuryHeader withLogout />
            <Container className="px-2 justify-start">
                <Back location="/judge" className="self-start" />
                <h1 className="text-xl font-bold mb-4">You&apos;ve seen all the projects!</h1>
                <p>Click &quot;Back&quot; to view all the projects you&apos;ve judged</p>
            </Container>
        </>
    );
};

export default Done;
