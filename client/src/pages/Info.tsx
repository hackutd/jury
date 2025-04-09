import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';

const Info = ({ children }: { children: React.ReactNode }) => {
    return (
        <>
            <JuryHeader />
            <Container>
                <p className="text-center">{children}</p>
            </Container>
        </>
    );
};

export default Info;
