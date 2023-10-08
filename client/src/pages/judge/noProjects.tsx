import Back from '../../components/Back';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';

const NoProjects = () => {
    return (
        <>
            <JuryHeader withLogout />
            <Container className="px-2 justify-start">
                <Back location="/judge" className="self-start" />
                <h1 className='text-xl font-bold mb-4'>There are no projects to judge</h1>
                <p>
                    You&apos;re early! There seems to be no projects currently in the system to
                    judge. Please let an organizer know if this is unexpected.
                </p>
            </Container>
        </>
    );
};

export default NoProjects;
